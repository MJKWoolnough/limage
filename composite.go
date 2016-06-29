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

func (c Composite) String() string {
	switch c {
	case CompositeDissolve:
		return "Dissolve"
	case CompositeBehind:
		return "Behind"
	case CompositeMultiply:
		return "Multiply"
	case CompositeScreen:
		return "Screen"
	case CompositeOverlay:
		return "Overlay"
	case CompositeDifference:
		return "Difference"
	case CompositeAddition:
		return "Addition"
	case CompositeSubtract:
		return "Subtract"
	case CompositeDarkenOnly:
		return "Darken Only"
	case CompositeLightenOnly:
		return "Lighten Only"
	case CompositeHue:
		return "Hue"
	case CompositeSaturation:
		return "Saturation"
	case CompositeColor:
		return "Color"
	case CompositeValue:
		return "Value"
	case CompositeDivide:
		return "Divide"
	case CompositeDodge:
		return "Dodge"
	case CompositeBurn:
		return "Burn"
	case CompositeHardLight:
		return "Hard Light"
	case CompositeSoftLight:
		return "Soft Light"
	case CompositeGrainExtract:
		return "Grain Extract"
	case CompositeGrainMerge:
		return "Grain Merge"
	default:
		return "Normal"
	}
}

func (c Composite) Composite(bottom, top color.NRGBA64) color.Color {
	var f func(uint16, uint16) uint16
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
	if bottom.A == 0 {
		return color.NRGBA{}
	}
	ma := min(bottom.A, top.A)
	return color.NRGBA64{
		R: blend(bottom.A, bottom.R, ma, f(bottom.R, top.R)),
		G: blend(bottom.A, bottom.G, ma, f(bottom.R, top.G)),
		B: blend(bottom.A, bottom.B, ma, f(bottom.R, top.B)),
		A: bottom.A,
	}
}

func compositeNormal(bottom, top color.NRGBA64) color.NRGBA64 {
	if bottom.A == 0 && top.A == 0 {
		return color.NRGBA64{}
	}
	return color.NRGBA64{
		R: blend(bottom.A, bottom.R, top.A, top.R),
		G: blend(bottom.A, bottom.G, top.A, top.G),
		B: blend(bottom.A, bottom.B, top.A, top.B),
		A: uint16(0xffff - (0xffff-uint32(bottom.A))*(0xffff-uint32(top.A))/0xffff),
	}
}

func compositeDissolve(bottom, top color.NRGBA64) color.NRGBA64 {
	if uint16(rand.Int31n(0xffff)) < bottom.A {
		top.A = 0xffff
		return top
	}
	return bottom
}

func compositeMultiply(x, y uint16) uint16 {
	return uint16(uint32(x) * uint32(y) / 0xffff)
}

func compositeScreen(x, y uint16) uint16 {
	return uint16(0xffff - uint32(0xffff-x)*uint32(0xffff-y)/0xffff)
}

func compositeOverlay(ax, ay uint16) uint16 {
	x := uint32(ax)
	y := uint32(ay)
	t := 0xffff - y
	return uint16((0xffff-y)*(x*x/0xffff)/0xffff + y*(0xffff-(t*t/0xffff))/0xffff)
}

func compositeDifference(x, y uint16) uint16 {
	if x > y {
		return x - y
	}
	return y - x
}

func compositeAddition(x, y uint16) uint16 {
	return clamp(uint32(x) + uint32(y))
}

func compositeSubtract(x, y uint16) uint16 {
	if y > x {
		return 0
	}
	return x - y
}

func compositeDarkenOnly(x, y uint16) uint16 {
	return min(x, y)
}

func compositeLightenOnly(x, y uint16) uint16 {
	return max(x, y)
}

func compositeDivide(x, y uint16) uint16 {
	if y == 0 {
		if x == 0 {
			return 0
		}
		return 0xffff
	}
	return clamp(0xffff * uint32(x) / uint32(y))
}

func compositeDodge(x, y uint16) uint16 {
	if y == 0xffff {
		if x == 0 {
			return 0
		}
		return 0xffff
	}
	return clamp(0xffff * uint32(x) / (0xffff - uint32(y)))
}

func compositeBurn(x, y uint16) uint16 {
	if y == 0 {
		if x == 0xffff {
			return 0
		}
		return 0xffff
	}
	return clamp(0xffff - 0xffff*(0xffff-uint32(x))/uint32(y))
}

func compositeHardLight(x, y uint16) uint16 {
	if y < 0x7fff {
		return uint16(uint32(x) * (uint32(y) << 1) / 0xffff)
	}
	return uint16(0xffff - (0xffff-uint32(x))*(0x1fffe-uint32(y)<<1)/0xffff)
}

func compositeSoftLight(x, y uint16) uint16 {
	return compositeOverlay(x, y)
}

func compositeGrainExtract(ax, ay uint16) uint16 {
	x := uint32(ax)
	y := uint32(ay)
	if x+0x7fff < y {
		return 0
	}
	return clamp(x - y + 0x7fff)
}

func compositeGrainMerge(x, y uint16) uint16 {
	return clamp(uint32(x) + uint32(y) + 0x7fff)
}

func compositeHue(bottom, top color.NRGBA64) color.Color {
	br, bg, bb, _ := top.RGBA()
	if br == bg && br == bb {
		return bottom
	}
	a := rgbToHSV(bottom)
	b := rgbToHSV(top)
	a.H = b.H
	return a
}

func compositeSaturation(bottom, top color.NRGBA64) color.Color {
	a := rgbToHSV(bottom)
	b := rgbToHSV(top)
	a.S = b.S
	return a
}

func compositeColor(bottom, top color.NRGBA64) color.Color {
	a := rgbToHSL(bottom)
	b := rgbToHSL(top)
	b.L = a.L
	return b
}

func compositeValue(bottom, top color.NRGBA64) color.Color {
	a := rgbToHSV(bottom)
	b := rgbToHSV(top)
	a.V = b.V
	return a
}

func min(n ...uint16) uint16 {
	var m uint16 = 0xffff
	for _, o := range n {
		if o < m {
			m = o
		}
	}
	return m
}

func max(n ...uint16) uint16 {
	var m uint16
	for _, o := range n {
		if o > m {
			m = o
		}
	}
	return m
}

func mid(n ...uint16) uint16 {
	return (min(n...) + max(n...)) >> 1
}

func clamp(n uint32) uint16 {
	if n > 0xffff {
		return 0xffff
	}
	return uint16(n)
}

func blend(aa1, ax1, aa2, ax2 uint16) uint16 {
	a1 := uint32(aa1)
	x1 := uint32(ax1)
	a2 := uint32(aa2)
	x2 := uint32(ax2)
	k := 0xffff * a2 / (0xffff - (0xffff-a1)*(0xffff-a2)/0xffff)
	return uint16((0xffff-k)*x1/0xffff + k*x2/0xffff)
}
