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

func main() {
	gHeader := new(globalHeader)
	pHeader := new(packetHeader)
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
		log.Fatal("Error reading global header", err)
	}
	fmt.Printf("%+v\n", gHeader)
	numPackets := 0
	for {
		err = binary.Read(buf, binary.LittleEndian, pHeader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal("Unexpected error reading packet header", err)
		}
		numPackets++
		if pHeader.LengthUntruncated != pHeader.LengthTruncated {
			log.Fatalf("Packet header is truncated: %+v", pHeader)
		}
		fmt.Printf("%+v\n", pHeader)
		_, err := buf.Seek(int64(pHeader.LengthTruncated), io.SeekCurrent)
		if err != nil {
			log.Fatal("Unexpected error seeking in file", err)
		}
	}
	fmt.Printf("Number of packets: %d", numPackets)
}
