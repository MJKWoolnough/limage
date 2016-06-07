package xcf

import (
	"image"
	"image/color"
)

type colourReader interface {
	ReadByte() byte
}

type rgbaImageReader struct {
	*image.NRGBA
}

func (rgba rgbaImageReader) ReadColour(x, y int, cr colourReader) {
	r := cr.ReadByte()
	g := cr.ReadByte()
	b := cr.ReadByte()
	a := cr.ReadByte()
	rgba.SetNRGBA(x, y, color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	})
}

type grayImageReader struct {
	*image.Gray
}

func (g grayImageReader) ReadColour(x, y int, cr colourReader) {
	yc := cr.ReadByte()
	g.SetGray(x, y, color.Gray{yc})
}

type indexedImageReader struct {
	*image.Paletted
}

func (p indexedImageReader) ReadColour(x, y, int, cr colourReader) {
	i := cr.ReadByte()
	p.SetColorIndex(x, y, i)
}

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

type grayAlpha struct {
	Y, A uint8
}

func (g grayAlpha) RGBA() (r, g, b, a uint32) {
	y := uint32(c.Y)
	y |= y << 8
	a := uint32(c.A)
	a |= a << 8
	return y, y, y, a
}

type grayAlphaImage struct {
	image.Gray
	Alpha []uint8
}

func newGrayAlpha(r image.Rect) *grayAlphaImage {

}

func (g *grayAlphaImage) At(x, y int) color.Color {

}

func (g *grayAlphaImage) ColorModel() color.Model {

}

func (g *grayAlphaImage) GrayAlphaAt(x, y int) grayAlpha {

}

func (g *grayAlphaImage) Opaque() bool {

}

func (g *grayAlphaImage) SetGrayAlpha(x, y int, ga grayAlpha) {

}

func (g *grayAlphaImage) SubImage(r image.Rectangle) image.Image {

}

type grayAlphaImageReader struct {
	*grayAlphaImage
}

func (g greyAlphaImageReader) ReadColour(x, y int, cr colourReader) {
	y := cr.ReadByte()
	a := cr.ReadByte()
	g.SetGray(x, y, grayAlpha{y, a})
}

type palettedAlpha struct {
	image.Paletted
	Alpha []uint8
}

type indexedAlpha struct {
	I, A uint8
}

func newPalettedAlpha(r image.Rect, p color.Palette) *palettedAlpha {

}

func (p *palettedAlpha) At(x, y int) color.Color {

}

func (g *palettedAlpha) ColorModel() {

}

func (g *palettedAlpha) IndexAlphaAt(x, y int) indexedAlpha {

}

func (g *palettedAlpha) Opaque() {
}

func (g *palettedAlpha) SetIndexAlpha(x, y int, ia indexedAlpha) {
}

func (g *palettedAlpha) SubImage(r image.Rectangle) image.Image {
}
