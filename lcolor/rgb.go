package lcolor // import "vimagination.zapto.org/limage/lcolor"

import "image/color"

// RGB is a standard colour type whose Alpha channel is always full
type RGB struct {
	R, G, B uint8
}

// RGBA implements the color.Color interface
func (rgb RGB) RGBA() (r, g, b, a uint32) {
	r = uint32(rgb.R)
	r |= r << 8
	g = uint32(rgb.G)
	g |= g << 8
	b = uint32(rgb.B)
	b |= b << 8
	return r, g, b, 0xFFFF
}

// ToNRGBA returns itself as a non-alpha-premultiplied value
// As the alpha is always full, this only returns the normal values
func (rgb RGB) ToNRGBA() color.NRGBA64 {
	r := uint16(rgb.R)
	r |= r << 8
	g := uint16(rgb.G)
	g |= g << 8
	b := uint16(rgb.B)
	b |= b << 8
	return color.NRGBA64{r, g, b, 0xffff}
}

func rgbColourModel(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	return RGB{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
	}
}
