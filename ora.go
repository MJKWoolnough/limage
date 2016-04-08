package ora

import (
	"image/color"
	"io"
)

type composite uint8

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

type Image struct {
	Width, Height int
	Name          string
	Stack
}

type Stack struct {
	X, Y    int
	Name    string
	Content []Content
}

type Content interface{}

type Layer struct {
	X, Y        int
	Name        string
	CompositeOp composite
	Opacity     float32
	//Filters     []Filter // Not needed for baseline
}

type Text struct {
	X, Y  int
	Name  string
	Data  string
	Font  string
	Size  uint16
	Color color.Color
}

func DecodeStack(r io.ReaderAt, size int64) (*Image, error) {
	return nil, nil
}

func EncodeStack(w io.Writer, s *Image) error {
	return nil
}
