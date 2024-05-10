package main

import (
	"distributed/pkg/networking"
	"distributed/pkg/networking/replication"
	"distributed/pkg/server"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

const (
	defaultDbDirectory   = "/tmp/distkv"
	routingSocketPath    = "routing.sock"
	orchestrationSocket  = "orchestration.sock"
	replicaCommandSocket = "replica.sock"
	primaryCommandSocket = "primary.sock"
)

var (
	dbDirectory         string
	launchedProcesses   []*os.Process
	openedSockets       []net.Listener
	orchestrationLogger *log.Logger
)

func main() {
	flag.StringVar(
		&dbDirectory,
		"db",
		defaultDbDirectory,
		"directory for database files (primary and replicaSockets)",
	)
	flag.Parse()
	orchestrationLogger = log.New(os.Stderr, "[orchestration] ", log.LstdFlags|log.Lmsgprefix)

	fPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		orchestrationLogger.Fatalf("error getting filepath: %v", err)
	}
	serverExecutable := filepath.Join(fPath, "../server/server")

	// teardown code
	var (
		endProgramWg sync.WaitGroup
		sigChan      = make(chan os.Signal)
	)
	endProgramWg.Add(1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func(sigChan <-chan os.Signal, wg *sync.WaitGroup) {
		<-sigChan
		for _, process := range launchedProcesses {
			if process != nil {
				if err := process.Kill(); err != nil {
					orchestrationLogger.Printf("error killing process %d: %v", process.Pid, err)
				}
			}
		}
		for _, socket := range openedSockets {
			if err := os.Remove(socket.Addr().String()); err != nil {
				orchestrationLogger.Printf("error closing socket %s: %v", socket, err)
			}
		}
		wg.Done()
		os.Exit(1)
	}(sigChan, &endProgramWg)

	// set up primary server
	var primaryServerArgs = []string{
		"-db", dbDirectory,
		"-r", orchestrationSocket,
		"primary", // server name
		primaryCommandSocket,
	}
	primaryServerCmd, err := spawnServer(
		serverExecutable,
		primaryServerArgs...,
	)
	if err != nil {
		orchestrationLogger.Fatalf("error starting primary server: %v", err)
	}
	launchedProcesses = append(launchedProcesses, primaryServerCmd.Process)

	// notify primary of replicas
	var conn net.Conn
	for {
		orchestrationSocketPath := filepath.Join(dbDirectory, orchestrationSocket)
		orchestrationLogger.Printf("dialling socket %s for orchestration", orchestrationSocketPath)
		conn, err = net.Dial(networking.Unix, orchestrationSocketPath)
		if err != nil {
			orchestrationLogger.Printf("error dialling primary's orchestration socket: %v", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		break
	}
	orchestrationLogger.Println("connected to orchestration socket")

	// TODO: accept `-r N` and range over that number
	var replicaSockets = []string{replicaCommandSocket}
	for _, replicaSocket := range replicaSockets {
		replicaCmd, err := spawnServer(
			serverExecutable,
			"replica",
			replicaSocket,
			fmt.Sprintf("-db=%s", dbDirectory),
		)
		if err != nil {
			orchestrationLogger.Printf("error starting replica for %s: %v", replicaSocket, err)
			continue
		}
		launchedProcesses = append(launchedProcesses, replicaCmd.Process)

		// this is probably a race condition of the replica creating its socket and
		// the server trying to dial it.
		if err := sendAddListenerRequest(conn, filepath.Join(dbDirectory, replicaSocket)); err != nil {
			orchestrationLogger.Printf("error sending request: %v", err)
		}

	}
	routingSocket, err := net.Listen(networking.Unix, filepath.Join(dbDirectory, routingSocketPath))
	if err != nil {
		log.Fatalf("error creating socket for routing: %v", err)
	}
	var servers = make([]*server.ForwardedServer, 2)
	for j, socketPath := range []string{
		primaryCommandSocket, // keep first in line for easy routing
		replicaCommandSocket,
	} {
		conn, err := net.Dial(networking.Unix, filepath.Join(dbDirectory, socketPath))
		if err != nil {
			log.Fatalf("error connecting to server %s: %v", socketPath, err)
		}
		forwardedServer, err := server.NewForwardedServer(conn)
		if err != nil {
			log.Fatalf("error initializing ForwardedServer: %v", err)
		}
		servers[j] = forwardedServer
	}

	router, err := server.NewRouter(orchestrationLogger, routingSocket, servers)
	if err != nil {
		log.Fatalf("error initializing router: %v", err)
	}
	router.Run()

	endProgramWg.Wait()
}

func spawnServer(executable string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(executable, args...)
	cmd.Stderr = os.Stderr
	go func(command *exec.Cmd) {
		if err := cmd.Run(); err != nil {
			orchestrationLogger.Printf("error executing command: %v", err)
		}
		if _, err := cmd.Output(); err != nil {
			orchestrationLogger.Printf("error in command output: %v", err)
		}
	}(cmd)
	return cmd, nil
}

func sendAddListenerRequest(conn net.Conn, listenerSocketPath string) error {
	if err := replication.WriteAddReplicaRequest(
		conn,
		&replication.AddReplicaRequest{SocketPath: listenerSocketPath},
	); err != nil {
		return err
	}

	var responsePtr = new(replication.AddReplicaResponse)
	err := replication.ReadAddReplicaResponse(conn, responsePtr)
	if err != nil {
		return err
	}

	if len(responsePtr.Error) > 0 {
		return errors.New(responsePtr.Error)
	}
	return nil
}
