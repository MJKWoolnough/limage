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
		//return compositeDifference(bottom, top)
	case CompositeAddition:
		//return compositeAddition(bottom, top)
	case CompositeSubtract:
		//return compositeSubtract(bottom, top)
	case CompositeDarkenOnly:
		//return compositeDarkenOnly(bottom, top)
	case CompositeLightenOnly:
		//return compositeLightenOnly(bottom, top)
	case CompositeHue:
		//return compositeHue(bottom, top)
	case CompositeSaturation:
		//return compositeSaturation(bottom, top)
	case CompositeColor:
		//return compositeColor(bottom, top)
	case CompositeValue:
		//return compositeValue(bottom, top)
	case CompositeDivide:
		//return compositeDivide(bottom, top)
	case CompositeDodge:
		//return compositeDodge(bottom, top)
	case CompositeBurn:
		//return compositeBurn(bottom, top)
	case CompositeHardLight:
		//return compositeHardLight(bottom, top)
	case CompositeSoftLight:
		//return compositeSoftLight(bottom, top)
	case CompositeGrainExtract:
		//return compositeGrainExtract(bottom, top)
	case CompositeGrainMerge:
		//return compositeGrainMerge(bottom, top)
	default: //Normal
		return compositeNormal(bottom, top)
	}
	ar, ag, ab, aa := bottom.RGBA()
	br, bg, bb, ba := top.RGBA()
	return color.RGBA64{
		R: uint16(blend(aa, br, min(aa, ba), f(ar, br))),
		G: uint16(blend(aa, bg, min(aa, ba), f(ag, bg))),
		B: uint16(blend(aa, bb, min(aa, ba), f(ab, bb))),
		A: aa,
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
