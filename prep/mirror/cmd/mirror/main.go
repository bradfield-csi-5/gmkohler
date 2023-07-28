package main

import (
	"flag"
	"fmt"
	"log"
	"mirror/pkg"
	"os"
)

var outDir = flag.String("o", "./mirrored", "output directory")

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatalf("Usage: ./mirror url")
	}
	fmt.Println("outputDir", *outDir)
	err := pkg.Mirror(flag.Arg(0), *outDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mirror: %v\n", err)
		os.Exit(1)
	}
}
