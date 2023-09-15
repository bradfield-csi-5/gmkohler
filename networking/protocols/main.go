package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type magicNumber uint32

func (m magicNumber) String() string {
	return strconv.FormatInt(int64(m), 16)
}

type globalHeader struct {
	MagicNumber    magicNumber
	MajorVersion   uint16
	MinorVersion   uint16
	Reserved1      uint32
	Reserved2      uint32
	SnapshotLength uint32
	LinkLayerInfo  uint32
}

type packetHeader struct {
	TimestampSeconds  int32
	TimestampMillis   int32
	LengthTruncated   uint32
	LengthUntruncated uint32
}

type macAddress [6]uint8

type macHeader struct {
	MacDestination macAddress
	MacSource      macAddress
	EtherType      uint16
}

func (m *macHeader) IPv4() bool {
	return m.EtherType == etherTypeIPv4
}
func (m *macHeader) IPv6() bool {
	return m.EtherType == etherTypeIPv6
}

type ipv4Version uint8

// String() unpacks the byte into the two fields it represents
func (i ipv4Version) String() string {
	return fmt.Sprintf(
		"{Version:%d HeaderLength:%d}",
		i.version(),
		i.headerLength(),
	)
}

func (i ipv4Version) version() uint8 {
	return uint8(i >> 4)
}
func (i ipv4Version) headerLength() uint8 {
	return uint8(i&0x0f) * 4
}

type ipAddress uint32

func (i ipAddress) String() string {
	return fmt.Sprintf(
		"%d.%d.%d.%d",
		i&0xff000000>>24,
		i&0x00ff0000>>16,
		i&0x0000ff00>>8,
		i&0x000000ff,
	)
}

type ipv4Header struct {
	VersionAndHeaderLength ipv4Version
	DscpEcn                uint8
	TotalLength            uint16
	Identification         uint16
	FlagsAndOffset         uint16
	Ttl                    uint8
	Protocol               uint8
	Checksum               uint16
	SourceIp               ipAddress
	DestIp                 ipAddress
	// options []byte
}

func (i ipv4Header) DataLength() uint64 {
	return uint64(
		i.TotalLength - uint16(i.VersionAndHeaderLength.headerLength()))
}

type tcpPayload []byte
type sequenceNumber uint32
type tcpHeader struct {
	SourcePort     uint16
	DestPort       uint16
	SequenceNumber sequenceNumber
	AckNumber      uint32
	DataOffset     uint8
	Flags          uint8
	WindowSize     uint16
	Checksum       uint16
	UrgentPointer  uint16
	// options []byte
}

// HeaderLength decodes DataOffset to tell us the length in bytes of the header,
// meaning that after this is the TCP Payload.
func (t *tcpHeader) HeaderLength() int64 {
	return int64(t.DataOffset>>4) * 4
}

const (
	etherTypeIPv4    = 0x0800
	etherTypeIPv6    = 0x86DD
	macHeaderSize    = 14
	tcpHeaderSize    = 20
	baseIpHeaderSize = 20
	ipV4             = 4
	tcpProtocol      = 6
)

