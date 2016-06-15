package xcf

import (
	"image"
	"image/color"
)

type Image struct {
	Group
	Comment string
	Opacity uint8
}

func (i *Image) At(x, y int) color.Color {
	if i.Opacity == 255 {
		return i.Group.At(x, y)
	}
	c := colourToNRGBA(i.Group.At(x, y))
	c.A = uint8((uint32(i.Opacity) * uint32(c.A)) >> 8)
	return c
}

type Layer struct {
	OffsetX, OffsetY int
	Mode             uint32
	image.Image
}

func (l *Layer) Bounds() image.Rectangle {
	max := l.Image.Bounds().Max
	return image.Rect(l.OffsetX, l.OffsetY, max.X+l.OffsetX, max.Y+l.OffsetY)
}

func (l *Layer) At(x, y int) color.Color {
	return l.Image.At(x-l.OffsetX, y-l.OffsetY)
}

type Group struct {
	Name          string
	Width, Height int
	colorModel    color.Model
	Layers        []Layer
}

func (g *Group) ColorModel() color.Model {
	return g.colorModel
}

func (g *Group) Bounds() image.Rectangle {
	return image.Rect(0, 0, g.Width, g.Height)
}

func (g *Group) At(x, y int) color.Color {
	return nil
}

type MaskedImage struct {
	image.Image
	Mask image.Image
}

func (m *MaskedImage) At(x, y int) color.Color {
	mask := m.Image.At(x, y).(color.Gray)
	switch i := m.Image.(type) {
	case *rgbImage:
		c := i.RGBAt(x, y)
		return color.NRGBA{
			R: c.R,
			G: c.G,
			B: c.B,
			A: mask.Y,
		}
	case *image.NRGBA:
		c := i.NRGBAAt(x, y)
		c.A = uint8((uint32(mask.Y) * uint32(c.A)) >> 8)
		return c
	case *image.Gray:
		c := i.GrayAt(x, y)
		return grayAlpha{
			Y: c.Y,
			A: mask.Y,
		}
	case *grayAlphaImage:
		c := i.GrayAlphaAt(x, y)
		c.A = uint8((uint32(mask.Y) * uint32(c.A)) >> 8)
		return c
	case *image.Paletted:
		c := i.Palette[i.ColorIndexAt(x, y)].(rgb)
		return color.NRGBA{
			R: c.R,
			G: c.G,
			B: c.B,
			A: mask.Y,
		}
	case *palettedAlpha:
		ca := i.IndexAlphaAt(x, y)
		c := i.Palette[ca.I].(rgb)
		return color.NRGBA{
			R: c.R,
			G: c.G,
			B: c.B,
			A: uint8((uint32(mask.Y) * uint32(ca.A)) >> 8),
		}
	default: // shouldn't happen (I think)
		c := colourToNRGBA(i.At(x, y))
		c.A = uint8((uint32(mask.Y) * uint32(c.A)) >> 8)
		return c
	}
}

func colourToNRGBA(c color.Color) color.NRGBA {
	r, g, b, a := c.RGBA()
	return color.NRGBA{
		R: uint8(((r * 0xffff) / a) >> 8),
		G: uint8(((g * 0xffff) / a) >> 8),
		B: uint8(((b * 0xffff) / a) >> 8),
		A: uint8(a >> 8),
	}
}

type Text struct {
	image.Image
	TextData
}

type TextData []TextDatum

func (t TextData) String() string {
	toRet := ""
	for _, d := range t {
		toRet += d.Data
	}
	return toRet
}

type TextDatum struct {
	ForeColor, BackColor                   color.Color
	Size, LetterSpacing, Rise              float64
	Bold, Italic, Underline, Strikethrough bool
	Font, Data                             string
	FontUnit                               uint8
}
