package main

import (
	"distributed/pkg/networking"
	"distributed/pkg/server"
	"distributed/pkg/server/storage"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

const (
	flagDbDirectory             = "db"
	flagOrchestrationSocketPath = "r"
)

var (
	dbDirectory             string
	orchestrationSocketPath string
	serverLogger            *log.Logger
)

func init() {
	const (
		defaultDbDirectory             = "/tmp/distkv"
		defaultOrchestrationSocketPath = ""
	)
	flag.StringVar(
		&dbDirectory,
		flagDbDirectory,
		defaultDbDirectory,
		"directory for database files (primary and replicas)",
	)
	flag.StringVar(
		&orchestrationSocketPath,
		flagOrchestrationSocketPath,
		defaultOrchestrationSocketPath,
		"if included, server listens here for accepting new replicas (effectively making it a primary)",
	)
}

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatalf(
			"usage: %s [-%s database_directory] [-%s replication_socket_path] name socket_path",
			os.Args[0],
			flagDbDirectory,
			flagOrchestrationSocketPath,
		)
	}

	var (
		serverName          string = flag.Arg(0)
		commandSocketPath   string = flag.Arg(1)
		replicationListener net.Listener
	)

	serverLogger = log.New(
		os.Stderr,
		fmt.Sprintf("[server/%s](pid=%d) ", serverName, os.Getpid()),
		log.LstdFlags|log.Lmsgprefix,
	)

	commandSocketFullPath := filepath.Join(dbDirectory, commandSocketPath)
	commandSocket, err := net.Listen(networking.Unix, commandSocketFullPath)
	if err != nil {
		serverLogger.Fatal("error opening unix socket:", err)
	}
	defer func(socket net.Listener) {
		if err := socket.Close(); err != nil {
			serverLogger.Printf("error closing commandSocket: %v", err)
		}
	}(commandSocket)
	serverLogger.Println("listening for commands at address", commandSocket.Addr())

	// cleanup based stuff now that we know which socket to remove
	var socketsToRemove = []string{commandSocketFullPath}
	serverLogger.Printf("orchestration socket: %q", orchestrationSocketPath)
	if len(orchestrationSocketPath) > 0 {
		fullReplicationSocket := filepath.Join(dbDirectory, orchestrationSocketPath)
		replicationListener, err = net.Listen(networking.Unix, fullReplicationSocket)

		if err != nil {
			serverLogger.Fatalf("error connecting to replication socket: %v", err)
		}
		socketsToRemove = append(socketsToRemove, fullReplicationSocket)
	}

	var sigChan = make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go cleanupSocketsOnExit(sigChan, socketsToRemove...)

	db, err := storage.NewPersistentStorage(dbDirectory, serverName)
	if err != nil {
		serverLogger.Fatalf("error initializing storage: %v", err)
	}
	defer func(db storage.Storage) {
		err := db.Close()
		if err != nil {
			serverLogger.Printf("error closing storage: %v\n", err)
		}
	}(db)

	s, err := server.New(db, commandSocket, replicationListener, serverLogger)
	if err != nil {
		serverLogger.Fatalf("error initializing server: %v", err)
	}
	s.Run()
}

func cleanupSocketsOnExit(sigChan <-chan os.Signal, socketPaths ...string) {
	<-sigChan
	serverLogger.Println("exiting")
	for _, socketPath := range socketPaths {
		if err := os.Remove(socketPath); err != nil {
			serverLogger.Printf("error closing socket %s: %v", socketPath, err)
		} else {
			serverLogger.Printf("socket %s closed\n", socketPath)
		}
	}
	os.Exit(1)
}
