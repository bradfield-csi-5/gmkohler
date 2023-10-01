package pkg

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"syscall"
)

var (
	endOfHeaders = [4]byte{'\r', '\n', '\r', '\n'}
)

type Socket struct {
	fd FileDescriptor
}
type FileDescriptor int
type Address struct {
	syscall.Sockaddr
}

func NewIpv4Address(port int, addr [4]byte) Address {
	return Address{
		&syscall.SockaddrInet4{
			Addr: addr,
			Port: port,
		},
	}
}

func (s *Socket) String() string {
	return strconv.Itoa(int(s.fd))
}

func NewTcpSocket() (*Socket, error) {
	fd, err := syscall.Socket(
		syscall.AF_INET,
		syscall.SOCK_STREAM,
		syscall.IPPROTO_TCP,
	)
	if err != nil {
		return nil, fmt.Errorf("openTcpSocket(): %w", err)
	}
	return &Socket{
		fd: FileDescriptor(fd),
	}, nil
}

// Accept returns a connection Socket from the process' TCP "welcome socket"
func (s *Socket) Accept() (*Socket, error) {
	fd, _, err := syscall.Accept(int(s.fd))
	if err != nil {
		return nil, fmt.Errorf("Socket.Accept(): %w", err)
	}
	return &Socket{fd: FileDescriptor(fd)}, nil
}

// Connect attaches an external Address to a Socket
func (s *Socket) Connect(a Address) error {
	if err := syscall.Connect(int(s.fd), a); err != nil {
		return fmt.Errorf("Socket.Connect(): %w", err)
	}
	return nil
}

func (s *Socket) Bind(a Address) error {
	if err := syscall.Bind(int(s.fd), a); err != nil {
		return fmt.Errorf("Socket.Bind(): %w", err)
	}
	return nil
}

func (s *Socket) Listen() error {
	err := syscall.Listen(int(s.fd), 128)
	if err != nil {
		return fmt.Errorf("Socket.Listen(): %w", err)
	}
	return nil
}

func (s *Socket) Close() {
	if err := syscall.Close(int(s.fd)); err != nil {
		log.Fatal("error closing socket: ", err)
	}
}

// ReadHttp assumes for now that we are reading fast enough that only one
// request is in the buffer
func (s *Socket) ReadHttp(buf *bytes.Buffer) error {
	recvBuf := make([]byte, 0x1000)
	for {
		bytesRead, _, err := syscall.Recvfrom(int(s.fd), recvBuf, 0)
		if err != nil {
			return fmt.Errorf("Socket.ReadHttp(): %w", err)
		}
		if bytesRead == 0 {
			fmt.Println("Socket.ReadHttp(): client disconnected")
			return nil
		}

		buf.Write(recvBuf[:bytesRead])
		fmt.Printf("buf:\n%s\n", buf.String())
		if bytes.Contains(buf.Bytes(), endOfHeaders[:]) {
			fmt.Println("found 2 CRLF.  ending reading")
			break
		}
	}

	return nil
}

func (s *Socket) Read(buf []byte) (int, error) {
	bytesRead, _, err := syscall.Recvfrom(int(s.fd), buf, 0)
	if err != nil {
		return 0, fmt.Errorf("Socket.Read(): %w", err)
	}
	return bytesRead, nil
}

func (s *Socket) Write(buf []byte) error {
	err := syscall.Sendto(int(s.fd), buf, 0, nil)
	if err != nil {
		return fmt.Errorf("Socket.Write(): %w", err)
	}
	return nil
}
