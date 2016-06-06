package xcf

import (
	"image"
	"image/color"
	"io"
)

type colourReader interface {
	ReadByte() byte
	io.Seeker
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

func newGrayAlpha(r image.Rect) grayAlphaImage {

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