func main() {
	gHeader := new(globalHeader)
	pHeader := new(packetHeader)
	mHeader := new(macHeader)
	ipHeader := new(ipv4Header)
	tHeader := new(tcpHeader)
	var clientIp ipAddress
	var hostIp ipAddress
	requestPayloads := make(map[sequenceNumber]tcpPayload)
	responsePayloads := make(map[sequenceNumber]tcpPayload)

	b, err := os.ReadFile("net.cap")
	if err != nil {
		panic(err)
	}
	buf := bytes.NewReader(b)
	err = binary.Read(
		buf,
		binary.LittleEndian,
		gHeader,
	)
	if err != nil {
		log.Fatal("Error reading global header: ", err)
	}
	fmt.Printf("%+v\n", gHeader)
	numPackets := 0
	for {
		// Packet Header
		err = binary.Read(buf, binary.LittleEndian, pHeader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal("Unexpected error reading packet header: ", err)
		}
		if pHeader.LengthUntruncated != pHeader.LengthTruncated {
			log.Fatalf("Packet header is truncated: %+v", pHeader)
		}
		remainingBytesInPacket := pHeader.LengthTruncated
		fmt.Printf("Packet Header: %+v\n", pHeader)

		// MAC Header
		err = binary.Read(buf, binary.BigEndian, mHeader)
		if err != nil {
			log.Fatal("error parsing MAC Header: ", err)
		}
		if !mHeader.IPv4() {
			log.Fatalf("Unexpected EtherType %x", mHeader.EtherType)
		}
		fmt.Printf("MAC Header: %+v\n", mHeader)
		remainingBytesInPacket -= macHeaderSize

		// IP Header
		err = binary.Read(buf, binary.BigEndian, ipHeader)
		if err != nil {
			log.Fatal("error parsing IP Header: ", err)
		}
		if ipHeader.VersionAndHeaderLength.version() != ipV4 {
			log.Fatalf(
				"unexpected IP version %d (expected %d)",
				ipHeader.VersionAndHeaderLength.version(),
				ipV4,
			)
		}
		if ipHeader.Protocol != tcpProtocol {
			log.Fatalf(
				"unexpected protocol number %d (expected %d)",
				ipHeader.Protocol,
				tcpProtocol,
			)
		}
		if numPackets == 0 {
			clientIp = ipHeader.SourceIp
			hostIp = ipHeader.DestIp
		}
		numPackets++
		fmt.Printf("IP Header: %+v\n", ipHeader)
		// Scan over options in header because we aren't using them in this
		// exercise (and none of the packets in this file even have the
		// options), and go straight to the Data section of the packet:
		_, err = buf.Seek(
			int64(ipHeader.VersionAndHeaderLength.headerLength()-baseIpHeaderSize),
			io.SeekCurrent,
		)
		// We didn't decrement when reading the base header,
		// so here we decrement based on the entire length encoded in said
		// header.
		remainingBytesInPacket -=
			uint32(ipHeader.VersionAndHeaderLength.headerLength())

		// TCP header
		err = binary.Read(buf, binary.BigEndian, tHeader)
		if err != nil {
			log.Fatal("error parsing TCP header: ", err)
		}

		fmt.Printf("TCP header: %+v\n", tHeader)
		// pass through the options that we ignore
		_, err := buf.Seek(tHeader.HeaderLength()-tcpHeaderSize, io.SeekCurrent)
		if err != nil {
			log.Fatal("error seeking past TCP options: ", err)
		}
		// we didn't decrement when reading the base headers, so we can use the
		// specified length here.
		remainingBytesInPacket -= uint32(tHeader.HeaderLength())
		if remainingBytesInPacket == 0 {
			continue
		}
		tPayload := make(tcpPayload, remainingBytesInPacket)
		read, err := buf.Read(tPayload)
		if err != nil || read != len(tPayload) {
			log.Fatal("error reading TCP Payload")
		}

		// keep request and response streams separate
		if ipHeader.SourceIp == clientIp {
			_, exists := requestPayloads[tHeader.SequenceNumber]
			if !exists {
				requestPayloads[tHeader.SequenceNumber] = tPayload
			}
		} else if ipHeader.SourceIp == hostIp {
			_, exists := responsePayloads[tHeader.SequenceNumber]
			if !exists {
				responsePayloads[tHeader.SequenceNumber] = tPayload
			}
		} else {
			log.Fatalf("unexpected source IP %d", ipHeader.SourceIp)
		}
	}

	fmt.Printf("Number of packets: %d\n", numPackets)

	// We have built a map of the packets (using map to deduplicate and manage
	// any out-of-order receipt).  Now, we sort by the sequence numbers to
	// stitch together the payloads and build the HTTP request/response objects
	// from these payloads.
	requestSequence := maps.Keys(requestPayloads)
	slices.Sort(requestSequence)
	fmt.Printf("request sequences: %v\n", requestSequence)
	var requestPayload tcpPayload
	for _, sn := range requestSequence {
		requestPayload = append(requestPayload, requestPayloads[sn]...)
	}
	req, err := http.ReadRequest(bufio.NewReaderSize(
		bytes.NewReader(requestPayload),
		len(requestPayload),
	))
	if err != nil {
		log.Fatal("error reading request: ", err)
	}
	fmt.Printf("request: %+v\n", req)
	responseSequence := maps.Keys(responsePayloads)
	slices.Sort(responseSequence)

	fmt.Printf("response sequences: %v\n", responseSequence)

	var responsePayload tcpPayload
	for _, sn := range responseSequence {
		responsePayload = append(responsePayload, responsePayloads[sn]...)
	}
	// slicing from 2: is a hack because I'm seeing $ and £ characters at the
	// beginning of the response payload.
	resp, err := http.ReadResponse(
		bufio.NewReaderSize(
			bytes.NewReader(responsePayload[2:]),
			len(responsePayload),
		),
		req,
	)
	if err != nil {
		log.Fatal("error reading response: ", err)
	}
	fmt.Printf("response: %+v\n", resp)

	// We have a response with a body, let's get the image from it and save it.
	decoded, err := jpeg.Decode(resp.Body)
	if err != nil {
		log.Fatal("error decoding jpeg from payload: ", err)
	}
	f, err := os.Create("img.jpeg")
	if err != nil {
		log.Fatal("error creating file for saving image: ", err)
	}

	err = jpeg.Encode(f, decoded, nil)
	if err != nil {
		log.Fatal("error encoding image: ", err)
	}
}
