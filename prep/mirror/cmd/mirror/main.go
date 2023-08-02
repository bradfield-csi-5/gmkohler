package main

import (
	"flag"
	"fmt"
	"log"
	"mirror/pkg"
	"net/url"
	"os"
	"path/filepath"
	"sync"
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
	fp, err := filepath.Abs(*outDirName)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"mirror: error determining absolute filepath for %q: %v\n",
			*outDirName,
			err,
		)
		os.Exit(1)
	}
	if err = os.Mkdir(fp, os.ModePerm); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"error creating directory %q: %v\n",
			fp,
			err,
		)
		os.Exit(1)
	}
	var wg sync.WaitGroup
	data := pkg.NewMirrorData(fp)

	wg.Add(1)
	go func() {
		e := pkg.Mirror(u, data, &wg)

		if e != nil {
			fmt.Fprintf(os.Stderr, "error mirroring page: %v\n", err)
		}
	}()
	wg.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "mirror: %v\n", err)
		os.Exit(1)
	}
}
