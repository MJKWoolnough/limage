package xcf

import "errors"

type layer struct {
	width, height uint32
}

func (d *decoder) ReadLayer() layer {
	var l layer
	l.width = d.ReadUint32()
	l.height = d.ReadUint32()
	typ := d.ReadUint()
	if typ>>1 != d.baseType {
		d.Err = ErrInvalidLayerType
		return l
	}
	for {

	}
}

// Errors
var (
	ErrInvalidLayerType = errors.New("invalid layer type")
)
