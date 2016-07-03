package xcf

import "image/color"

type HSLA struct {
	H, S, L, A uint16
}

func rgbToHSL(cl color.NRGBA64) HSLA {
	mn := uint32(min(cl.R, cl.G, cl.B))
	mx := uint32(max(cl.R, c.G, cl.B))
	l := mx + mn
	hsl := HSLA{
		L: uint16(l >> 1),
		A: cl.A,
	}
	c := max - min
	if c == 0 {
		return hsl
	}
	if l <= 0xffff {
		hsl.S = uint16(0xffff * c / l)
	} else {
		hsl.S = uint16(0xffff * c / (0x1fffe - l))
	}
	var h uint32
	switch uint16(mx) {
	case cl.R:
		if cl.G > cl.B {
			h = 0xffff * (cl.G - cl.B) / c
		} else {
			h = 0x5fffa - 0xffff*(cl.B-cl.G)/c
		}
	case cl.G:
		if cl.B > cl.R {
			h = 0x1fffe + 0xffff*(cl.B-cl.R)/c
		} else {
			h = 0x1fffe - 0xffff*(cl.R-cl.B)/c
		}
	case cl.B:
		if cl.R > cl.G {
			h = 0x3fffc + 0xffff*(cl.R-cl.G)/c
		} else {
			h = 0x3fffc - 0xffff*(cl.G-cl.R)/c
		}
	}

	hsl.H = h / 6

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
	return HSVA{}
}

func (h HSVA) RGBA() (uint32, uint32, uint32, uint32) {
	return h.ToNRGBA().RGBA()
}

func (h HSVA) ToNRGBA() color.NRGBA64 {
	return color.NRGBA64{}
}
