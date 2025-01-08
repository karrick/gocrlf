package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/karrick/gocrlf"
)

var progname = filepath.Base(os.Args[0])

func main() {
	// implement dos2unix like command
	flag.Parse()

	var ior io.Reader
	var err error

	if flag.NArg() == 0 {
		_, err = io.Copy(os.Stdout, &gocrlf.LFfromCRLF{R: ior})
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", progname, err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	for _, arg := range flag.Args() {
		fh, err := os.Open(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", progname, err)
			continue
		}
		_, err = io.Copy(os.Stdout, &gocrlf.LFfromCRLF{R: ior})
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", progname, err)
		}
		if err = fh.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", progname, err)
		}
	}
}
