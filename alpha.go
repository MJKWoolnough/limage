package xcf

import (
	"image"
	"image/color"
)

type rgb struct {
	R, G, B uint8
}

func (c rgb) RGBA() (uint32, uint32, uint32, uint32) {
	r := uint32(c.R)
	r |= c.R
	g := uint32(c.G)
	g |= c.G
	b := uint32(c.B)
	b |= c.B
	return r, g, b, 65535
}

type PalettedAlpha struct {
	image.Paletted
	alpha []color.Alpha
}

func NewPalettedAlpha(r image.Rectangle, p color.Palette) *PalettedAlpha {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 1*w*h)
	alpha := make([]color.Alpha, 1*w*h)
	return &PalettedAlpha{
		Paletted: image.Paletted{pix, 1 * w, r, p},
		alpha:    alpha,
	}
}

func (p *PalettedAlpha) At(x, y int) color.Color {
	if len(p.Palette) == 0 {
		return nil
	}
	if !(Point{x, y}.In(p.Rect)) {
		return p.Palette[0]
	}
	i := p.PixOffset(x, y)
	c, _ := p.Palette[p.Pix[i]].(rgb)
	a := p.alpha[i].A
	return color.NRGBA{
		R: c.R,
		G: c.G,
		B: c.B,
		A: a,
	}
}

func (p *PalettedAlpha) Opaque() bool {
	for _, a := range p.alpha {
		if a != 255 {
			return false
		}
	}
	return true
}

func (p *PalettedAlpha) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	r, g, b, a := c.RGBA()
	p.Pix[i] = uint8(p.Palette.Index(rgb{uint8(r), uint8(g), uint8(b)}))
	p.alpha[i].A = uint8(a)
}

func (p *PalettedAlpha) SubImage(r Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	if r.Empty() {
		return &PalettedAlpha{
			Paletted: image.Paletted{
				Palette: p.Palette,
			},
		}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &PalettedAlpha{
		Palette: image.Palette{
			Pix:     p.Pix[i:],
			Stride:  p.Stride,
			Rect:    p.Rect.Intersect(r),
			Palette: p.Palette,
		},
		alpha: p.alpha[i:],
	}
}
