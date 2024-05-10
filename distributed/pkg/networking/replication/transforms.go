package replication

import (
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"net"
)

func WriteAddReplicaRequest(conn net.Conn, requestPtr *AddReplicaRequest) error {
	encoded, err := proto.Marshal(requestPtr)
	if err != nil {
		return fmt.Errorf("networking/replication: WriteAddReplicaRequest: %w", err)
	}
	if err := binary.Write(conn, binary.LittleEndian, int32(len(encoded))); err != nil {
		return fmt.Errorf("networking/replication: WriteAddReplicaRequest: %w", err)
	}
	if _, err := conn.Write(encoded); err != nil {
		return fmt.Errorf("networking/replication: WriteAddReplicaRequest: %w", err)
	}
	return nil
}

func ReadAddReplicaRequest(conn net.Conn, requestPtr *AddReplicaRequest) error {
	var msgLen int32
	if err := binary.Read(conn, binary.LittleEndian, &msgLen); err != nil {
		return fmt.Errorf("networking/replication: ReadAddReplicaRequest: %w", err)
	}
	var reqBuf = make([]byte, msgLen)
	if _, err := conn.Read(reqBuf); err != nil {
		return fmt.Errorf("networking/replication: ReadAddReplicaRequest: %w", err)
	}
	if err := proto.Unmarshal(reqBuf, requestPtr); err != nil {
		return fmt.Errorf("networking/replication: ReadAddReplicaRequest: %w", err)
	}
	return nil
}

func ReadAddReplicaResponse(conn net.Conn, responsePtr *AddReplicaResponse) error {
	var messageLen int32
	if err := binary.Read(conn, binary.LittleEndian, &messageLen); err != nil {
		return fmt.Errorf("networking/replication: ReadAddReplicaResponse: %w", err)
	}
	var responseBuf = make([]byte, messageLen)
	if _, err := conn.Read(responseBuf); err != nil {
		return fmt.Errorf("networking/replication: ReadAddReplicaResponse: %w", err)
	}
	if err := proto.Unmarshal(responseBuf, responsePtr); err != nil {
		return fmt.Errorf("networking/replication: ReadAddReplicaResponse: %w", err)
	}
	return nil
}

func WriteAddReplicaResponse(conn net.Conn, responsePtr *AddReplicaResponse) error {
	encodedResponse, err := proto.Marshal(responsePtr)
	if err != nil {
		return fmt.Errorf("error marshalling response: %w", err)
	}
	if err := binary.Write(conn, binary.LittleEndian, int32(len(encodedResponse))); err != nil {
		return fmt.Errorf("error writing response length to connection: %w", err)
	}
	if _, err = conn.Write(encodedResponse); err != nil {
		return fmt.Errorf("error writing response length to connection: %w", err)
	}
	return nil
}
