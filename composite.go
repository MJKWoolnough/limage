package xcf

import (
	"image/color"
	"math/rand"
)

type Composite uint32

const (
	CompositeNormal       Composite = 0
	CompositeDissolve     Composite = 1
	CompositeBehind       Composite = 2
	CompositeMultiply     Composite = 3
	CompositeScreen       Composite = 4
	CompositeOverlay      Composite = 5
	CompositeDifference   Composite = 6
	CompositeAddition     Composite = 7
	CompositeSubtract     Composite = 8
	CompositeDarkenOnly   Composite = 9
	CompositeLightenOnly  Composite = 10
	CompositeHue          Composite = 11
	CompositeSaturation   Composite = 12
	CompositeColor        Composite = 13
	CompositeValue        Composite = 14
	CompositeDivide       Composite = 15
	CompositeDodge        Composite = 16
	CompositeBurn         Composite = 17
	CompositeHardLight    Composite = 18
	CompositeSoftLight    Composite = 19
	CompositeGrainExtract Composite = 20
	CompositeGrainMerge   Composite = 21
)

func (c Composite) Composite(bottom, top color.Color) color.Color {
	var f func(uint32, uint32) uint32
	switch c {
	case CompositeDissolve:
		return compositeDissolve(bottom, top)
	case CompositeBehind:
		return bottom
	case CompositeMultiply:
		f = compositeMultiply
	case CompositeScreen:
		f = compositeScreen
	case CompositeOverlay:
		f = compositeOverlay
	case CompositeDifference:
		f = compositeDifference
	case CompositeAddition:
		f = compositeAddition
	case CompositeSubtract:
		f = compositeSubtract
	case CompositeDarkenOnly:
		f = compositeDarkenOnly
	case CompositeLightenOnly:
		f = compositeLightenOnly
	case CompositeHue:
		return compositeHue(bottom, top)
	case CompositeSaturation:
		return compositeSaturation(bottom, top)
	case CompositeColor:
		return compositeColor(bottom, top)
	case CompositeValue:
		return compositeValue(bottom, top)
	case CompositeDivide:
		f = compositeDivide
	case CompositeDodge:
		f = compositeDodge
	case CompositeBurn:
		f = compositeBurn
	case CompositeHardLight:
		f = compositeHardLight
	case CompositeSoftLight:
		f = compositeSoftLight
	case CompositeGrainExtract:
		f = compositeGrainExtract
	case CompositeGrainMerge:
		f = compositeGrainMerge
	default: //Normal
		return compositeNormal(bottom, top)
	}
	ar, ag, ab, aa := bottom.RGBA()
	br, bg, bb, ba := top.RGBA()
	return color.RGBA64{
		R: uint16(blend(aa, ar, min(aa, ba), f(ar, br))),
		G: uint16(blend(aa, ag, min(aa, ba), f(ag, bg))),
		B: uint16(blend(aa, ab, min(aa, ba), f(ab, bb))),
		A: uint16(aa),
	}
}

func compositeNormal(bottom, top color.Color) color.Color {
	ar, ag, ab, aa := bottom.RGBA()
	br, bg, bb, ba := top.RGBA()
	return color.RGBA64{
		R: uint16(blend(aa, ar, ba, br)),
		G: uint16(blend(aa, ag, ba, bg)),
		B: uint16(blend(aa, ab, ba, bb)),
		A: uint16(0xffff - (0xffff-aa)*(0xffff-ba)),
	}
}

func compositeDissolve(bottom, top color.Color) color.Color {
	r, g, b, a := top.RGBA()
	if uint32(rand.Int31n(0xffff)) < a {
		return color.RGBA64{
			R: uint16(r),
			G: uint16(g),
			B: uint16(b),
			A: 0xffff,
		}
	}
	return bottom
}

func compositeMultiply(x, y uint32) uint32 {
	return x * y / 0xffff
}

func compositeScreen(x, y uint32) uint32 {
	return 0xffff - (0xffff-x)*(0xffff-y)/0xffff
}

func compositeOverlay(x, y uint32) uint32 {
	t := 0xffff - y
	return (0xffff-y)*(x*x/0xffff)/0xffff + y*(0xffff-(t*t/0xffff))/0xffff
}

func compositeDifference(x, y uint32) uint32 {
	if x > y {
		return x - y
	}
	return y - x
}

func compositeAddition(x, y uint32) uint32 {
	return clamp(x + y)
}

func compositeSubtract(x, y uint32) uint32 {
	if y > x {
		return 0
	}
	return x - y
}

func compositeDarkenOnly(x, y uint32) uint32 {
	return min(x, y)
}

func compositeLightenOnly(x, y uint32) uint32 {
	return max(x, y)
}

func compositeDivide(x, y uint32) uint32 {
	if y == 0 {
		if x == 0 {
			return 0
		}
		return 0xffff
	}
	return clamp(0xffff * x / y)
}

func compositeDodge(x, y uint32) uint32 {
	if y == 0xffff {
		if x == 0 {
			return 0
		}
		return 0xffff
	}
	return clamp(0xffff * x / (0xffff - y))
}

func compositeBurn(x, y uint32) uint32 {
	if y == 0 {
		if x == 0xffff {
			return 0
		}
		return 0xffff
	}
	return clamp(0xffff - 0xffff*(0xffff-x)/y)
}

func compositeHardLight(x, y uint32) uint32 {
	if y < 0x7fff {
		return x * (y << 1) / 0xffff
	}
	return 0xffff - (0xffff-x)*(0x1fffe-y<<1)/0xffff
}

func compositeSoftLight(x, y uint32) uint32 {
	return compositeOverlay(x, y)
}

func compositeGrainExtract(x, y uint32) uint32 {
	if x+0x7fff < y {
		return 0
	}
	return clamp(x - y + 0x7fff)
}

func compositeGrainMerge(x, y uint32) uint32 {
	return clamp(x + y + 0x7fff)
}

func compositeHue(bottom, top color.Color) color.Color {
	br, bg, bb, _ := top.RGBA()
	if br == bg && br == bb {
		return bottom
	}
	a := rgbToHSV(bottom)
	b := rgbToHSV(top)
	a.H = b.H
	return a
}

func compositeSaturation(bottom, top color.Color) color.Color {
	a := rgbToHSV(bottom)
	b := rgbToHSV(top)
	a.S = b.S
	return a
}

func compositeColor(bottom, top color.Color) color.Color {
	a := rgbToHSL(bottom)
	b := rgbToHSL(top)
	b.L = a.L
	return b
}

func compositeValue(bottom, top color.Color) color.Color {
	a := rgbToHSV(bottom)
	b := rgbToHSV(top)
	a.V = b.V
	return a
}

func min(n ...uint32) uint32 {
	var m uint32 = 0xffffffff
	for _, o := range n {
		if o < m {
			m = o
		}
	}
	return m
}

func max(n ...uint32) uint32 {
	var m uint32
	for _, o := range n {
		if o > m {
			m = o
		}
	}
	return m
}

func mid(n ...uint32) uint32 {
	return (min(n...) + max(n...)) >> 1
}

func clamp(n uint32) uint32 {
	if n > 0xffff {
		return 0xffff
	}
	return n
}

func blend(a1, x1, a2, x2 uint32) uint32 {
	k := a2 / (0xffff - (0xffff-a1)*(0xffff-a2))
	return (1-k)*x1 + k*x2
}
