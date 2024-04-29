package main

import (
	"bufio"
	"distributed/storage"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		reader  *bufio.Reader = bufio.NewReader(os.Stdin)
		writer  *bufio.Writer = bufio.NewWriter(os.Stdout)
		sigChan               = make(chan os.Signal)
	)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		if _, err := writer.WriteString("exiting"); err != nil {
			log.Fatalf("failure to write to stdout: %v", err)
		}
		if err := writer.Flush(); err != nil {
			log.Fatalf("failure to flush stdout: %v", err)
		}
		os.Exit(1)
	}()

	var (
		db      storage.Storage
		err     error
		input   string
		command *Command
	)

	if len(os.Args) > 1 {
		fileName := os.Args[1]
		log.Printf("opening test database %q\n", fileName)
		db, err = storage.NewPersistentStorage(fileName)
	} else {
		db, err = storage.NewInMemoryStorage()
	}
	if err != nil {
		log.Fatalf("error intializing storage: %v", err)
	}

	// REPL
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

		command, err = ParseInput(input)
		if err != nil {
			log.Printf("error parsing input: %v", err)
			continue
		}

		var output string
		switch command.Operation {
		case Get:
			value, err := db.Get(command.Key)
			if err != nil {
				output = fmt.Sprintf("error getting value: %v\n", err)
			} else {
				output = fmt.Sprintf("%s\n", value)
			}
		case Put:
			value, err := db.Put(command.Key, command.Value)
			if err != nil {
				output = fmt.Sprintf("error putting value: %v\n", err)
			} else {
				output = fmt.Sprintf("%s\n", value)
			}
		default: // shouldn't happen because ParseInput ensures we know the operation
			output = "unknown error: could not recognize operation\n"
		}

		if _, err = writer.WriteString(output); err != nil {
			log.Fatalf("error writing to stdout: %v", err)
		}
		if err = writer.Flush(); err != nil {
			log.Fatalf("error flushing stdout: %v", err)
		}
	}
}
