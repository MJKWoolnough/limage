package xcf

import (
	"image"
	"image/color"
)

type rgb struct {
	R, G, B uint8
}

func (rgb rgb) RGBA() (r, g, b, a uint32) {
	r = uint32(rgb.R)
	r |= r << 8
	g = uint32(rgb.G)
	g |= g << 8
	b = uint32(rgb.B)
	b |= b << 8
	return r, g, b, 0xFFFF
}

func rgbColourModel(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	return rgb{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
	}
}

type rgbImage struct {
	Pix    []rgb
	Stride int
	Rect   image.Rectangle
}

func newRGB(r image.Rectangle) *rgbImage {
	w, h := r.Dx(), r.Dy()
	return &rgbImage{
		Pix:    make([]rgb, w*h),
		Stride: w,
		Rect:   r,
	}
}

func (r *rgbImage) At(x, y int) color.Color {
	return r.RGBAt(x, y)
}

func (r *rgbImage) Bounds() image.Rectangle {
	return r.Rect
}

func (r *rgbImage) ColorModel() color.Model {
	return color.ModelFunc(rgbColourModel)
}

func (r *rgbImage) RGBAt(x, y int) rgb {
	if !(image.Point{x, y}.In(r.Rect)) {
		return rgb{}
	}
	return r.Pix[r.PixOffset(x, y)]
}

func (r *rgbImage) Opaque() bool {
	return true
}

func (r *rgbImage) PixOffset(x, y int) int {
	return (y-r.Rect.Min.Y)*r.Stride + x - r.Rect.Min.X
}

func (rg *rgbImage) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(rg.Rect)) {
		return
	}
	r, g, b, _ := c.RGBA()
	rg.Pix[rg.PixOffset(x, y)] = rgb{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
	}
}

func (r *rgbImage) SetRGB(x, y int, rgb rgb) {
	if !(image.Point{x, y}.In(r.Rect)) {
		return
	}
	r.Pix[r.PixOffset(x, y)] = rgb
}

func (rgb *rgbImage) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(rgb.Rect)
	if r.Empty() {
		return &rgbImage{}
	}
	return &rgbImage{
		Pix:    rgb.Pix[rgb.PixOffset(r.Min.X, r.Min.Y):],
		Stride: rgb.Stride,
		Rect:   r,
	}
}

type rgbImageReader struct {
	*rgbImage
}

func (rg rgbImageReader) ReadColour(x, y int, pixels []byte) {
	rg.SetRGB(x, y, rgb{R: pixels[0], G: pixels[1], B: pixels[2]})
}
