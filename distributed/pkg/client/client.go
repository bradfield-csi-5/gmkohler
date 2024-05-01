package client

import (
	"distributed/pkg/networking"
	"distributed/pkg/storage"
	"encoding/gob"
	"errors"
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

func (c *Client) ExecuteCommand(command networking.Command) (storage.Value, error) {
	if err := gob.NewEncoder(c.conn).Encode(command); err != nil {
		return "", err
	}
	var response networking.ExecuteCommandResponse
	if err := gob.NewDecoder(c.conn).Decode(&response); err != nil {
		return "", err
	}

	var responseErr error
	if len(response.Err) > 0 {
		responseErr = errors.New(response.Err)
	}

	return response.Value, responseErr
}
