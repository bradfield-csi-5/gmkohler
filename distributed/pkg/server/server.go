package server

import (
	"distributed/pkg/client"
	"distributed/pkg/networking"
	"distributed/pkg/networking/commands"
	"distributed/pkg/networking/replication"
	"distributed/pkg/server/storage"
	"fmt"
	"log"
	"net"
	"time"
)

type Server struct {
	logger         *log.Logger
	listener       net.Listener
	db             storage.Storage
	replicaClients []*client.Client
}

func New(
	db storage.Storage,
	commandListener net.Listener,
	replicaListener net.Listener,
	logger *log.Logger,
) (*Server, error) {
	server := &Server{
		logger:   logger,
		listener: commandListener,
		db:       db,
	}

	if replicaListener != nil {
		go server.handleReplicaListening(replicaListener)
	}

	return server, nil
}

func (s *Server) Run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.Println("Server.Run(): error accepting connection", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	s.logger.Printf("handling connection from %s", conn.RemoteAddr())
	for {
		var requestPtr = new(commands.ExecuteCommandRequest)
		err := commands.ReadExecuteCommandRequest(conn, requestPtr)
		if err != nil {
			s.logger.Printf("error reading message length from connection: %v", err)
			conn.Close()
			break
		}

		var value storage.Value
		switch req := requestPtr.Operation.(type) {
		case *commands.ExecuteCommandRequest_Get:
			s.logger.Printf("received get request for %q", req.Get.Key)
			value, err = s.db.Get(storage.Key(req.Get.Key))
		case *commands.ExecuteCommandRequest_Put:
			value, err = s.db.Put(storage.Key(req.Put.Key), storage.Value(req.Put.Value))
			for _, replica := range s.replicaClients {
				if _, err := replica.SendRequest(requestPtr); err != nil {
					s.logger.Printf("error sending request to replica: %v", err)
				}
			}
		default:
			err = fmt.Errorf("unrecognized operation: %v", req)
		}

		var responsePtr *commands.ExecuteCommandResponse
		if err != nil {
			responsePtr = &commands.ExecuteCommandResponse{
				Result: &commands.ExecuteCommandResponse_Error{
					Error: err.Error(),
				},
			}
		} else {
			responsePtr = &commands.ExecuteCommandResponse{
				Result: &commands.ExecuteCommandResponse_Value{
					Value: []byte(value),
				},
			}
		}

		if err := commands.WriteExecuteCommandResponse(conn, responsePtr); err != nil {
			s.logger.Printf("error marshalling response: %v", err)
			conn.Close()
			break
		}
	}
}

func (s *Server) handleReplicaListening(socket net.Listener) {
	s.logger.Printf("listening at socket %s", socket.Addr())
	for {
		conn, err := socket.Accept()
		if err != nil {
			s.logger.Printf("error accepting connection: %v", err)
		}
		var requestPtr = new(replication.AddReplicaRequest)
		if err := replication.ReadAddReplicaRequest(conn, requestPtr); err != nil {
			s.logger.Printf("error reading message length from connection: %v", err)
			conn.Close()
			break
		}

		var responsePtr = new(replication.AddReplicaResponse)
		if err = s.addReplica(requestPtr.SocketPath); err != nil {
			responsePtr.Error = err.Error()
		}

		if err := replication.WriteAddReplicaResponse(conn, responsePtr); err != nil {
			s.logger.Printf("error writing response to connection: %v", err)
			continue
		}
	}
}

func (s *Server) addReplica(replicaSocketPath string) error {
	var (
		replicaConn net.Conn
		err         error
		retries     int = 10
	)
	for replicaConn == nil && retries >= 0 {
		replicaConn, err = net.Dial(networking.Unix, replicaSocketPath)
		if err != nil {
			s.logger.Printf(
				"server.addReplica: error connecting to replica %s: %v",
				replicaSocketPath,
				err,
			)
			time.Sleep(500 * time.Millisecond)
		}
		retries--
	}

	if replicaConn == nil {
		return fmt.Errorf("error connecting to replica %s", replicaSocketPath)
	}

	s.logger.Printf("connected to replica at %s", replicaConn.RemoteAddr())

	replicaClient, err := client.New(replicaConn)
	if err != nil {
		return fmt.Errorf(
			"server.New: error building replica client for %s: %w",
			replicaSocketPath,
			err,
		)

	}
	s.replicaClients = append(s.replicaClients, replicaClient)
	return nil
}
