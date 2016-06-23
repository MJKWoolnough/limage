package xcf

import "image/color"

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

func (c Composite) Composite(a, b color.Color) color.Color {
	switch c {
	case CompositeNormal:
		return compositeNormal(a, b)
	case CompositeDissolve:
		//return compositeDissolve(a, b)
	case CompositeBehind:
		//return compositeBehin(a, b)
	case CompositeMultiply:
		//return compositeMultiple(a, b)
	case CompositeScreen:
		//return compositeScreen(a, b)
	case CompositeOverlay:
		//return compositeOverlay(a, b)
	case CompositeDifference:
		//return compositeDifference(a, b)
	case CompositeAddition:
		//return compositeAddition(a, b)
	case CompositeSubtract:
		//return compositeSubtract(a, b)
	case CompositeDarkenOnly:
		//return compositeDarkenOnly(a, b)
	case CompositeLightenOnly:
		//return compositeLightenOnly(a, b)
	case CompositeHue:
		//return compositeHue(a, b)
	case CompositeSaturation:
		//return compositeSaturation(a, b)
	case CompositeColor:
		//return compositeColor(a, b)
	case CompositeValue:
		//return compositeValue(a, b)
	case CompositeDivide:
		//return compositeDivide(a, b)
	case CompositeDodge:
		//return compositeDodge(a, b)
	case CompositeBurn:
		//return compositeBurn(a, b)
	case CompositeHardLight:
		//return compositeHardLight(a, b)
	case CompositeSoftLight:
		//return compositeSoftLight(a, b)
	case CompositeGrainExtract:
		//return compositeGrainExtract(a, b)
	case CompositeGrainMerge:
		//return compositeGrainMerge(a, b)
	}
	return color.Alpha{}
}

func compositeNormal(a, b color.Color) color.Color {

	return nil
}

func min(n ...uint32) uint32 {
	var m uint32 = 0xffffffff
	for _, o := range n {
		if o < m {
			m = o
		}
	}
}

func max(n ...uint32) uint32 {
	var m uint32
	for _, o := range n {
		if o > m {
			m = o
		}
	}
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
	k := a2 / (1 - (1-a1)*(1-a2))
	return (1-k)*x1 + k*x2
}
