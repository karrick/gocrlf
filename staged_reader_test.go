package gocrlf_test

import (
	"bytes"
	"io"
	"testing"
)

// stagedReader is an io.Reader that on each Read invocation returns the next
// in a series of canned responses.
type stagedReader struct {
	Stages    []string
	readCount int
}

func (cb *stagedReader) Read(buf []byte) (int, error) {
	if cb.readCount == len(cb.Stages) {
		return 0, io.EOF
	}
	nr := copy(buf, cb.Stages[cb.readCount])
	cb.readCount++
	return nr, nil
}

func readAllWithChecks(tb testing.TB, r io.Reader) []byte {
	tb.Helper()
	dst := new(bytes.Buffer)
	n, err := io.Copy(dst, r)
	if got, want := err, error(nil); got != want {
		tb.Errorf("(GOT): %v; (WNT): %v", got, want)
	}
	if got, want := n, int64(dst.Len()); got != want {
		tb.Errorf("(GOT): %v; (WNT): %v", got, want)
	}
	return dst.Bytes()
}
