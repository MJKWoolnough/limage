// Package limage introduces structures to read and build layered images
package limage

import (
	"image"
	"image/color"

	"github.com/MJKWoolnough/limage/internal"
)

// Layer represents a single layer of a multilayered image
type Layer struct {
	Name         string
	LayerBounds  image.Rectangle // Bounds within the layer group
	Mode         Composite
	Invisible    bool
	Transparency uint8
	image.Image
}

// Bounds returns the limits for the dimensions of the layer
func (l Layer) Bounds() image.Rectangle {
	return l.LayerBounds
}

// At returns the colour at the specified coords
func (l Layer) At(x, y int) color.Color {
	return transparency(l.Image.At(x-l.LayerBounds.Min.X, y-l.LayerBounds.Min.Y), 255-l.Transparency)
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image
func (l Layer) SubImage(r image.Rectangle) image.Image {
	l.LayerBounds = r.Intersect(l.LayerBounds)
	return l
}

// Image represents a collection of layers
type Image []Layer

// ColorModel represents the color model of the group. It uses the first layer
// to determine the color model
func (g Image) ColorModel() color.Model {
	if len(g) == 0 {
		return color.AlphaModel
	}
	return g[0].ColorModel()
}

// Bounds returns the limits for the dimensions of the group
func (g Image) Bounds() image.Rectangle {
	var r image.Rectangle
	for _, l := range g {
		r = r.Union(l.Bounds())
	}
	return r
}

// At returns the colour at the specified coords
func (g Image) At(x, y int) color.Color {
	c := color.Color(color.Alpha{})
	point := image.Point{x, y}
	for i := len(g) - 1; i >= 0; i-- {
		if g[i].Invisible {
			continue
		}
		if !point.In(g[i].Bounds()) {
			continue
		}
		if _, ok := g.ColorModel().(color.Palette); g[i].Mode != CompositeDissolve && ok {
			if d := internal.ColourToNRGBA(g[i].At(x, y)); d.A > 0x7fff {
				d.A = 0xffff
				c = d
			}
		} else {
			c = g[i].Mode.Composite(c, g[i].At(x, y))
		}
	}
	return c
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image
func (g Image) SubImage(r image.Rectangle) image.Image {
	h := make(Image, len(g))
	for n, l := range g {
		h[n] = l.SubImage(r).(Layer)
	}
	return h
}

// MaskedImage represents an image that has a to-be-applied mask
type MaskedImage struct {
	image.Image
	Mask *image.Gray
}

// At returns the colour at the specified coords after masking
func (m MaskedImage) At(x, y int) color.Color {
	return transparency(m.Image.At(x, y), m.Mask.GrayAt(x, y).Y)
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
