package limage // import "vimagination.zapto.org/limage"

import (
	"image"
	"image/color"

	"vimagination.zapto.org/limage/internal"
)

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
		h[n] = l.SubImageLayer(r)
	}
	return h
}
