package main

import (
	"flag"
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
		buf := make([]byte, 0x1000)
		bytesRead, fromAddr, err := syscall.Recvfrom(nfd, buf, 0)
		if err != nil {
			log.Fatal("error receiving from file descriptor:", err)
		}

		if err = syscall.Sendto(
			nfd,
			buf[:bytesRead],
			0,
			fromAddr,
		); err != nil {
			log.Fatal("error sending data back: ", err)
		}

		if err = syscall.Close(nfd); err != nil {
			log.Fatal("error closing client file descriptor: ", err)
		}
	}
}
