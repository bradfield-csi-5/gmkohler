package server

import (
	"distributed/pkg/networking"
	"distributed/pkg/networking/commands"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"math/rand"
	"net"
)

func NewRouter(logger *log.Logger, socket net.Listener, servers []*ForwardedServer) (*Router, error) {
	return &Router{
		logger:  logger,
		socket:  socket,
		servers: servers,
	}, nil
}

type Router struct {
	logger  *log.Logger
	socket  net.Listener
	servers []*ForwardedServer
}

func (r *Router) Run() {
	for {
		conn, err := r.socket.Accept()
		if err != nil {
			r.logger.Printf("error accepting connection: %v", err)
			continue
		}
		r.logger.Printf("handling connection from %s", conn.RemoteAddr())
		var requestPtr = new(commands.ExecuteCommandRequest)

		for { // TODO do not actually unpack the requests by encoding the type in the binary and just relay bytes
			err := commands.ReadExecuteCommandRequest(conn, requestPtr)
			if err != nil {
				r.logger.Printf("error reading request: %v", err)
				conn.Close()
				break
			}
			encodedRequest, err := proto.Marshal(requestPtr)
			if err != nil {
				r.logger.Printf("error marshalling proto: %v", err)
			}

			var encodedResponse []byte
			switch req := requestPtr.Operation.(type) {
			case *commands.ExecuteCommandRequest_Get:
				chosenServer := r.servers[rand.Intn(len(r.servers))]
				encodedResponse, err = chosenServer.SendMessage(encodedRequest)
			case *commands.ExecuteCommandRequest_Put:
				encodedResponse, err = r.servers[0].SendMessage(encodedRequest)
			default:
				err = fmt.Errorf("unrecognized operation: %v", req)
			}

			if err != nil {
				r.logger.Printf("error forwarding response to server: %v", err)
			}

			if err := networking.WriteMessageRaw(conn, encodedResponse); err != nil {
				r.logger.Printf("error forwarding request to client: %v", err)
				conn.Close()
				break
			}
		}
	}
}
