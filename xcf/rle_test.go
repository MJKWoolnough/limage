package xcf

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/MJKWoolnough/byteio"
)

func TestReads(t *testing.T) {
	tests := []struct {
		Input, Output string
	}{
		{},
		{
			"\x00A",
			"A",
		},
		{
			"\x01A",
			"AA",
		},
		{
			"\x7eA",
			"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		},
		{
			"\x7f\x00\x00A",
			"",
		},
		{
			"\x7f\x00\x01A",
			"A",
		},
		{
			"\x7f\x00\x0aA",
			"AAAAAAAAAA",
		},
		{
			"\x7f\x01\x00A",
			"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		},
		{
			"\x80\x00\x01A",
			"A",
		},
		{
			"\x80\x00\x02AA",
			"AA",
		},
		{
			"\xffA",
			"A",
		},
		{
			"\xfeAB",
			"AB",
		},
		{
			"\x00A\x01B\x7f\x00\x01C\x7f\x00\x0aD\x80\x00\x0a1234567890",
			"ABBCDDDDDDDDDD1234567890",
		},
	}
	for n, test := range tests {
		r := rle{
			Reader: byteio.StickyReader{
				Reader: byteio.BigEndianReader{Reader: strings.NewReader(test.Input)},
			},
		}
		var buf bytes.Buffer
		io.Copy(&buf, &r)
		if str := buf.String(); str != test.Output {
			t.Errorf("test %d: expecting %q, got %q", n+1, test.Output, str)
		}
	}
}
