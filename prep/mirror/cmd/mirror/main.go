package main

import (
	"flag"
	"fmt"
	"log"
	"mirror/pkg"
	"os"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatalf("Usage: ./mirror url")
	}
	err := pkg.Mirror(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "mirror: %v\n", err)
		os.Exit(1)
	}
}
