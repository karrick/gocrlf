package gocrlf

import (
	"bytes"
	"io"
)

var crlf = []byte("\r\n")

// LFfromCRLF is an io.Reader whose Read method converts CRLF byte sequences
// to LF bytes.
type LFfromCRLF struct {
	R          io.Reader
	isPrevByteCR bool
}

// Read consumes bytes from the structure's source io.Reader to fill the
// specified slice of bytes. It converts all CRLF byte sequences to LF, and
// handles cases where CR and LF straddle across two Read operations.
func (f *LFfromCRLF) Read(buf []byte) (int, error) {
	buflen := len(buf)
	if f.isPrevByteCR {
		// Read one fewer bytes so we have room if the first byte of the
		// upcoming Read is not a LF, in which case we will need to insert
		// trailing CR from previous read.
		buflen--
	}
	nr, er := f.R.Read(buf[:buflen])
	if nr > 0 {
		if f.isPrevByteCR && buf[0] != '\n' {
			// Having a CRLF split across two Read operations is rare, so the
			// performance impact of copying entire buffer to the right by one
			// byte, while suboptimal, will at least will not happen very
			// often. This negative performance impact is mitigated somewhat on
			// many Go compilation architectures, GOARCH, because the `copy`
			// builtin uses a machine opcode for performing the memory copy on
			// possibly overlapping regions of memory. This machine opcodes is
			// not instantaneous and does require multiple CPU cycles to
			// complete, but is significantly faster than the application
			// looping through bytes.
			copy(buf[1:nr+1], buf[:nr]) // shift data to right one byte
			buf[0] = '\r'               // insert CR at first byte
			nr++						// pretend we read one additional byte
		}

		// Remove any CRLF byte sequences in the buffer using `bytes.Index`
		// because, like the `copy` builtin on many GOARCHs, it also takes
		// advantage of a machine opcode to search for byte patterns.

		// searchOffset is index within buffer from whence the search will
		// commence for each loop; set to the index of the end of the previous
		// loop.
		var searchOffset int

		// shiftCount is each subsequenct shift operation needs to shift bytes
		// to the left by one more position than the shift that preceded it.
		var shiftCount int

		// previousIndex is the index of previously found CRLF; -1 means no previous index
		previousIndex := -1

		for {
			index := bytes.Index(buf[searchOffset:nr], crlf)
			if index == -1 {
				break
			}
			index += searchOffset // convert relative index to absolute
			if previousIndex != -1 {
				// shift substring between previous index and this index
				copy(buf[previousIndex-shiftCount:], buf[previousIndex+1:index])
				shiftCount++ // next shift needs to be 1 byte to the left
			}
			previousIndex = index
			searchOffset = index + 2 // start next search after len(crlf)
		}
		if previousIndex != -1 {
			// handle final shift
			copy(buf[previousIndex-shiftCount:], buf[previousIndex+1:nr])
			shiftCount++
		}
		nr -= shiftCount // shorten byte read count by number of shifts executed

		// When final byte from a read operation is CR, do not emit it until
		// ensure first byte on next read is not LF.
		if f.isPrevByteCR = buf[nr-1] == '\r'; f.isPrevByteCR {
			nr-- // pretend byte was never read from source
		}
	} else if f.isPrevByteCR {
		// Reading from source returned nothing, but this struct is sitting on
		// a trailing CR byte from previous Read, so let's give the CR to
		// client.
		buf[0] = '\r'
		nr = 1
		er = nil
		f.isPrevByteCR = false // prevent infinite loop
	}
	return nr, er
}

// LFfromCRorCRLF is an io.Reader whose Read method converts bare CR bytes or
// CRLF byte sequences to LF bytes.
type LFfromCRorCRLF struct {
	R            io.Reader
	isPrevByteCR bool
}

// Read consumes bytes from the structure's source io.Reader to fill the
// specified slice of bytes. It converts all bare CR bytes and CRLF byte
// sequences to LF, and handles cases where CR and LF straddle across two Read
// operations.
func (f *LFfromCRorCRLF) Read(buf []byte) (int, error) {
	buflen := len(buf)
	if f.isPrevByteCR {
		// Read one fewer bytes so we have room if the first byte of the
		// upcoming Read is not a LF, in which case we will need to insert
		// LF for the trailing CR from previous read.
		buflen--
	}
	nr, er := f.R.Read(buf[:buflen])
	if nr > 0 {
		// index is index within buffer from whence the search will commence
		// for next CR.
		var index int

		if f.isPrevByteCR {
			f.isPrevByteCR = false
			if buf[0] != '\n' {
				// Having a CRLF split across two Read operations is rare, so
				// the performance impact of copying entire buffer to the
				// right by one byte, while suboptimal, will at least will not
				// happen very often. This negative performance impact is
				// mitigated somewhat on many Go compilation architectures,
				// GOARCH, because the `copy` builtin uses a machine opcode
				// for performing the memory copy on possibly overlapping
				// regions of memory. This machine opcodes is not
				// instantaneous and requires multiple CPU cycles to complete,
				// but is significantly faster than the application looping
				// through bytes.
				copy(buf[1:nr+1], buf[:nr]) // shift data to right one byte
				buf[0] = '\n'				// insert LF at first byte
				nr++						// pretend we read one additional byte
				index++ // optimization
			}
		}

		// Remove any CR and CRLF sequences in the buffer using
		// `bytes.IndexByte` because, like the `copy` builtin on many GOARCHs,
		// it takes advantage of machine opcodes to search for bytes.
		for {
			relativeIndex := bytes.IndexByte(buf[index:nr], '\r')
			if relativeIndex == -1 {
				break
			}
			index += relativeIndex // convert relative to absolute index

			indexPlusOne := index + 1
			f.isPrevByteCR = indexPlusOne == nr
			if f.isPrevByteCR {
				// This is the final byte read into the buffer.
				nr--          // ignore final byte of the line (CR)
				return nr, er // optimization over break statement
			} else {
				// There are more bytes after this byte in the buffer.
				if buf[indexPlusOne] == '\n' {
					copy(buf[index:], buf[indexPlusOne:nr])
					nr--
				} else {
					buf[index] = '\n' // Replace lone CR by LF.
				}
			}
			index = indexPlusOne // start search at next byte
		}
	} else if f.isPrevByteCR {
		// Reading from source returned nothing, but this struct is sitting on
		// a trailing CR byte from previous Read, so let's give a LF to
		// client.
		buf[0] = '\n'
		nr = 1
		er = nil
		f.isPrevByteCR = false
	}
	return nr, er
}
