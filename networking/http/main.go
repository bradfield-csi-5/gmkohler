package main

import (
	"flag"
	"fmt"
	"log"
	"syscall"
)

var (
	clientFd   int
	clientPort int
	serverFd   int
	serverPort int
)

func init() {
	flag.IntVar(&clientPort, "p", 80, "port to listen on")
	flag.IntVar(&serverPort, "f", 81, "port to forward to")
}

func main() {
	flag.Parse()
	var err error
	clientFd, err = openTcpSocket()
	defer closeSocket(clientFd)
	if err != nil {
		panic(err)
	}
	if err = syscall.Bind(
		clientFd,
		&syscall.SockaddrInet4{
			Port: clientPort,
			Addr: [4]byte{},
		},
	); err != nil {
		panic(fmt.Errorf(
			"error binding fd %d: %w",
			clientFd,
			err,
		))
	}

	serverFd, err = openTcpSocket()
	defer closeSocket(serverFd)
	if err != nil {
		panic(err)
	}
	if err = syscall.Connect(
		serverFd,
		&syscall.SockaddrInet4{
			Port: serverPort,
			Addr: [4]byte{},
		},
	); err != nil {
		panic(err)
	}

	err = syscall.Listen(clientFd, 128)
	if err != nil {
		panic(fmt.Errorf("error listening on fd %d: %w", clientFd, err))
	}

	for {
		nfd, _, err := syscall.Accept(clientFd)
		if err != nil {
			log.Fatal("error receiving message: ", err)
		}
		go handleConnection(nfd)
	}
}

func closeSocket(fd int) {
	if err := syscall.Close(fd); err != nil {
		log.Fatal("error closing socket: ", err)
	}
}

func openTcpSocket() (int, error) {
	fd, err := syscall.Socket(
		syscall.AF_INET,
		syscall.SOCK_STREAM,
		syscall.IPPROTO_TCP,
	)
	if err != nil {
		return -1, fmt.Errorf("openTcpSocket(): %w", err)
	}
	return fd, nil
}

func handleConnection(fd int) {
	defer closeSocket(fd)

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
			serverFd,
			buf[:bytesRead],
			0,
			fromAddr,
		); err != nil {
			log.Fatal("error sending data to proxied server: ", err)
		}
	}
}
