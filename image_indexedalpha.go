package xcf

import (
	"image"
	"image/color"
)

type palettedAlpha struct {
	image.Paletted
	Alpha []uint8
}

type indexedAlpha struct {
	I, A uint8
}

func newPalettedAlpha(r image.Rect, p color.Palette) *palettedAlpha {

}

func (p *palettedAlpha) At(x, y int) color.Color {

}

func (g *palettedAlpha) ColorModel() {

}

func (g *palettedAlpha) IndexAlphaAt(x, y int) indexedAlpha {

}

func (g *palettedAlpha) Opaque() {
}

func (g *palettedAlpha) SetIndexAlpha(x, y int, ia indexedAlpha) {
}

func (g *palettedAlpha) SubImage(r image.Rectangle) image.Image {
}
