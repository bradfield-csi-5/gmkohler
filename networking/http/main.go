package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"proxy/pkg"
)

var (
	proxyPort     int
	proxySocket   *pkg.Socket // socket talking to the real server
	welcomePort   int
	welcomeSocket *pkg.Socket // socket welcoming new client connections
)

func init() {
	flag.IntVar(&welcomePort, "p", 80, "port to listen on")
	flag.IntVar(&proxyPort, "f", 8080, "port to forward to")
}

// This program assumes that the proxy socket connection stays alive
func main() {
	flag.Parse()
	var err error
	welcomeSocket, err = pkg.NewTcpSocket()
	if err != nil {
		panic(err)
	}
	defer welcomeSocket.Close()

	welcomeAddr := pkg.NewIpv4Address(welcomePort, [4]byte{})
	if err = welcomeSocket.Bind(welcomeAddr); err != nil {
		panic(fmt.Errorf(
			"error binding fd %v: %w",
			welcomeSocket,
			err,
		))
	}

	proxySocket, err = pkg.NewTcpSocket()
	if err != nil {
		panic(err)
	}
	defer proxySocket.Close()

	proxyAddr := pkg.NewIpv4Address(proxyPort, [4]byte{})
	if err = proxySocket.Connect(proxyAddr); err != nil {
		panic(err)
	}

	err = welcomeSocket.Listen()
	if err != nil {
		panic(fmt.Errorf("error listening on socket %v: %w", welcomeSocket,
			err))
	}

	for {
		connSocket, err := welcomeSocket.Accept()
		if err != nil {
			log.Fatal("error receiving message: ", err)
		}
		go handleConnection(connSocket)
	}
}

func handleConnection(connSocket *pkg.Socket) {
	defer connSocket.Close()

	for {
		var requestBuf bytes.Buffer
		err := connSocket.ReadHttp(&requestBuf)
		if err != nil {
			log.Fatal("error receiving data from client:", err)
		}
		fmt.Printf(
			"request from client: %s\n",
			requestBuf.String(),
		)
		fmt.Println("proxying request to server")
		if err = proxySocket.Write(requestBuf.Bytes()); err != nil {
			log.Fatal("error sending data to proxied server: ", err)
		}
		fmt.Println("request written to server")
		fmt.Println("reading response from server")
		var responseBuf bytes.Buffer
		err = proxySocket.ReadHttp(&responseBuf)
		fmt.Printf(
			"response from server:\n%s\n",
			responseBuf.String(),
		)
		if err != nil {
			log.Fatalf("error reading data from server: %v", err)
		}
		fmt.Println("sending response to client")
		responseBuf.WriteByte('\n')
		if err = connSocket.Write(responseBuf.Bytes()); err != nil {
			log.Fatal("error sending data to client")
		}
	}
}
