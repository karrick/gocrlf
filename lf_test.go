package gocrlf_test

import (
	"bytes"
	"testing"

	"github.com/karrick/gocrlf"
)

func TestLFfromCRLF(t *testing.T) {
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
		t.Run(testCase.output, func(t *testing.T) {
			got := readAllWithChecks(t, &gocrlf.LFfromCRLF{Source: &stagedReader{Stages: testCase.input}})
			if want := []byte(testCase.output); !bytes.Equal(got, want) {
				t.Errorf("Input: %#v; (GOT): %#q; (WNT): %#q", testCase.input, got, want)
			}
		})
	}
}

func TestLFfromCRorCRLF(t *testing.T) {
	testCases := []struct {
		input  []string
		output string
	}{
		{[]string{"\r"}, "\n"},
		{[]string{"\r\n"}, "\n"},
		{[]string{"now is the time\r\n"}, "now is the time\n"},
		{[]string{"now is the time\r\n(trailing data)"}, "now is the time\n(trailing data)"},
		{[]string{"now is the time\n"}, "now is the time\n"},
		{[]string{"now is the time\r"}, "now is the time\n"},     // trailing CR ought to be converted
		{[]string{"\rnow is the time"}, "\nnow is the time"},     // CR not followed by LF ought to be converted
		{[]string{"\rnow is the time\r"}, "\nnow is the time\n"}, // CR not followed by LF ought to be converted

		// no line splits
		{[]string{"first", "second", "third"}, "firstsecondthird"},

		// 1->2 and 2->3 both break across a CRLF
		{[]string{"first\r", "\nsecond\r", "\nthird"}, "first\nsecond\nthird"},

		// 1->2 breaks across CRLF and 2->3 does not
		{[]string{"first\r", "\nsecond", "third"}, "first\nsecondthird"},

		// 1->2 breaks across CRLF and 2 ends in CR but 3 does not begin LF
		{[]string{"first\r", "\nsecond\r", "third"}, "first\nsecond\nthird"},

		// 1 ends in CR but 2 does not begin LF, and 2->3 breaks across CRLF
		{[]string{"first\r", "second\r", "\nthird"}, "first\nsecond\nthird"},

		// 1 ends in CR but 2 does not begin LF, and 2->3 does not break across CRLF
		{[]string{"first\r", "second\r", "\nthird"}, "first\nsecond\nthird"},

		// 1->2 and 2->3 both break across a CRLF, but 3->4 does not
		{[]string{"first\r", "\nsecond\r", "\nthird\r", "fourth"}, "first\nsecond\nthird\nfourth"},
		{[]string{"first\r", "\nsecond\r", "\nthird\n", "fourth"}, "first\nsecond\nthird\nfourth"},

		{[]string{"alpha" + "\n\n" + "bravo" + "\r"},
			"alpha" + "\n\n" + "bravo" + "\n",
		},

		{[]string{"alpha" + "\n\n" + "bravo" + "\r" + "charlie"},
			"alpha" + "\n\n" + "bravo" + "\n" + "charlie",
		},

		{[]string{"alpha" + "\n\n" + "bravo" + "\r\n"},
			"alpha" + "\n\n" + "bravo" + "\n",
		},

		{[]string{"alpha" + "\r\n" + "bravo" + "\r"},
			"alpha" + "\n" + "bravo" + "\n",
		},

		{[]string{"alpha" + "\n" + "bravo" + "\r\n"},
			"alpha" + "\n" + "bravo" + "\n",
		},

		{[]string{"this is the result\r\nfrom the first read\r", "\nthis is the result\r\nfrom the second read\r"},
			"this is the result\nfrom the first read\nthis is the result\nfrom the second read\n"},
		{[]string{"now is the time\r\nfor all good engineers\r\nto improve their test coverage!\r\n"},
			"now is the time\nfor all good engineers\nto improve their test coverage!\n"},
		{[]string{"now is the time\r\nfor all good engineers\r", "\nto improve their test coverage!\r\n"},
			"now is the time\nfor all good engineers\nto improve their test coverage!\n"},

		{[]string{
			"Lorem ipsum RFC 2049 § 2 dolor sit amet, consectetur adipiscing elit,\r\n",
			"sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim\r\n",
			"ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip\r\n",
			"ex ea commodo consequat. Duis aute irure dolor in reprehenderit in\r\n",
			"voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint\r\n",
			"occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit\r\n",
			"anim id est laborum.\r\n",
		},
			`Lorem ipsum RFC 2049 § 2 dolor sit amet, consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim
ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip
ex ea commodo consequat. Duis aute irure dolor in reprehenderit in
voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint
occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit
anim id est laborum.
`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.output, func(t *testing.T) {
			got := readAllWithChecks(t, &gocrlf.LFfromCRorCRLF{R: &stagedReader{Stages: testCase.input}})
			if want := []byte(testCase.output); !bytes.Equal(got, want) {
				t.Errorf("Input: %#v; (GOT): %#q; (WNT): %#q", testCase.input, got, want)
			}
		})
	}
}

func BenchmarkLFfromCRorCRLF(b *testing.B) {
	for i := 0; i < b.N; i++ {
		input := []string{
			"Lorem ipsum RFC 2049 § 2 dolor sit amet, consectetur adipiscing elit,\r\n",
			"sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim\r\n",
			"ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip\r\n",
			"ex ea commodo consequat. Duis aute irure dolor in reprehenderit in\r\n",
			"voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint\r\n",
			"occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit\r\n",
			"anim id est laborum.\r\n",
		}
		output := `Lorem ipsum RFC 2049 § 2 dolor sit amet, consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim
ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip
ex ea commodo consequat. Duis aute irure dolor in reprehenderit in
voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint
occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit
anim id est laborum.
`
		got := readAllWithChecks(b, &gocrlf.LFfromCRorCRLF{R: &stagedReader{Stages: input}})
		if want := []byte(output); !bytes.Equal(got, want) {
			b.Errorf("Input: %#v; (GOT): %#q; (WNT): %#q", input, got, want)
		}
	}
}
