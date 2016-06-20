package xcf

import (
	"image"
	"image/color"
)

type GrayAlpha struct {
	Y, A uint8
}

func (c GrayAlpha) RGBA() (r, g, b, a uint32) {
	y := uint32(c.Y)
	y |= y << 8
	a = uint32(c.A)
	a |= a << 8
	return y, y, y, a
}

func grayAlphaColourModel(c color.Color) color.Color {
	_, _, _, a := c.RGBA()
	return GrayAlpha{
		Y: color.GrayModel.Convert(c).(color.Gray).Y,
		A: uint8(a >> 8),
	}
}

type GrayAlphaImage struct {
	Pix    []GrayAlpha
	Stride int
	Rect   image.Rectangle
}

func newGrayAlpha(r image.Rectangle) *GrayAlphaImage {
	w, h := r.Dx(), r.Dy()
	return &GrayAlphaImage{
		Pix:    make([]GrayAlpha, w*h),
		Stride: w,
		Rect:   r,
	}
}

func (g *GrayAlphaImage) At(x, y int) color.Color {
	return g.GrayAlphaAt(x, y)
}

func (g *GrayAlphaImage) Bounds() image.Rectangle {
	return g.Rect
}

func (g *GrayAlphaImage) ColorModel() color.Model {
	return color.ModelFunc(grayAlphaColourModel)
}

func (g *GrayAlphaImage) GrayAlphaAt(x, y int) GrayAlpha {
	if !(image.Point{x, y}.In(g.Rect)) {
		return GrayAlpha{}
	}
	return g.Pix[g.PixOffset(x, y)]
}

func (g *GrayAlphaImage) Opaque() bool {
	for _, c := range g.Pix {
		if c.A != 255 {
			return false
		}
	}
	return true
}

func (g *GrayAlphaImage) PixOffset(x, y int) int {
	return (y-g.Rect.Min.Y)*g.Stride + x - g.Rect.Min.X
}

func (ga *GrayAlphaImage) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(ga.Rect)) {
		return
	}
	ga.Pix[ga.PixOffset(x, y)] = grayAlphaColourModel(c).(GrayAlpha)
}

func (g *GrayAlphaImage) SetGrayAlpha(x, y int, ga GrayAlpha) {
	if !(image.Point{x, y}.In(g.Rect)) {
		return
	}
	g.Pix[g.PixOffset(x, y)] = ga
}

func (g *GrayAlphaImage) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(g.Rect)
	if r.Empty() {
		return &GrayAlphaImage{}
	}
	return &GrayAlphaImage{
		Pix:    g.Pix[g.PixOffset(r.Min.X, r.Min.Y):],
		Stride: g.Stride,
		Rect:   r,
	}
}

type grayAlphaImageReader struct {
	*GrayAlphaImage
}

func (ga grayAlphaImageReader) ReadColour(x, y int, pixels []byte) {
	ga.SetGrayAlpha(x, y, GrayAlpha{pixels[0], pixels[1]})
}
