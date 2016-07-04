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

func (hsv HSVA) ToNRGBA() color.NRGBA64 {
	if hsv.S == 0 {
		return color.NRGBA64{
			R: hsv.V,
			G: hsv.V,
			B: hsv.V,
			A: hsv.A,
		}
	}
	cl := color.NRGBA64{
		A: hsv.A,
	}
	c := uint16(uint32(hsv.V) * uint32(hsv.S) / 0xffff)
	var h uint32
	if hsv.H != 0xffff {
		h = uint32(hsv.H) * 6
	}
	x := h % 0x1fffe
	if x >= 0xffff {
		x = 0x1fffe - ha
	}
	x = x * c / 0xffff
	switch h / 0xffff {
	case 0:
		cl.R = c
		cl.G = x
	case 1:
		cl.R = x
		cl.G = c
	case 2:
		cl.G = c
		cl.B = x
	case 3:
		cl.G = x
		cl.B = c
	case 4:
		cl.R = x
		cl.B = c
	case 5:
		cl.R = c
		cl.B = x
	}
	m := hsv.V - uint16(c)
	cl.R += m
	cl.G += m
	cl.B += m
	return cl
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
