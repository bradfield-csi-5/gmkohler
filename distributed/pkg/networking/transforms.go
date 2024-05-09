package networking

import (
	"encoding/binary"
	"fmt"
	"net"
)

func WriteMessageRaw(conn net.Conn, data []byte) error {
	if err := binary.Write(conn, binary.LittleEndian, int32(len(data))); err != nil {
		return fmt.Errorf("networking: WriteMessageRaw: %w", err)
	}
	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("networking: WriteMessageRaw: %w", err)
	}

	return nil
}

func ReadMessageRaw(conn net.Conn) ([]byte, error) {
	var messageLen int32
	if err := binary.Read(conn, binary.LittleEndian, &messageLen); err != nil {
		return nil, err
	}
	var responseBuf = make([]byte, messageLen)
	if _, err := conn.Read(responseBuf); err != nil {
		return nil, err
	}
	return responseBuf, nil
}
