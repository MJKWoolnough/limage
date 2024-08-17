package lcolor

import "image/color"

// AlphaPalette is an extension to the normal color.Palette type.
type AlphaPalette color.Palette

// Convert converts the given color to the nearest color in the palette, but
// preserves the alpha channel.
func (ap AlphaPalette) Convert(c color.Color) color.Color {
	r, g, b, _ := color.Palette(ap).Convert(c).RGBA()
	_, _, _, a := c.RGBA()

	return color.NRGBA64{
		R: uint16(r),
		G: uint16(g),
		B: uint16(b),
		A: uint16(a),
	}
}

// Index returns the palette index of the nearest color.
func (ap AlphaPalette) Index(c color.Color) int {
	return color.Palette(ap).Index(c)
}

// IndexedAlpha is the combination of a palette index and an Alpha channel.
type IndexedAlpha struct {
	I, A uint8
}
