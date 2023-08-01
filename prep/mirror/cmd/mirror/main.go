package main

import (
	"flag"
	"fmt"
	"log"
	"mirror/pkg"
	"net/url"
	"os"
)

var outDirName = flag.String("o", "./mirrored", "output directory")

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatalf("Usage: ./mirror url")
	}
	u, err := url.Parse(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "mirror: %v\n", err)
		os.Exit(1)
	}
	// make the directory.  Improvement can be wiping its contents if it exists
	// rather than failing
	if err = os.Mkdir(*outDirName, os.ModePerm); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"error creating directory %q: %v\n",
			*outDirName,
			err,
		)
		os.Exit(1)
	}
	data := pkg.NewMirrorData(*outDirName)
	err = pkg.Mirror(u, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mirror: %v\n", err)
		os.Exit(1)
	}
}
