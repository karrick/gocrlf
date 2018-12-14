package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/karrick/gocrlf"
	"github.com/karrick/gorill"
)

func main() {
	// implement dos2unix like command

	var ior io.Reader
	if flag.NArg() == 0 {
		ior = os.Stdin
	} else {
		ior = &gorill.FilesReader{Pathnames: flag.Args()}
	}

	_, err := io.Copy(os.Stdout, &gocrlf.LFfromCRLF{Source: ior})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
