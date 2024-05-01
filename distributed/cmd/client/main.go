package main

import (
	"bufio"
	"distributed/pkg/client"
	"distributed/pkg/networking"
	"distributed/pkg/storage"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// cmd/client is a REPL for interacting with a cmd/server process
func main() {
	var (
		reader  = bufio.NewReader(os.Stdin)
		writer  = bufio.NewWriter(os.Stdout)
		sigChan = make(chan os.Signal)
		err     error
	)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func(sigChan <-chan os.Signal) {
		<-sigChan
		log.Println("exiting")
		os.Exit(1)
	}(sigChan)

	conn, err := net.Dial(networking.Unix, networking.SocketPath)
	if err != nil {
		log.Fatalf("error establishing connection: %v", err)
	}
	dbClient, err := client.New(conn)
	if err != nil {
		log.Fatalf("error initializing client: %v", err)
	}

	// REPL
	var (
		input   string
		command *networking.Command
		output  string
		value   storage.Value
	)
	for {
		if _, err = writer.WriteString("> "); err != nil {
			log.Fatalf("error writing to stdout: %v", err)
		}
		if err = writer.Flush(); err != nil {
			log.Fatalf("error flushing to stdout: %v", err)
		}

		input, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error reading from stdin: %v", err)
		}

		command, err = Input(input)
		if err != nil {
			log.Printf("error parsing input: %v", err)
			continue
		}

		value, err = dbClient.ExecuteCommand(*command)
		if err != nil {
			output = fmt.Sprintf("error executing command: %v\n", err)
		} else {
			output = fmt.Sprintf("%s\n", value)
		}

		if _, err = writer.WriteString(output); err != nil {
			log.Fatalf("error writing to stdout: %v", err)
		}
		if err = writer.Flush(); err != nil {
			log.Fatalf("error flushing stdout: %v", err)
		}
	}
}
