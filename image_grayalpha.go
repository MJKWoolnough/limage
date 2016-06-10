package xcf

import (
	"image"
	"image/color"
)

type grayAlpha struct {
	Y, A uint8
}

func (c grayAlpha) RGBA() (r, g, b, a uint32) {
	y := uint32(c.Y)
	y |= y << 8
	a = uint32(c.A)
	a |= a << 8
	return y, y, y, a
}

func grayAlphaColourModel(c color.Color) color.Color {
	_, _, _, a := c.RGBA()
	return grayAlpha{
		Y: color.GrayModel.Convert(c).(color.Gray).Y,
		A: uint8(a >> 8),
	}
}

type grayAlphaImage struct {
	Pix    []grayAlpha
	Stride int
	Rect   image.Rectangle
}

func newGrayAlpha(r image.Rectangle) *grayAlphaImage {
	w, h := r.Dx(), r.Dy()
	return &grayAlphaImage{
		Pix:    make([]grayAlpha, w*h),
		Stride: w,
		Rect:   r,
	}
}

func (g *grayAlphaImage) At(x, y int) color.Color {
	return g.GrayAlphaAt(x, y)
}

func (g *grayAlphaImage) Bounds() image.Rectangle {
	return g.Rect
}

func (g *grayAlphaImage) ColorModel() color.Model {
	return color.ModelFunc(grayAlphaColourModel)
}

func (g *grayAlphaImage) GrayAlphaAt(x, y int) grayAlpha {
	if !(image.Point{x, y}.In(g.Rect)) {
		return grayAlpha{}
	}
	return g.Pix[g.PixOffset(x, y)]
}

func (g *grayAlphaImage) Opaque() bool {
	for _, c := range g.Pix {
		if c.A != 255 {
			return false
		}
	}
	return true
}

func (g *grayAlphaImage) PixOffset(x, y int) int {
	return (y-g.Rect.Min.Y)*g.Stride + x - g.Rect.Min.X
}

func (ga *grayAlphaImage) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(ga.Rect)) {
		return
	}
	ga.Pix[ga.PixOffset(x, y)] = grayAlphaColourModel(c).(grayAlpha)
}

func (g *grayAlphaImage) SetGrayAlpha(x, y int, ga grayAlpha) {
	if !(image.Point{x, y}.In(g.Rect)) {
		return
	}
	g.Pix[g.PixOffset(x, y)] = ga
}

func (g *grayAlphaImage) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(g.Rect)
	if r.Empty() {
		return &grayAlphaImage{}
	}
	return &grayAlphaImage{
		Pix:    g.Pix[g.PixOffset(r.Min.X, r.Min.Y):],
		Stride: g.Stride,
		Rect:   r,
	}
}

type grayAlphaImageReader struct {
	*grayAlphaImage
}

func (ga grayAlphaImageReader) ReadColour(x, y int, cr colourReader) {
	ya := cr.ReadByte()
	a := cr.ReadByte()
	ga.SetGrayAlpha(x, y, grayAlpha{ya, a})
}
