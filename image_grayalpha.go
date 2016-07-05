package xcf

import (
	"image"
	"image/color"
)

// GrayAlpha represents a Gray color with an Alpha channel
type GrayAlpha struct {
	Y, A uint8
}

// RGBA implements the color.Color interface
func (c GrayAlpha) RGBA() (r, g, b, a uint32) {
	y := uint32(c.Y)
	y |= y << 8
	a = uint32(c.A)
	y *= a
	y /= 0xff
	a |= a << 8
	return y, y, y, a
}

// ToNRGBA converts the HSV color into the RGB colorspace
func (c GrayAlpha) ToNRGBA() color.NRGBA64 {
	y := uint16(c.Y)
	y |= y << 8
	a := uint16(c.A)
	a |= a << 8
	return color.NRGBA64{y, y, y, a}
}

func grayAlphaColourModel(c color.Color) color.Color {
	_, _, _, a := c.RGBA()
	return GrayAlpha{
		Y: color.GrayModel.Convert(c).(color.Gray).Y,
		A: uint8(a >> 8),
	}
}

// GrayAlphaImage is an image of GrayAlpha pixels
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

// At returns the color for the pixel at the specified coords
func (g *GrayAlphaImage) At(x, y int) color.Color {
	return g.GrayAlphaAt(x, y)
}

// Bounds returns the limits of the image
func (g *GrayAlphaImage) Bounds() image.Rectangle {
	return g.Rect
}

// ColorModel returns a color model to transform arbitrary colours into a
// GrayAlpha color
func (g *GrayAlphaImage) ColorModel() color.Model {
	return color.ModelFunc(grayAlphaColourModel)
}

// GrayAlphaAt returns a GrayAlpha colr for the specified coords
func (g *GrayAlphaImage) GrayAlphaAt(x, y int) GrayAlpha {
	if !(image.Point{x, y}.In(g.Rect)) {
		return GrayAlpha{}
	}
	return g.Pix[g.PixOffset(x, y)]
}

// Opaque returns true if all pixels have full alpha
func (g *GrayAlphaImage) Opaque() bool {
	for _, c := range g.Pix {
		if c.A != 255 {
			return false
		}
	}
	return true
}

// PixOffset returns the index of the element of Pix corresponding to the given
// coords
func (g *GrayAlphaImage) PixOffset(x, y int) int {
	return (y-g.Rect.Min.Y)*g.Stride + x - g.Rect.Min.X
}

// Set converts the given colour to a GrayAlpha colour and sets it at the given
// coords
func (g *GrayAlphaImage) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(g.Rect)) {
		return
	}
	g.Pix[g.PixOffset(x, y)] = grayAlphaColourModel(c).(GrayAlpha)
}

// SetGrayAlpha sets the colour at the given coords
func (g *GrayAlphaImage) SetGrayAlpha(x, y int, ga GrayAlpha) {
	if !(image.Point{x, y}.In(g.Rect)) {
		return
	}
	g.Pix[g.PixOffset(x, y)] = ga
}

// SubImage retuns the Image viewable through the given bounds
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
