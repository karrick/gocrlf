package gocrlf_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/karrick/gocrlf"
)

func TestLFfromCRLF(t *testing.T) {
	stream := func(t *testing.T, iterations []string) []byte {
		dst := new(bytes.Buffer)
		n, err := io.Copy(dst, &gocrlf.LFfromCRLF{Source: &crossBuffer{iterations: iterations}})
		if got, want := err, error(nil); got != want {
			t.Errorf("(GOT): %v; (WNT): %v", got, want)
		}
		if got, want := n, int64(dst.Len()); got != want {
			t.Errorf("(GOT): %v; (WNT): %v", got, want)
		}
		return dst.Bytes()
	}

	testCases := []struct {
		input  []string
		output string
	}{
		{[]string{"\r"}, "\r"},
		{[]string{"\r\n"}, "\n"},
		{[]string{"now is the time\r\n"}, "now is the time\n"},
		{[]string{"now is the time\r\n(trailing data)"}, "now is the time\n(trailing data)"},
		{[]string{"now is the time\n"}, "now is the time\n"},
		{[]string{"now is the time\r"}, "now is the time\r"},     // trailing CR ought to convey
		{[]string{"\rnow is the time"}, "\rnow is the time"},     // CR not followed by LF ought to convey
		{[]string{"\rnow is the time\r"}, "\rnow is the time\r"}, // CR not followed by LF ought to convey

		// no line splits
		{[]string{"first", "second", "third"}, "firstsecondthird"},

		// 1->2 and 2->3 both break across a CRLF
		{[]string{"first\r", "\nsecond\r", "\nthird"}, "first\nsecond\nthird"},

		// 1->2 breaks across CRLF and 2->3 does not
		{[]string{"first\r", "\nsecond", "third"}, "first\nsecondthird"},

		// 1->2 breaks across CRLF and 2 ends in CR but 3 does not begin LF
		{[]string{"first\r", "\nsecond\r", "third"}, "first\nsecond\rthird"},

		// 1 ends in CR but 2 does not begin LF, and 2->3 breaks across CRLF
		{[]string{"first\r", "second\r", "\nthird"}, "first\rsecond\nthird"},

		// 1 ends in CR but 2 does not begin LF, and 2->3 does not break across CRLF
		{[]string{"first\r", "second\r", "\nthird"}, "first\rsecond\nthird"},

		// 1->2 and 2->3 both break across a CRLF, but 3->4 does not
		{[]string{"first\r", "\nsecond\r", "\nthird\r", "fourth"}, "first\nsecond\nthird\rfourth"},
		{[]string{"first\r", "\nsecond\r", "\nthird\n", "fourth"}, "first\nsecond\nthird\nfourth"},

		{[]string{"this is the result\r\nfrom the first read\r", "\nthis is the result\r\nfrom the second read\r"},
			"this is the result\nfrom the first read\nthis is the result\nfrom the second read\r"},
		{[]string{"now is the time\r\nfor all good engineers\r\nto improve their test coverage!\r\n"},
			"now is the time\nfor all good engineers\nto improve their test coverage!\n"},
		{[]string{"now is the time\r\nfor all good engineers\r", "\nto improve their test coverage!\r\n"},
			"now is the time\nfor all good engineers\nto improve their test coverage!\n"},
	}

	for _, testCase := range testCases {
		got := stream(t, testCase.input)
		if want := []byte(testCase.output); !bytes.Equal(got, want) {
			t.Errorf("Input: %#v; (GOT): %#q; (WNT): %#q", testCase.input, got, want)
		}
	}
}
