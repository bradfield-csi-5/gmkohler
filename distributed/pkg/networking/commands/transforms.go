package commands

import (
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"net"
)

func ReadExecuteCommandRequest(conn net.Conn, requestPtr *ExecuteCommandRequest) error {
	var messageLength int32

	if err := binary.Read(conn, binary.LittleEndian, &messageLength); err != nil {
		return fmt.Errorf(
			"networking/commands: ReadExecuteCommandRequest: error reading message length from connection: %w",
			err,
		)
	}
	var requestBuf = make([]byte, messageLength)
	if _, err := conn.Read(requestBuf); err != nil {
		return fmt.Errorf(
			"networking/commands: ReadExecuteCommandRequest: error reading message from connection: %w",
			err,
		)
	}
	if err := proto.Unmarshal(requestBuf, requestPtr); err != nil {
		return fmt.Errorf(
			"networking/commands: ReadExecuteCommandRequest: error unmarshalling message: %w",
			err,
		)
	}

	return nil
}

func WriteExecuteCommandRequest(conn net.Conn, requestPtr *ExecuteCommandRequest) error {
	encoded, err := proto.Marshal(requestPtr)
	if err != nil {
		return fmt.Errorf("networking/commands: WriteExecuteCommandRequest: error marshalling message: %w", err)
	}
	if err := binary.Write(conn, binary.LittleEndian, int32(len(encoded))); err != nil {
		return fmt.Errorf("networking/commands: WriteExecuteCommandRequest: error writing message length to connection: %w", err)
	}
	if _, err := conn.Write(encoded); err != nil {
		return fmt.Errorf("networking/commands: WriteExecuteCommandRequest: error writing message to connection: %w", err)
	}
	return nil
}

func ReadExecuteCommandResponse(conn net.Conn, responsePtr *ExecuteCommandResponse) error {
	var messageLen int32
	if err := binary.Read(conn, binary.LittleEndian, &messageLen); err != nil {
		return fmt.Errorf("SendRequest: error reading message length from connection: %w", err)
	}
	var responseBuf = make([]byte, messageLen)
	if _, err := conn.Read(responseBuf); err != nil {
		return fmt.Errorf("SendRequest: error reading message from connection: %w", err)
	}
	if err := proto.Unmarshal(responseBuf, responsePtr); err != nil {
		return fmt.Errorf("SendRequest: error unmarshalling message: %w", err)
	}
	return nil
}

func WriteExecuteCommandResponse(conn net.Conn, responsePtr *ExecuteCommandResponse) error {
	encodedResponse, err := proto.Marshal(responsePtr)
	if err != nil {
		return fmt.Errorf("error marshalling response: %w", err)
	}
	if err := binary.Write(conn, binary.LittleEndian, int32(len(encodedResponse))); err != nil {
		return fmt.Errorf("error writing response length to connection: %w", err)
	}
	if _, err := conn.Write(encodedResponse); err != nil {
		return fmt.Errorf("error writing response to connection: %w", err)
	}
	return nil
}
