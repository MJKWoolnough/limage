package xcf

import (
	"image"
	"image/color"
)

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

func newGrayAlpha(r image.Rectangle) *grayAlphaImage {
	w, h := r.Dx(), r.Dy()
	return &grayAlphaImage{
		Gray: image.Gray{
			Pix:    make([]uint8, w*h),
			Stride: w,
			Rect:   r,
		},
		Alpha: make([]uint8, w*h),
	}
}

func (g *grayAlphaImage) At(x, y int) color.Color {
	return g.GrayAt(x, y)
}

func (g *grayAlphaImage) ColorModel() color.Model {

}

func (g *grayAlphaImage) GrayAlphaAt(x, y int) grayAlpha {
	p := g.PixOffset(x, y)
	return grayAlpha{
		Y: g.Pix[p],
		A: g.Alpha[p],
	}
}

func (g *grayAlphaImage) Opaque() bool {
	for _, a := range g.Alpha {
		if a != 255 {
			return true
		}
	}
	return false
}

func (g *grayAlphaImage) SetGrayAlpha(x, y int, ga grayAlpha) {
	p := g.PixOffset(x, y)
	g.Pix[p] = ga.Y
	g.Alpha[p] = ga.A
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
