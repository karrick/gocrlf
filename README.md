# gocrlf

io.Reader that converts CRLF sequence to LF for Go 

## Description

LineEndingReader is a `io.Reader` that converts CRLF sequences to LF.

This structure wraps an io.Reader that modifies the file's contents
when it is read, translating all CRLF sequences to LF.

## Usage

```Go
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
```
