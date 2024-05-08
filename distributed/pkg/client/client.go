package client

import (
	"distributed/pkg/networking"
	"distributed/pkg/server/storage"
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

// SendRequest is exported so primary servers can relay requests to replicas without having to deal with
// re-serialization
func (c *Client) SendRequest(requestPtr *networking.ExecuteCommandRequest) (*networking.ExecuteCommandResponse, error) {
	encoded, err := proto.Marshal(requestPtr)
	if err != nil {
		return nil, fmt.Errorf("SendRequest: error marshalling message: %w", err)
	}
	if err := binary.Write(c.conn, binary.LittleEndian, int32(len(encoded))); err != nil {
		return nil, fmt.Errorf("SendRequest: error writing message length to connection: %w", err)
	}
	if _, err := c.conn.Write(encoded); err != nil {
		return nil, fmt.Errorf("SendRequest: error writing message to connection: %w", err)
	}

	var (
		responsePtr   = new(networking.ExecuteCommandResponse)
		messageLength int32
	)
	if err := binary.Read(c.conn, binary.LittleEndian, &messageLength); err != nil {
		return nil, fmt.Errorf("SendRequest: error reading message length from connection: %w", err)
	}
	var responseBuf = make([]byte, messageLength)
	if _, err := c.conn.Read(responseBuf); err != nil {
		return nil, fmt.Errorf("SendRequest: error reading message from connection: %w", err)
	}
	if err := proto.Unmarshal(responseBuf, responsePtr); err != nil {
		return nil, fmt.Errorf("SendRequest: error unmarshalling message: %w", err)
	}
	return responsePtr, nil

}

// ExecuteCommand is a wrapper for clients to send "domain objects" through the wire; more of a "facade" but will think
// of the best way to organize this later.
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
	responsePtr, err := c.SendRequest(requestPtr)
	if err != nil {
		return "", err
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
