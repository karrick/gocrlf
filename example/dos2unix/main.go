package main

import (
	"fmt"
	"io"
	"os"

	"github.com/karrick/gocrlf"
)

func main() {
	// implement dos2unix like command
	_, err := io.Copy(os.Stdout, &gocrlf.LFfromCRLF{Source: os.Stdin})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
