package main

import (
	"distributed/pkg/networking"
	"distributed/pkg/server"
	"distributed/pkg/storage"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const defaultDbDirectory = "/tmp/distkv"

var dbFile string

func main() {
	flag.StringVar(
		&dbFile,
		"db",
		defaultDbDirectory,
		"directory for database files (primary and replicas)",
	)
	flag.Parse()

	var sigChan = make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go cleanupSocketOnExit(sigChan, networking.SocketPath)
	listener, err := net.Listen(networking.Unix, networking.SocketPath)
	if err != nil {
		log.Fatal("error opening unix socket:", err)
	}
	defer func(listener net.Listener) {
		if err := listener.Close(); err != nil {
			log.Printf("error closing listener: %v", err)
		}
	}(listener)

	log.Println("listening at address", listener.Addr())
	if err != nil {
		log.Fatal("error opening unix socket:", err)
	}
	db, err := storage.NewPersistentStorage(dbFile)
	if err != nil {
		log.Fatalf("error initializing storage: %v", err)
	}
	defer func(db storage.Storage) {
		err := db.Close()
		if err != nil {
			log.Printf("error closing storage: %v\n", err)
		}
	}(db)

	s, err := server.New(listener, db)
	if err != nil {
		log.Fatalf("error initializing server: %v", err)
	}
	s.Run()
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
