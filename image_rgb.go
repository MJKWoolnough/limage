package xcf

import (
	"image"
	"image/color"
)

// RGB is a standard colour type whose Alpha channel is always full
type RGB struct {
	R, G, B uint8
}

// RGBA implements the color.Color interface
func (rgb RGB) RGBA() (r, g, b, a uint32) {
	r = uint32(rgb.R)
	r |= r << 8
	g = uint32(rgb.G)
	g |= g << 8
	b = uint32(rgb.B)
	b |= b << 8
	return r, g, b, 0xFFFF
}

// ToNRGBA returns itself as a non-alpha-premultiplied value
// As the alpha is always full, this only returns the normal values
func (rgb RGB) ToNRGBA() color.NRGBA64 {
	r := uint16(rgb.R)
	r |= r << 8
	g := uint16(rgb.G)
	g |= g << 8
	b := uint16(rgb.B)
	b |= b << 8
	return color.NRGBA64{r, g, b, 0xffff}
}

func rgbColourModel(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	return RGB{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
	}
}

// RGBImage is an image of RGB colours
type RGBImage struct {
	Pix    []RGB
	Stride int
	Rect   image.Rectangle
}

func newRGB(r image.Rectangle) *RGBImage {
	w, h := r.Dx(), r.Dy()
	return &RGBImage{
		Pix:    make([]RGB, w*h),
		Stride: w,
		Rect:   r,
	}
}

// At returns the colour at the given coords
func (r *RGBImage) At(x, y int) color.Color {
	return r.RGBAt(x, y)
}

// Bounds returns the limits of the image
func (r *RGBImage) Bounds() image.Rectangle {
	return r.Rect
}

// ColorModel returns a colour model that converts arbitrary colours to the RGB
// space
func (r *RGBImage) ColorModel() color.Model {
	return color.ModelFunc(rgbColourModel)
}

// RGBAt returns the exact RGB colour at the given coords
func (r *RGBImage) RGBAt(x, y int) RGB {
	if !(image.Point{x, y}.In(r.Rect)) {
		return RGB{}
	}
	return r.Pix[r.PixOffset(x, y)]
}

// Opaque just returns true as the alpha channel is fixed.
func (r *RGBImage) Opaque() bool {
	return true
}

// PixOffset returns the index of the Pix array correspinding to the given
// coords
func (r *RGBImage) PixOffset(x, y int) int {
	return (y-r.Rect.Min.Y)*r.Stride + x - r.Rect.Min.X
}

// Set converts the given colour to the RGB space and sets it at the given
// coords
func (r *RGBImage) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(r.Rect)) {
		return
	}
	rr, g, b, _ := c.RGBA()
	r.Pix[r.PixOffset(x, y)] = RGB{
		R: uint8(rr >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
	}
}

// SetRGB directly set an RGB colour to the given coords
func (r *RGBImage) SetRGB(x, y int, rgb RGB) {
	if !(image.Point{x, y}.In(r.Rect)) {
		return
	}
	r.Pix[r.PixOffset(x, y)] = rgb
}

// SubImage retuns the Image viewable through the given bounds
func (r *RGBImage) SubImage(rt image.Rectangle) image.Image {
	rt = rt.Intersect(r.Rect)
	if rt.Empty() {
		return &RGBImage{}
	}
	return &RGBImage{
		Pix:    r.Pix[r.PixOffset(rt.Min.X, rt.Min.Y):],
		Stride: r.Stride,
		Rect:   rt,
	}
}

type rgbImageReader struct {
	*RGBImage
}

func (rg rgbImageReader) ReadColour(x, y int, pixels []byte) {
	rg.SetRGB(x, y, RGB{R: pixels[0], G: pixels[1], B: pixels[2]})
}
