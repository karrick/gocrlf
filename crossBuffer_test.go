package gocrlf_test

import "io"

// crossBuffer is a test io.Reader that emits a few canned responses.
type crossBuffer struct {
	readCount  int
	iterations []string
}

func (cb *crossBuffer) Read(buf []byte) (int, error) {
	if cb.readCount == len(cb.iterations) {
		return 0, io.EOF
	}
	cb.readCount++
	return copy(buf, cb.iterations[cb.readCount-1]), nil
}
