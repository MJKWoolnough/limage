package limage

import (
	"image/color"
	"math/rand"

	"vimagination.zapto.org/limage/internal"
	"vimagination.zapto.org/limage/lcolor"
)

// Composite determines how two layers are composed together
type Composite uint32

// Composite constants
const (
	CompositeNormal Composite = iota
	CompositeDissolve
	CompositeBehind
	CompositeMultiply
	CompositeScreen
	CompositeOverlay
	CompositeDifference
	CompositeAddition
	CompositeSubtract
	CompositeDarkenOnly
	CompositeLightenOnly
	CompositeHue
	CompositeSaturation
	CompositeColor
	CompositeValue
	CompositeDivide
	CompositeDodge
	CompositeBurn
	CompositeHardLight
	CompositeSoftLight
	CompositeGrainExtract
	CompositeGrainMerge
	CompositeLuminosity
	CompositePlus
	CompositeDestinationIn
	CompositeDestinationOut
	CompositeSourceAtop
	CompositeDestinationAtop
	CompositeColorErase
	CompositeChroma
	CompositeLightness
	CompositeVividLight
	CompositePinLight
	CompositeLinearLight
	CompositeHardMix
	CompositeExclusion
	CompositeLinearBurn
	CompositeLuminance
	CompositeErase
	CompositeMerge
	CompositeSplit
	CompositePassThrough
)

var compositeNames = [...]string{
	"Normal",
	"Dissolve",
	"Behind",
	"Multiply",
	"Screen",
	"Overlay",
	"Difference",
	"Addition",
	"Subtract",
	"Darken Only",
	"Lighten Only",
	"Hue",
	"Saturation",
	"Color",
	"Value",
	"Divide",
	"Dodge",
	"Burn",
	"Hard Light",
	"Soft Light",
	"Grain Extract",
	"Grain Merge",
	"Luminosity",
	"Plus",
	"Destination In",
	"Destination Out",
	"Source Atop",
	"Destination Atop",
}

// String returns the name of the composition
func (c Composite) String() string {
	if int(c) < len(compositeNames) {
		return compositeNames[c]
	}
	return compositeNames[0]
}

// Composite performs the composition of two layers
func (c Composite) Composite(b, t color.Color) color.Color {
	bottom := internal.ColourToNRGBA(b)
	top := internal.ColourToNRGBA(t)
	var f func(uint16, uint16) uint16
	switch c {
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
	case CompositeBehind:
		return bottom
	case CompositeDissolve:
		return compositeDissolve(bottom, top)
	case CompositeHue:
		return compositeHue(bottom, top)
	case CompositeSaturation:
		return compositeSaturation(bottom, top)
	case CompositeColor:
		return compositeColor(bottom, top)
	case CompositeValue:
		return compositeValue(bottom, top)
	case CompositeLuminosity:
		return compositeColor(top, bottom)
	case CompositePlus:
		return compositePlus(bottom, top)
	case CompositeDestinationIn:
		return compositeDstIn(bottom, top)
	case CompositeDestinationOut:
		return compositeDstOut(bottom, top)
	case CompositeSourceAtop:
		return compositeAtop(bottom, top)
	case CompositeDestinationAtop:
		return compositeAtop(top, bottom)
	default: //Normal
		return compositeNormal(bottom, top)
	}
	if bottom.A == 0 {
		return color.NRGBA{}
	}
	ma := internal.Min(bottom.A, top.A)
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
	return internal.Min(x, y)
}

func compositeLightenOnly(x, y uint16) uint16 {
	return internal.Max(x, y)
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
	a := lcolor.RGBToHSV(bottom)
	b := lcolor.RGBToHSV(top)
	a.H = b.H
	return a
}

func compositeSaturation(bottom, top color.NRGBA64) color.Color {
	a := lcolor.RGBToHSV(bottom)
	b := lcolor.RGBToHSV(top)
	a.S = b.S
	return a
}

func compositeColor(bottom, top color.NRGBA64) color.Color {
	a := lcolor.RGBToHSL(bottom)
	b := lcolor.RGBToHSL(top)
	b.L = a.L
	return b
}

func compositeValue(bottom, top color.NRGBA64) color.Color {
	a := lcolor.RGBToHSV(bottom)
	b := lcolor.RGBToHSV(top)
	a.V = b.V
	return a
}

func compositePlus(bottom, top color.NRGBA64) color.NRGBA64 {
	return compositeNormal(bottom, top)
}

func compositeDstIn(bottom, top color.NRGBA64) color.NRGBA64 {
	return compositeNormal(bottom, top)
}

func compositeDstOut(bottom, top color.NRGBA64) color.NRGBA64 {
	return compositeNormal(bottom, top)
}

func compositeAtop(bottom, top color.NRGBA64) color.NRGBA64 {
	return compositeNormal(bottom, top)
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
