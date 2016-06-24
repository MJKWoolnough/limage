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
	c.A = uint16((uint32(i.Opacity) * uint32(c.A)) / 0xff)
	return c
}

type Layer struct {
	Name             string
	OffsetX, OffsetY int
	Mode             uint32
	Visible          bool
	Opacity          uint8
	image.Image
}

func (l *Layer) Bounds() image.Rectangle {
	max := l.Image.Bounds().Max
	return image.Rect(l.OffsetX, l.OffsetY, max.X+l.OffsetX, max.Y+l.OffsetY)
}

func (l *Layer) At(x, y int) color.Color {
	if l.Opacity == 255 {
		return l.Image.At(x-l.OffsetX, y-l.OffsetY)
	}
	c := colourToNRGBA(l.Image.At(x-l.OffsetX, y-l.OffsetY))
	c.A = uint16((uint32(l.Opacity) * uint32(c.A)) / 0xff)
	return c
}

type Group struct {
	image.Config
	Layers []Layer
}

func (g *Group) ColorModel() color.Model {
	return g.Config.ColorModel
}

func (g *Group) Bounds() image.Rectangle {
	return image.Rect(0, 0, g.Width, g.Height)
}

func (g *Group) At(x, y int) color.Color {
	var c color.Color = color.Alpha{}
	point := image.Point{x, y}
	for i := len(g.Layers) - 1; i >= 0; i-- {
		if !g.Layers[i].Visible {
			continue
		}
		if !point.In(g.Layers[i].Bounds()) {
			continue
		}
		ca := g.Layers[i].At(x, y)
		ar, ag, ab, aa := ca.RGBA()
		if aa == 0xffff {
			c = ca
			continue
		} else if aa == 0 {
			continue
		}

		br, bg, bb, ba := c.RGBA()

		ma := 0xffff - aa
		c = color.RGBA64{
			R: uint16((br*ma + ar*0xffff) / 0xffff),
			G: uint16((bg*ma + ag*0xffff) / 0xffff),
			B: uint16((bb*ma + ab*0xffff) / 0xffff),
			A: uint16((ba*ma + aa*0xffff) / 0xffff),
		}

	}
	return c
}

type MaskedImage struct {
	image.Image
	Mask *image.Gray
}

func (m *MaskedImage) At(x, y int) color.Color {
	mask := m.Mask.GrayAt(x, y)
	if mask.Y == 0 {
		return color.Alpha{}
	} else if mask.Y == 0xff {
		return m.Image.At(x, y)
	}
	switch i := m.Image.(type) {
	case *RGBImage:
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
		return GrayAlpha{
			Y: c.Y,
			A: mask.Y,
		}
	case *GrayAlphaImage:
		c := i.GrayAlphaAt(x, y)
		c.A = uint8((uint32(mask.Y) * uint32(c.A)) >> 8)
		return c
	case *image.Paletted:
		c := i.Palette[i.ColorIndexAt(x, y)].(RGB)
		return color.NRGBA{
			R: c.R,
			G: c.G,
			B: c.B,
			A: mask.Y,
		}
	case *PalettedAlpha:
		ca := i.IndexAlphaAt(x, y)
		c := i.Palette[ca.I].(RGB)
		return color.NRGBA{
			R: c.R,
			G: c.G,
			B: c.B,
			A: uint8((uint32(mask.Y) * uint32(ca.A)) >> 8),
		}
	default: // shouldn't happen (I think)
		c := colourToNRGBA(i.At(x, y))
		c.A = uint16((uint32(mask.Y) * uint32(c.A)))
		return c
	}
}

func colourToNRGBA(c color.Color) color.NRGBA64 {
	r, g, b, a := c.RGBA()
	if a == 0 {
		return color.NRGBA64{}
	}
	return color.NRGBA64{
		R: uint16(((r * 0xffff) / a)),
		G: uint16(((g * 0xffff) / a)),
		B: uint16(((b * 0xffff) / a)),
		A: uint16(a),
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
