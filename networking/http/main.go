package main

import (
	"flag"
	"fmt"
	"log"
	"syscall"
)

var port int

func init() {
	flag.IntVar(&port, "p", 80, "port to listen on")
}

func main() {
	flag.Parse()
	s, err := syscall.Socket(
		syscall.AF_INET,
		syscall.SOCK_STREAM,
		syscall.IPPROTO_TCP,
	)
	defer func(fd int) {
		if err := syscall.Close(fd); err != nil {
			log.Fatal("error closing socket: ", err)
		}
	}(s)
	if err != nil {
		log.Fatal("error listening on port: ", err)
	}
	err = syscall.Bind(
		s,
		&syscall.SockaddrInet4{Port: port, Addr: [4]byte{}},
	)
	if err != nil {
		log.Fatal("error binding socket: ", err)
	}
	err = syscall.Listen(s, 128)
	if err != nil {
		log.Fatal("error listening on socket: ", err)
	}
	for {
		nfd, _, err := syscall.Accept(s)
		if err != nil {
			log.Fatal("error receiving message: ", err)
		}
		go handleConnection(nfd)
	}
}

func handleConnection(fd int) {
	defer func(fd int) {
		if err := syscall.Close(fd); err != nil {
			log.Fatalf(
				"failure closing socket %d: %v",
				fd,
				err,
			)
		}
	}(fd)

	for {
		buf := make([]byte, 0x1000)
		bytesRead, fromAddr, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			log.Fatal("error receiving from file descriptor:", err)
		}
		if bytesRead == 0 {
			fmt.Println("no bytes to be read. closing connection.")
			break
		}
		if err = syscall.Sendto(
			fd,
			buf[:bytesRead],
			0,
			fromAddr,
		); err != nil {
			log.Fatal("error sending data back: ", err)
		}
	}
}
