package main

import (
	"dns/pkg"
	"fmt"
	"log"
	"os"
	"syscall"
)

const (
	dnsPort           = 53
	destinationServer = "8.8.8.8" // google public DNS
)

func main() {
	ip, err := pkg.IpAddrFromString(destinationServer)
	if err != nil {
		log.Fatal("error parsing IP address: ", err)
	}
	dnsAddr := syscall.SockaddrInet4{
		Addr: ip,
		Port: dnsPort,
	}
	s, err := syscall.Socket(
		syscall.AF_INET,
		syscall.SOCK_DGRAM,
		syscall.IPPROTO_UDP,
	)
	if err != nil {
		log.Fatal("error opening socket: ", err)
	}
	fmt.Printf("opened socket %d\n", s)
	defer func(fd int) {
		fmt.Printf("closing socket %d\n", fd)
		closeErr := syscall.Close(fd)
		if closeErr != nil {
			log.Fatal("error closing socket: ", closeErr)
		}
	}(s)
	msg := pkg.Query(
		pkg.ResourceName(os.Args[1]),
		pkg.ResourceTypeFromString(os.Args[2]),
	)
	fmt.Printf("message: %+v\n", msg)

	encoded, err := msg.Encode()
	if err != nil {
		log.Fatalf("error encoding message: %v", err)
	}
	err = syscall.Bind(
		s,
		&syscall.SockaddrInet4{Port: 0, Addr: [4]byte{0, 0, 0, 0}},
	)
	if err != nil {
		log.Fatalf("error binding to a socket: %v", err)
	}
	err = syscall.Sendto(s, encoded, 0, &dnsAddr)
	fmt.Println("message sent to socket")

	if err != nil {
		log.Fatal("error sending to socket: ", err)
	}
	received := make([]byte, 0x1111)
	for {
		_, recvAddr, err := syscall.Recvfrom(s, received, 0)
		if err != nil {
			log.Fatal("error receiving from socket: ", err)
		}
		ip4Addr, ok := recvAddr.(*syscall.SockaddrInet4)
		if !ok {
			continue
		}
		if ip4Addr.Addr != dnsAddr.Addr || ip4Addr.Port != dnsAddr.Port {
			continue
		}

		//fmt.Printf("received from socket: %d %+v %+v\n", n, recvAddr, received)
		response, err := pkg.DecodeMessage(received)
		fmt.Printf("received message: %+v\n", response)
		if err != nil {
			log.Fatalf("error decoding message: %+v", err)
		}
		if code := response.Header.Flags.Code(); code != pkg.Ok {
			fmt.Printf("non-OK response code: %v\n", code)
		}
		if response.Header.Identification != msg.Header.Identification {
			continue
		}
		break
	}
}
