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
	Mode             Composite
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
		if _, ok := g.Config.ColorModel.(color.Palette); g.Layers[i].Mode != CompositeDissolve && ok {
			d := g.Layers[i].At(x, y)
			if ar, ag, ab, aa := d.RGBA(); aa > 0x7fff {
				c = color.RGBA64{
					R: uint16(ar),
					G: uint16(ag),
					B: uint16(ab),
					A: 0xffff,
				}
			}
		} else {
			c = g.Layers[i].Mode.Composite(colourToNRGBA(c), colourToNRGBA(g.Layers[i].At(x, y)))
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
	switch c := c.(type) {
	case color.NRGBA:
		var d color.NRGBA64
		d.R = uint16(c.R)
		d.R |= d.R << 8
		d.G = uint16(c.G)
		d.G |= d.G << 8
		d.B = uint16(c.B)
		d.B |= d.B << 8
		d.A = uint16(c.A)
		d.A |= d.A << 8
		return d
	case color.NRGBA64:
		return c
	}
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
