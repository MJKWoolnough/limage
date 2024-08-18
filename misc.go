// Package limage introduces structures to read and build layered images.
package limage

import (
	"image"
	"image/color"
	"strings"

	"vimagination.zapto.org/limage/internal"
)

// MaskedImage represents an image that has a to-be-applied mask.
type MaskedImage struct {
	image.Image
	Mask *image.Gray
}

// At returns the colour at the specified coords after masking.
func (m MaskedImage) At(x, y int) color.Color {
	return transparency(m.Image.At(x, y), m.Mask.GrayAt(x, y).Y)
}

// Text represents a text layer.
type Text struct {
	image.Image
	TextData
}

// TextData represents the stylised text.
type TextData []TextDatum

// String returns a flattened string.
func (t TextData) String() string {
	var sb strings.Builder

	for _, d := range t {
		sb.WriteString(d.Data)
	}

	return sb.String()
}

// TextDatum is a collection of styling for a single piece of text.
type TextDatum struct {
	ForeColor, BackColor                   color.Color
	Size, LetterSpacing, Rise              uint32
	Bold, Italic, Underline, Strikethrough bool
	Font, Data                             string
}

func transparency(ac color.Color, ao uint8) color.Color {
	if ao == 0xff {
		return ac
	} else if ao == 0 {
		return color.NRGBA64{}
	}

	c := internal.ColourToNRGBA(ac)
	o := uint32(ao)
	o |= o << 8
	c.A = uint16(o * uint32(c.A) / 0xffff)

	return c
}
