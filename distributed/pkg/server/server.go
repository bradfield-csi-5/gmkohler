package server

import (
	"distributed/pkg/networking"
	"distributed/pkg/storage"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type Server struct {
	listener net.Listener
	db       storage.Storage
}

func New(listener net.Listener, db storage.Storage) (*Server, error) {
	return &Server{
		listener: listener,
		db:       db,
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
	var cmdPtr = new(networking.Command)
	for {
		if err := gob.NewDecoder(conn).Decode(cmdPtr); err != nil {
			log.Printf("error decoding message: %v", err)
			conn.Close()
			break
		}
		var (
			value storage.Value
			err   error
		)
		switch cmdPtr.Operation {
		case networking.OpGet:
			value, err = s.db.Get(cmdPtr.Key)
		case networking.OpPut:
			value, err = s.db.Put(cmdPtr.Key, cmdPtr.Value)
		default:
			err = fmt.Errorf("unrecognized operation: %v", cmdPtr.Operation)
		}

		var response networking.ExecuteCommandResponse
		if err != nil {
			response = networking.ExecuteCommandResponse{Err: err.Error()}
		} else {
			response = networking.ExecuteCommandResponse{Value: value}
		}

		if err = gob.NewEncoder(conn).Encode(response); err != nil {
			log.Printf("error encoding response: %v", err)
			conn.Close()
			break
		}
	}
}
