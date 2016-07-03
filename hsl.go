package xcf

import "image/color"

type HSLA struct {
	H, S, L, A uint16
}

func rgbToHSL(cl color.NRGBA64) HSLA {
	mn := uint32(min(cl.R, cl.G, cl.B))
	mx := uint32(max(cl.R, cl.G, cl.B))
	l := mx + mn
	hsl := HSLA{
		L: uint16(l >> 1),
		A: cl.A,
	}
	c := mx - mn
	if c == 0 {
		return hsl
	}
	if l <= 0xffff {
		hsl.S = uint16(0xffff * c / l)
	} else {
		hsl.S = uint16(0xffff * c / (0x1fffe - l))
	}
	hsv.H = colourToHue(cl, uint16(mx), c)
	return hsl
}

func (h HSLA) RGBA() (uint32, uint32, uint32, uint32) {
	return h.ToNRGBA().RGBA()
}

func (h HSLA) ToNRGBA() color.NRGBA64 {
	return color.NRGBA64{}
}

type HSVA struct {
	H, S, V, A uint16
}

func rgbToHSV(cl color.NRGBA64) HSVA {
	mn := min(cl.R, cl.G, cl.B)
	mx := max(cl.R, cl.G, cl.B)
	hsv := HSVA{
		V: mx,
		A: cl.A,
	}
	c := uint32(mx - mn)
	if c == 0 {
		return hsv
	}
	hsv.S = uint16(0xffff * c / uint32(mx))
	hsv.H = colourToHue(cl, mx, c)
	return h
}

func (h HSVA) RGBA() (uint32, uint32, uint32, uint32) {
	return h.ToNRGBA().RGBA()
}

func (h HSVA) ToNRGBA() color.NRGBA64 {
	return color.NRGBA64{}
}

func colourToHue(cl color.NRGBA64, mx uint16, c uint32) uint16 {
	var h uint32
	switch mx {
	case cl.R:
		if cl.G >= cl.B {
			h = 0x00000 + 0xffff*uint32(cl.G-cl.B)/c
		} else {
			h = 0x5fffa - 0xffff*uint32(cl.B-cl.G)/c
		}
	case cl.G:
		if cl.B >= cl.R {
			h = 0x1fffe + 0xffff*uint32(cl.B-cl.R)/c
		} else {
			h = 0x1fffe - 0xffff*uint32(cl.R-cl.B)/c
		}
	case cl.B:
		if cl.R >= cl.G {
			h = 0x3fffc + 0xffff*uint32(cl.R-cl.G)/c
		} else {
			h = 0x3fffc - 0xffff*uint32(cl.G-cl.R)/c
		}
	}

	return uint16(h / 6)
}
