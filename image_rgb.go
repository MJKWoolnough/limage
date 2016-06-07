package xcf

import (
	"image"
	"image/color"
)

type rgb struct {
	R, G, B uint8
}

func (rgb rgb) RGBA() (r, g, b, a uint32) {
	r = rgb.R
	r |= r << 8
	g = rgb.G
	g |= g << 8
	b = rgb.B
	b |= b << 8
	return r, g, b, 0xFFFF
}

func rgbColourModel(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	return rgb{
		R: r >> 8,
		G: g >> 8,
		B: b >> 8,
	}
}

type rgbImage struct {
	Pix    []rgb
	Stride int
	Rect   image.Rectangle
}

func newRGBImage(r image.Rectangle) *rgbImage {
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

func (r *rgbImage) ColorModel() color.Model {
	return rgbColourModel
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

func (rgb *rgbImage) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(rgb.Rect)) {
		return
	}
	r, g, b, _ := c.RGBA()
	rgb.Pix[r.PixOffset(x, y)] = rgb{
		R: r >> 8,
		G: g >> 8,
		B: b >> 8,
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
		Pix:    rgb.Pix[rgb.PixOffset(r.Min.X, r.Min.Y)],
		Stride: rgb.Stride,
		Rect:   r,
	}
}

type rgbImageReader struct {
	*rgb
}

func (rgb rgbImageReader) ReadColour(x, y int, cr colourReader) {
	r := cr.ReadByte()
	g := cr.ReadByte()
	b := cr.ReadByte()
	rgb.SetRGB(x, y, rgb{r, g, b})
}
