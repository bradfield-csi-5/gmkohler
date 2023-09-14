package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
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

type tcpHeader struct {
	SourcePort     uint16
	DestPort       uint16
	SequenceNumber uint32
	AckNumber      uint32
	DataOffset     uint8
	Flags          uint8
	WindowSize     uint8
	Checksum       uint8
	UrgentPointer  uint8
}

func (t *tcpHeader) HeaderLength() int64 {
	return int64(t.DataOffset>>4) * 4
}
func (t *tcpHeader) DataLength() int64 {
	return t.HeaderLength() - tcpHeaderSize
}

const (
	etherTypeIPv4 = 0x0800
	etherTypeIPv6 = 0x86DD
	macHeaderSize = 14
	tcpHeaderSize = 20
	ipHeaderSize  = 20
	ipV4          = 4
	tcpProtocol   = 6
)

func main() {
	gHeader := new(globalHeader)
	pHeader := new(packetHeader)
	mHeader := new(macHeader)
	ipHeader := new(ipv4Header)
	//tHeader := new(tcpHeader)

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
		fmt.Printf("IP Header: %+v\n", ipHeader)

		_, err = buf.Seek(
			int64(pHeader.LengthTruncated)-macHeaderSize-ipHeaderSize,
			io.SeekCurrent,
		)
		if err != nil {
			log.Fatal("Unexpected error seeking in file: ", err)
		}
		numPackets++
	}
	fmt.Printf("Number of packets: %d", numPackets)
}
