// Package limage introduces structures to read and build layered images
package limage

import (
	"image"
	"image/color"
)

// Image represents a layered image
type Image struct {
	Group
	Comment      string
	Transparency uint8
}

// At returns the colour at the specified coords
func (i *Image) At(x, y int) color.Color {
	return transparency(i.Group.At(x, y), 255-i.Transparency)
}

// Layer represents a single layer of a multilayered image
type Layer struct {
	Name             string
	OffsetX, OffsetY int
	Mode             Composite
	Invisible        bool
	Transparency     uint8
	image.Image
}

// Bounds returns the limits for the dimensions of the layer
func (l *Layer) Bounds() image.Rectangle {
	max := l.Image.Bounds().Max
	return image.Rect(l.OffsetX, l.OffsetY, max.X+l.OffsetX, max.Y+l.OffsetY)
}

// At returns the colour at the specified coords
func (l *Layer) At(x, y int) color.Color {
	return transparency(l.Image.At(x-l.OffsetX, y-l.OffsetY), 255-l.Transparency)
}

// Group represents a collection of layers
type Group struct {
	image.Config
	Layers []Layer
}

// ColorModel represents the color model of the group
func (g *Group) ColorModel() color.Model {
	return g.Config.ColorModel
}

// Bounds returns the limites for the dimensions of the group
func (g *Group) Bounds() image.Rectangle {
	return image.Rect(0, 0, g.Width, g.Height)
}

// At returns the colour at the specified coords
func (g *Group) At(x, y int) color.Color {
	var c color.Color = color.Alpha{}
	point := image.Point{x, y}
	for i := len(g.Layers) - 1; i >= 0; i-- {
		if g.Layers[i].Invisible {
			continue
		}
		if !point.In(g.Layers[i].Bounds()) {
			continue
		}
		if _, ok := g.Config.ColorModel.(color.Palette); g.Layers[i].Mode != CompositeDissolve && ok {
			if d := colourToNRGBA(g.Layers[i].At(x, y)); d.A > 0x7fff {
				d.A = 0xffff
				c = d
			}
		} else {
			c = g.Layers[i].Mode.Composite(c, g.Layers[i].At(x, y))
		}
	}
	return c
}

// MaskedImage represents an image that has a to-be-applied mask
type MaskedImage struct {
	image.Image
	Mask *image.Gray
}

// At returns the colour at the specified coords after masking
func (m *MaskedImage) At(x, y int) color.Color {
	return transparency(m.Image.At(x, y), m.Mask.GrayAt(x, y).Y)
}

func colourToNRGBA(c color.Color) color.NRGBA64 {
	switch c := c.(type) {
	case color.NRGBA:
		var d color.NRGBA64
		d.R = uint16(c.R)
		d.R |= d.R << 8
		d.G = uint16(c.G)
		d.G |= d.G << 8
		d.B = uint16(c.B)
		d.B |= d.B << 8
		d.A = uint16(c.A)
		d.A |= d.A << 8
		return d
	case color.NRGBA64:
		return c
	}
	if n, ok := c.(interface {
		ToNRGBA() color.NRGBA64
	}); ok {
		return n.ToNRGBA()
	}
	return color.NRGBA64Model.Convert(c).(color.NRGBA64)
}

// Text represents a text layer
type Text struct {
	image.Image
	TextData
}

// TextData represents the stylised text
type TextData []TextDatum

// String returns a flattened string
func (t TextData) String() string {
	toRet := ""
	for _, d := range t {
		toRet += d.Data
	}
	return toRet
}

// TextDatum is a collection of styling for a single piece of text
type TextDatum struct {
	ForeColor, BackColor                   color.Color
	Size, LetterSpacing, Rise              float64
	Bold, Italic, Underline, Strikethrough bool
	Font, Data                             string
	FontUnit                               uint8
}

func transparency(ac color.Color, ao uint8) color.Color {
	if ao == 0xff {
		return ac
	} else if ao == 0 {
		return color.NRGBA64{}
	}
	c := colourToNRGBA(ac)
	o := uint32(ao)
	o |= o << 8
	c.A = uint16(o * uint32(c.A) / 0xffff)
	return c
}
