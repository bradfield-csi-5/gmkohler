package client

import (
	"distributed/pkg/networking/commands"
	"distributed/pkg/server/storage"
	"errors"
	"fmt"
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
func (c *Client) SendRequest(requestPtr *commands.ExecuteCommandRequest) (*commands.ExecuteCommandResponse, error) {
	if err := commands.WriteExecuteCommandRequest(c.conn, requestPtr); err != nil {
		return nil, err
	}

	var responsePtr = new(commands.ExecuteCommandResponse)
	if err := commands.ReadExecuteCommandResponse(c.conn, responsePtr); err != nil {
		return nil, err
	}

	return responsePtr, nil
}

// ExecuteCommand is a wrapper for clients to send "domain objects" through the wire; more of a "facade" but will think
// of the best way to organize this later.
func (c *Client) ExecuteCommand(command Command) (storage.Value, error) {
	var requestPtr *commands.ExecuteCommandRequest
	switch command.Operation {
	case OpGet:
		requestPtr = &commands.ExecuteCommandRequest{
			Operation: &commands.ExecuteCommandRequest_Get{
				Get: &commands.GetRequest{Key: string(command.Key)},
			},
		}
	case OpPut:
		requestPtr = &commands.ExecuteCommandRequest{
			Operation: &commands.ExecuteCommandRequest_Put{
				Put: &commands.PutRequest{
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
	case *commands.ExecuteCommandResponse_Value:
		return storage.Value(result.Value), nil
	case *commands.ExecuteCommandResponse_Error:
		return "", errors.New(result.Error)
	default:
		return "", fmt.Errorf("ExecuteCommand: unrecognized result %T", result)
	}
}
