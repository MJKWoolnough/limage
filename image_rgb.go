package xcf

import (
	"image"
	"image/color"
)

type RGB struct {
	R, G, B uint8
}

func (rgb RGB) RGBA() (r, g, b, a uint32) {
	r = uint32(rgb.R)
	r |= r << 8
	g = uint32(rgb.G)
	g |= g << 8
	b = uint32(rgb.B)
	b |= b << 8
	return r, g, b, 0xFFFF
}

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

func (r *RGBImage) At(x, y int) color.Color {
	return r.RGBAt(x, y)
}

func (r *RGBImage) Bounds() image.Rectangle {
	return r.Rect
}

func (r *RGBImage) ColorModel() color.Model {
	return color.ModelFunc(rgbColourModel)
}

func (r *RGBImage) RGBAt(x, y int) RGB {
	if !(image.Point{x, y}.In(r.Rect)) {
		return RGB{}
	}
	return r.Pix[r.PixOffset(x, y)]
}

func (r *RGBImage) Opaque() bool {
	return true
}

func (r *RGBImage) PixOffset(x, y int) int {
	return (y-r.Rect.Min.Y)*r.Stride + x - r.Rect.Min.X
}

func (rg *RGBImage) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(rg.Rect)) {
		return
	}
	r, g, b, _ := c.RGBA()
	rg.Pix[rg.PixOffset(x, y)] = RGB{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
	}
}

func (r *RGBImage) SetRGB(x, y int, rgb RGB) {
	if !(image.Point{x, y}.In(r.Rect)) {
		return
	}
	r.Pix[r.PixOffset(x, y)] = rgb
}

func (rgb *RGBImage) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(rgb.Rect)
	if r.Empty() {
		return &RGBImage{}
	}
	return &RGBImage{
		Pix:    rgb.Pix[rgb.PixOffset(r.Min.X, r.Min.Y):],
		Stride: rgb.Stride,
		Rect:   r,
	}
}

type rgbImageReader struct {
	*RGBImage
}

func (rg rgbImageReader) ReadColour(x, y int, pixels []byte) {
	rg.SetRGB(x, y, RGB{R: pixels[0], G: pixels[1], B: pixels[2]})
}
