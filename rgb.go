package limage

import (
	"image"
	"image/color"

	"vimagination.zapto.org/limage/lcolor"
)

// RGB is an image of RGB colours.
type RGB struct {
	Pix    []lcolor.RGB
	Stride int
	Rect   image.Rectangle
}

// NewRGB create a new RGB image with the given bounds.
func NewRGB(r image.Rectangle) *RGB {
	w, h := r.Dx(), r.Dy()

	return &RGB{
		Pix:    make([]lcolor.RGB, w*h),
		Stride: w,
		Rect:   r,
	}
}

// At returns the colour at the given coords.
func (r *RGB) At(x, y int) color.Color {
	return r.RGBAt(x, y)
}

// Bounds returns the limits of the image.
func (r *RGB) Bounds() image.Rectangle {
	return r.Rect
}

// ColorModel returns a colour model that converts arbitrary colours to the RGB
// space.
func (r *RGB) ColorModel() color.Model {
	return lcolor.RGBModel
}

// RGBAt returns the exact RGB colour at the given coords.
func (r *RGB) RGBAt(x, y int) lcolor.RGB {
	if !(image.Point{x, y}.In(r.Rect)) {
		return lcolor.RGB{}
	}

	return r.Pix[r.PixOffset(x, y)]
}

// Opaque just returns true as the alpha channel is fixed.
func (r *RGB) Opaque() bool {
	return true
}

// PixOffset returns the index of the Pix array correspinding to the given
// coords.
func (r *RGB) PixOffset(x, y int) int {
	return (y-r.Rect.Min.Y)*r.Stride + x - r.Rect.Min.X
}

// Set converts the given colour to the RGB space and sets it at the given
// coords.
func (r *RGB) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(r.Rect)) {
		return
	}

	rr, g, b, _ := c.RGBA()
	r.Pix[r.PixOffset(x, y)] = lcolor.RGB{
		R: uint8(rr >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
	}
}

// SetRGB directly set an RGB colour to the given coords.
func (r *RGB) SetRGB(x, y int, rgb lcolor.RGB) {
	if !(image.Point{x, y}.In(r.Rect)) {
		return
	}

	r.Pix[r.PixOffset(x, y)] = rgb
}

// SubImage retuns the Image viewable through the given bounds.
func (r *RGB) SubImage(rt image.Rectangle) image.Image {
	rt = rt.Intersect(r.Rect)

	if rt.Empty() {
		return &RGB{}
	}

	return &RGB{
		Pix:    r.Pix[r.PixOffset(rt.Min.X, rt.Min.Y):],
		Stride: r.Stride,
		Rect:   rt,
	}
}
