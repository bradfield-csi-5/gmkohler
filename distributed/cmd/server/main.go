package main

import (
	"distributed/pkg"
	"distributed/pkg/networking"
	"distributed/pkg/server"
	storage2 "distributed/pkg/server/storage"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const defaultDbDirectory = "/tmp/distkv"

var dbDirectory string
var isReplica bool

type roleBasedConfig struct {
	dbFileName         string
	socketPath         string
	replicaSocketPaths []string
}

func main() {
	flag.StringVar(
		&dbDirectory,
		"db",
		defaultDbDirectory,
		"directory for database files (primary and replicas)",
	)
	flag.BoolVar(
		&isReplica,
		"r",
		false,
		"specify whether server should run as replica (affects which socket it connects to)",
	)
	flag.Parse()

	var role pkg.Role
	if isReplica {
		role = pkg.RoleReplica
	} else {
		role = pkg.RolePrimary
	}
	roleBasedCfg, err := getRoleBasedConfig(role)
	if err != nil {
		log.Fatalf("error getting config: %v", err)
	}

	listener, err := net.Listen(networking.Unix, roleBasedCfg.socketPath)
	if err != nil {
		log.Fatal("error opening unix socket:", err)
	}
	defer func(listener net.Listener) {
		if err := listener.Close(); err != nil {
			log.Printf("error closing listener: %v", err)
		}
	}(listener)
	log.Println("listening at address", listener.Addr())
	// cleanup based stuff now that we know which socket to remove
	var sigChan = make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go cleanupSocketOnExit(sigChan, roleBasedCfg.socketPath)

	db, err := storage2.NewPersistentStorage(dbDirectory, roleBasedCfg.dbFileName)
	if err != nil {
		log.Fatalf("error initializing storage: %v", err)
	}
	defer func(db storage2.Storage) {
		err := db.Close()
		if err != nil {
			log.Printf("error closing storage: %v\n", err)
		}
	}(db)

	s, err := server.New(listener, db, roleBasedCfg.replicaSocketPaths...)
	if err != nil {
		log.Fatalf("error initializing server: %v", err)
	}
	s.Run()
}

func getRoleBasedConfig(role pkg.Role) (*roleBasedConfig, error) {
	var (
		replicaSocketPaths []string
		dbFileName         string
		socketPath         string
		err                error
	)
	socketPath, err = networking.SocketPath(role)
	if err != nil {
		return nil, fmt.Errorf("error getting socket path: %w", err)
	}
	dbFileName, err = storage2.FileName(role)
	if err != nil {
		return nil, fmt.Errorf("error getting db file path: %w", err)
	}

	if isReplica {
		replicaSocketPaths = make([]string, 0)
	} else {
		replicaSocketPath, err := networking.SocketPath(pkg.RoleReplica)
		if err != nil {
			return nil, fmt.Errorf("error getting socket path for replica: %w", err)
		}
		replicaSocketPaths = []string{replicaSocketPath}
	}

	return &roleBasedConfig{
		dbFileName:         dbFileName,
		socketPath:         socketPath,
		replicaSocketPaths: replicaSocketPaths,
	}, nil
}

func cleanupSocketOnExit(sigChan <-chan os.Signal, socketPath string) {
	<-sigChan
	log.Println("exiting")
	if err := os.Remove(socketPath); err != nil {
		log.Printf("error closing socket %s: %v", socketPath, err)
	} else {
		log.Printf("socket %s closed\n", socketPath)
	}
	os.Exit(1)
}
