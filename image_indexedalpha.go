package xcf

import (
	"image"
	"image/color"
)

type palettedAlpha struct {
	Pix     []indexedAlpha
	Stride  int
	Rect    image.Rectangle
	Palette color.Palette
}

type indexedAlpha struct {
	I, A uint8
}

func newPalettedAlpha(r image.Rectangle, p color.Palette) *palettedAlpha {
	w, h := r.Dx(), r.Dy()
	return &palettedAlpha{
		Pix:     make([]indexedAlpha, w*h),
		Stride:  w,
		Rect:    r,
		Palette: p,
	}
}

func (p *palettedAlpha) At(x, y int) color.Color {
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

func (p *palettedAlpha) Bounds() image.Rectangle {
	return p.Rect
}

func (p *palettedAlpha) ColorModel() color.Model {
	return p.Palette
}

func (p *palettedAlpha) IndexAlphaAt(x, y int) indexedAlpha {
	if !(image.Point{x, y}.In(p.Rect)) {
		return indexedAlpha{}
	}
	return p.Pix[p.PixOffset(x, y)]
}

func (p *palettedAlpha) Opaque() bool {
	for _, c := range p.Pix {
		if c.A != 255 {
			return true
		}
	}
	return false
}
func (p *palettedAlpha) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*1
}

func (p *palettedAlpha) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	_, _, _, a := c.RGBA()
	p.Pix[p.PixOffset(x, y)] = indexedAlpha{
		I: uint8(p.Palette.Index(c)),
		A: uint8(a >> 8),
	}
}

func (p *palettedAlpha) SetIndexAlpha(x, y int, ia indexedAlpha) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	p.Pix[p.PixOffset(x, y)] = ia
}

func (p *palettedAlpha) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	if r.Empty() {
		return &palettedAlpha{}
	}
	return &palettedAlpha{
		Pix:     p.Pix[p.PixOffset(r.Min.X, r.Min.Y):],
		Stride:  p.Stride,
		Rect:    r,
		Palette: p.Palette,
	}
}

type palettedAlphaReader struct {
	*palettedAlpha
}

func (p palettedAlphaReader) ReadColour(x, y int, pixels []byte) {
	p.SetIndexAlpha(x, y, indexedAlpha{
		I: pixels[0],
		A: pixels[1],
	})
}
