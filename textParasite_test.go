package xcf

import (
	"image/color"
	"testing"
)

func (t Text) Equal(u Text) bool {
	if u.ForeColor == nil {
		u.ForeColor = color.Gray{}
	}
	if u.BackColor == nil {
		u.BackColor = color.Alpha{}
	}
	r1, g1, b1, a1 := t.BackColor.RGBA()
	r2, g2, b2, a2 := u.BackColor.RGBA()
	r3, g3, b3, a3 := t.ForeColor.RGBA()
	r4, g4, b4, a4 := u.ForeColor.RGBA()
	return t.Bold == u.Bold &&
		t.Data == u.Data &&
		t.Font == u.Font &&
		t.FontUnit == u.FontUnit &&
		t.Italic == u.Italic &&
		t.LetterSpacing == u.LetterSpacing &&
		t.Rise == u.Rise &&
		t.Size == u.Size &&
		t.Strikethrough == u.Strikethrough &&
		t.Underline == u.Underline &&
		r1 == r2 &&
		r3 == r4 &&
		g1 == g2 &&
		g3 == g4 &&
		b1 == b2 &&
		b3 == b4 &&
		a1 == a2 &&
		a3 == a4
}

func TestParseTextParasite(t *testing.T) {
	tests := []struct {
		input  string
		output TextData
		err    error
	}{
		{
			"(markup \"<markup>Hello, World</markup>\")",
			TextData{
				{Data: "Hello, World"},
			},
			nil,
		},
		{
			"(markup \"<markup>Hello, <b>World</b></markup>\")",
			TextData{
				{Data: "Hello, "},
				{Data: "World", Bold: true},
			},
			nil,
		},
		{
			"(markup \"<markup><i>H<b>e</b>l<u>l</u></i>o, <b>World</b></markup>\")\n(hinting true)",
			TextData{
				{Data: "H", Italic: true},
				{Data: "e", Italic: true, Bold: true},
				{Data: "l", Italic: true},
				{Data: "l", Italic: true, Underline: true},
				{Data: "o, "},
				{Data: "World", Bold: true},
			},
			nil,
		},
		{
			"(markup \"<markup>Hello, World</markup>\")\n(color (color-rgb 255 0 0))",
			TextData{
				{Data: "Hello, World", ForeColor: color.RGBA{255, 0, 0, 255}},
			},
			nil,
		},
	}

	for n, test := range tests {
		o, err := parseTextParasite(test.input)
		if test.err != nil {
			if test.err != err {
				t.Errorf("test %d: expecting error %q, got %q", n+1, test.err, err)
			}
		} else if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
		} else {
			if len(o) != len(test.output) {
				t.Errorf("test %d: expecting length %d, got %d", len(test.output), len(o))
				continue
			}
			for m := range o {
				if !o[m].Equal(test.output[m]) {
					t.Errorf("test %d: expecting %v\ngot: %v", n+1, []Text(test.output), []Text(o))
					break
				}
			}
		}

	}
}
