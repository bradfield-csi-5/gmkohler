package client

import (
	"distributed/pkg/networking"
	"distributed/pkg/storage"
	"encoding/binary"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"net"
)

type Client struct {
	conn net.Conn
}

func New(conn net.Conn) (*Client, error) {
	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) ExecuteCommand(command Command) (storage.Value, error) {
	var requestPtr *networking.ExecuteCommandRequest
	switch command.Operation {
	case OpGet:
		requestPtr = &networking.ExecuteCommandRequest{
			Operation: &networking.ExecuteCommandRequest_Get{
				Get: &networking.GetRequest{Key: string(command.Key)},
			},
		}
	case OpPut:
		requestPtr = &networking.ExecuteCommandRequest{
			Operation: &networking.ExecuteCommandRequest_Put{
				Put: &networking.PutRequest{
					Key:   string(command.Key),
					Value: []byte(command.Value),
				},
			},
		}
	default:
		return "", fmt.Errorf("unrecognized operation %s", command.Operation)
	}
	encoded, err := proto.Marshal(requestPtr)
	if err != nil {
		return "", fmt.Errorf("ExecuteCommand: error marshalling message: %w", err)
	}
	if err := binary.Write(c.conn, binary.LittleEndian, int32(len(encoded))); err != nil {
		return "", fmt.Errorf("ExecuteCommand: error writing message length to connection: %w", err)
	}
	if _, err := c.conn.Write(encoded); err != nil {
		return "", fmt.Errorf("ExecuteCommand: error writing message to connection: %w", err)
	}

	var (
		responsePtr   = new(networking.ExecuteCommandResponse)
		messageLength int32
	)
	if err := binary.Read(c.conn, binary.LittleEndian, &messageLength); err != nil {
		return "", fmt.Errorf("ExecuteCommand: error reading message length from connection: %w", err)
	}
	var responseBuf = make([]byte, messageLength)
	if _, err := c.conn.Read(responseBuf); err != nil {
		return "", fmt.Errorf("ExecuteCommand: error reading message from connection: %w", err)
	}
	if err := proto.Unmarshal(responseBuf, responsePtr); err != nil {
		return "", fmt.Errorf("ExecuteCommand: error unmarshalling message: %w", err)
	}

	switch result := responsePtr.Result.(type) {
	case *networking.ExecuteCommandResponse_Value:
		return storage.Value(result.Value), nil
	case *networking.ExecuteCommandResponse_Error:
		return "", errors.New(result.Error)
	default:
		return "", fmt.Errorf("ExecuteCommand: unrecognized result %T", result)
	}
}
