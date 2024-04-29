package main

import (
	"bufio"
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
		err     error
		input   string
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
		if _, err = writer.WriteString(input); err != nil {
			log.Fatalf("error writing to stdout: %v", err)
		}
		if err = writer.Flush(); err != nil {
			log.Fatalf("error flushing stdout: %v", err)
		}
	}
}
