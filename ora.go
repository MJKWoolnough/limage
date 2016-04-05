package ora

import "io"

const (
	CompositeSrcOver composite = iota
	CompositeMultiply
	CompositeScreen
	CompositeOverlay
	CompositeDarken
	CompositeLighten
	CompositeColorDodge
	CompositeColorBurn
	CompositeHardLight
	CompositeSoftLight
	CompositeDifference
	CompositeColor
	CompositeLuminosity
	CompositeHue
	CompositeSaturation
	CompositePlus
	CompositeDstIn
	CompositeDstOut
	CompositeSrcAtop
	CompositeDstAtop
)

type Layers struct {
}

func DecodeLayers(r io.ReaderAt) (*Layers, error) {
	return nil, nil
}
