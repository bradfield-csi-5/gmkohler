package server

import (
	"distributed/pkg/networking"
	"net"
)

func NewForwardedServer(connection net.Conn) (*ForwardedServer, error) {
	return &ForwardedServer{conn: connection}, nil
}

// ForwardedServer is a simple wrapper over a connection for reading/writing raw data.  Not sure I need this struct now
// that I have the networking.(Read|Write)MessageRaw functions but will leave it around in case more needs to be known
// about each of the downstream servers e.g. partitions.
type ForwardedServer struct {
	conn net.Conn
}

func (fs *ForwardedServer) SendMessage(encodedRequest []byte) ([]byte, error) {
	if err := networking.WriteMessageRaw(fs.conn, encodedRequest); err != nil {
		return nil, err
	}

	rawResponse, err := networking.ReadMessageRaw(fs.conn)
	if err != nil {
		return nil, err
	}

	return rawResponse, nil
}
