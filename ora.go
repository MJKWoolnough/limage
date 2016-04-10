package ora

import (
	"archive/zip"
	"encoding/xml"
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
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	required := 0
	files := make(map[string]*zip.File)
	var stack imageContent
	for _, f := range zr.File {
		switch f.Name {
		case "stack.xml":
			ms, err := f.Open()
			if err != nil {
				return nil, err
			}
			err = xml.NewDecoder(ms).Decode(&stack)
			ms.Close()
			if err != nil {
				return nil, err
			}
			required++
		case "mimetype":
			if !checkMime(f) {
				return nil, ErrInvalidMimeType
			}
			required++
		case "data", "Thumbnails/thumbnail.png", "mergedimage.png":
			required++
		default:
			files[f.Name] = f
		}
	}
	if required < 5 {
		return nil, ErrMissingRequired
	}

	return nil, nil
}

func EncodeStack(w io.Writer, s *Image) error {
	return nil
}
