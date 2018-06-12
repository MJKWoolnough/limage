package limage

import (
	"image"
	"image/color"

	"vimagination.zapto.org/limage/lcolor"
)

// GrayAlpha is an image of GrayAlpha pixels
type GrayAlpha struct {
	Pix    []lcolor.GrayAlpha
	Stride int
	Rect   image.Rectangle
}

// NewGrayAlpha create a new GrayAlpha image with the given bounds
func NewGrayAlpha(r image.Rectangle) *GrayAlpha {
	w, h := r.Dx(), r.Dy()
	return &GrayAlpha{
		Pix:    make([]lcolor.GrayAlpha, w*h),
		Stride: w,
		Rect:   r,
	}
}

// At returns the color for the pixel at the specified coords
func (g *GrayAlpha) At(x, y int) color.Color {
	return g.GrayAlphaAt(x, y)
}

// Bounds returns the limits of the image
func (g *GrayAlpha) Bounds() image.Rectangle {
	return g.Rect
}

// ColorModel returns a color model to transform arbitrary colours into a
// GrayAlpha color
func (g *GrayAlpha) ColorModel() color.Model {
	return lcolor.GrayAlphaModel
}

// GrayAlphaAt returns a GrayAlpha colr for the specified coords
func (g *GrayAlpha) GrayAlphaAt(x, y int) lcolor.GrayAlpha {
	if !(image.Point{x, y}.In(g.Rect)) {
		return lcolor.GrayAlpha{}
	}
	return g.Pix[g.PixOffset(x, y)]
}

// Opaque returns true if all pixels have full alpha
func (g *GrayAlpha) Opaque() bool {
	for _, c := range g.Pix {
		if c.A != 255 {
			return false
		}
	}
	return true
}

// PixOffset returns the index of the element of Pix corresponding to the given
// coords
func (g *GrayAlpha) PixOffset(x, y int) int {
	return (y-g.Rect.Min.Y)*g.Stride + x - g.Rect.Min.X
}

// Set converts the given colour to a GrayAlpha colour and sets it at the given
// coords
func (g *GrayAlpha) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(g.Rect)) {
		return
	}
	g.Pix[g.PixOffset(x, y)] = lcolor.GrayAlphaModel.Convert(c).(lcolor.GrayAlpha)
}

// SetGrayAlpha sets the colour at the given coords
func (g *GrayAlpha) SetGrayAlpha(x, y int, ga lcolor.GrayAlpha) {
	if !(image.Point{x, y}.In(g.Rect)) {
		return
	}
	g.Pix[g.PixOffset(x, y)] = ga
}

// SubImage retuns the Image viewable through the given bounds
func (g *GrayAlpha) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(g.Rect)
	if r.Empty() {
		return &GrayAlpha{}
	}
	return &GrayAlpha{
		Pix:    g.Pix[g.PixOffset(r.Min.X, r.Min.Y):],
		Stride: g.Stride,
		Rect:   r,
	}
}
