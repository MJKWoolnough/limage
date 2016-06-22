package xcf

import (
	"image"
	"image/color"
)

type PalettedAlpha struct {
	Pix     []IndexedAlpha
	Stride  int
	Rect    image.Rectangle
	Palette color.Palette
}

type alphaPalette struct {
	color.Palette
}

func (ap alphaPalette) Convert(c color.Color) color.Color {
	r, g, b, _ := ap.Palette.Convert(c).RGBA()
	_, _, _, a := c.RGBA()
	return image.RGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}

type IndexedAlpha struct {
	I, A uint8
}

func newPalettedAlpha(r image.Rectangle, p color.Palette) *PalettedAlpha {
	w, h := r.Dx(), r.Dy()
	return &PalettedAlpha{
		Pix:     make([]IndexedAlpha, w*h),
		Stride:  w,
		Rect:    r,
		Palette: p,
	}
}

func (p *PalettedAlpha) At(x, y int) color.Color {
	if p.Palette == nil {
		return nil
	}
	ia := p.IndexAlphaAt(x, y)
	r, g, b, _ := p.Palette[ia.I].RGBA()
	return color.NRGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: ia.A,
	}
}

func (p *PalettedAlpha) Bounds() image.Rectangle {
	return p.Rect
}

func (p *PalettedAlpha) ColorModel() color.Model {
	return alphaPalette{p.Palette}
}

func (p *PalettedAlpha) IndexAlphaAt(x, y int) IndexedAlpha {
	if !(image.Point{x, y}.In(p.Rect)) {
		return IndexedAlpha{}
	}
	return p.Pix[p.PixOffset(x, y)]
}

func (p *PalettedAlpha) Opaque() bool {
	for _, c := range p.Pix {
		if c.A != 255 {
			return true
		}
	}
	return false
}
func (p *PalettedAlpha) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*1
}

func (p *PalettedAlpha) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	_, _, _, a := c.RGBA()
	p.Pix[p.PixOffset(x, y)] = IndexedAlpha{
		I: uint8(p.Palette.Index(c)),
		A: uint8(a >> 8),
	}
}

func (p *PalettedAlpha) SetIndexAlpha(x, y int, ia IndexedAlpha) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	p.Pix[p.PixOffset(x, y)] = ia
}

func (p *PalettedAlpha) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	if r.Empty() {
		return &PalettedAlpha{}
	}
	return &PalettedAlpha{
		Pix:     p.Pix[p.PixOffset(r.Min.X, r.Min.Y):],
		Stride:  p.Stride,
		Rect:    r,
		Palette: p.Palette,
	}
}

type palettedAlphaReader struct {
	*PalettedAlpha
}

func (p palettedAlphaReader) ReadColour(x, y int, pixels []byte) {
	p.SetIndexAlpha(x, y, IndexedAlpha{
		I: pixels[0],
		A: pixels[1],
	})
}
