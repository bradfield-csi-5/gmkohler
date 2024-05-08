package server

import (
	"distributed/pkg/client"
	"distributed/pkg/networking"
	"distributed/pkg/server/storage"
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
)

type Server struct {
	listener       net.Listener
	db             storage.Storage
	replicaClients []*client.Client
}

func New(listener net.Listener, db storage.Storage, replicaSocketPaths ...string) (*Server, error) {
	var replicaClients = make([]*client.Client, len(replicaSocketPaths))
	if len(replicaSocketPaths) > 0 {
		for j, replicaSocketPath := range replicaSocketPaths {
			replicaConn, err := net.Dial(networking.Unix, replicaSocketPath)
			if err != nil {
				return nil, fmt.Errorf(
					"server.New: error connecting to replica %s: %w",
					replicaSocketPath,
					err,
				)
			}
			log.Printf("connected to replica at %s", replicaConn.RemoteAddr())

			replicaClient, err := client.New(replicaConn)
			if err != nil {
				return nil, fmt.Errorf(
					"server.New: error building replica client for %s: %w",
					replicaSocketPath,
					err,
				)

			}
			replicaClients[j] = replicaClient
		}
	}

	return &Server{
		listener:       listener,
		db:             db,
		replicaClients: replicaClients,
	}, nil
}

func (s *Server) Run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("Server.Run(): error accepting connection", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	log.Printf("handling connection from %s", conn.RemoteAddr())
	for {
		var (
			requestPtr    = new(networking.ExecuteCommandRequest)
			messageLength int32
		)

		if err := binary.Read(conn, binary.LittleEndian, &messageLength); err != nil {
			log.Printf("error reading message length from connection: %v", err)
			conn.Close()
			break
		}
		var requestBuf = make([]byte, messageLength)
		if _, err := conn.Read(requestBuf); err != nil {
			log.Printf("error reading message from connection: %v", err)
			conn.Close()
			break
		}
		if err := proto.Unmarshal(requestBuf, requestPtr); err != nil {
			log.Printf("error decoding request: %v", err)
			continue
		}
		var (
			value storage.Value
			err   error
		)
		switch req := requestPtr.Operation.(type) {
		case *networking.ExecuteCommandRequest_Get:
			value, err = s.db.Get(storage.Key(req.Get.Key))
		case *networking.ExecuteCommandRequest_Put:
			value, err = s.db.Put(storage.Key(req.Put.Key), storage.Value(req.Put.Value))
			for _, replica := range s.replicaClients {
				if _, err := replica.SendRequest(requestPtr); err != nil {
					log.Printf("error sending request to replica: %v", err)
				}
			}
		default:
			err = fmt.Errorf("unrecognized operation: %v", req)
		}

		var responsePtr *networking.ExecuteCommandResponse
		if err != nil {
			responsePtr = &networking.ExecuteCommandResponse{
				Result: &networking.ExecuteCommandResponse_Error{
					Error: err.Error(),
				},
			}
		} else {
			responsePtr = &networking.ExecuteCommandResponse{
				Result: &networking.ExecuteCommandResponse_Value{
					Value: []byte(value),
				},
			}
		}
		encodedResponse, err := proto.Marshal(responsePtr)
		if err != nil {
			log.Printf("error marshalling response: %v", err)
			continue
		}
		if err := binary.Write(conn, binary.LittleEndian, int32(len(encodedResponse))); err != nil {
			log.Printf("error writing response length to connection: %v", err)
			conn.Close()
			break
		}
		if _, err := conn.Write(encodedResponse); err != nil {
			log.Printf("error writing response to connection: %v", err)
			conn.Close()
			break
		}
	}
}
