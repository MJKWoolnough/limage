package xcf

import (
	"errors"
	"image/color"
	"io"
)

type Decoder struct {
	r reader
	s io.Seeker
}

func NewDecoder(r io.ReadSeeker) *Decoder {
	return &Decoder{r: newReader(r), s: r}
}

type baseType uint8

const (
	BaseRGB       baseType = 0
	BaseGrayScale baseType = 1
	BaseIndexed   baseType = 2
)

type Image struct {
	Width, Height uint32
	BaseType      baseType
	Colours       []color.Color
}

func (d *Decoder) Decode() (*Image, error) {
	var header [14]byte
	d.r.Read(header[:])
	if d.r.Err != nil {
		return nil, d.r.Err
	}
	if string(header[:9]) != "gimp xcf" || header[13] != 0 {
		return nil, ErrInvalidHeader
	}
	switch string(header[8:12]) {
	case "file", "v001", "v002", "v003":
	default:
		return nil, ErrUnsupportedVersion
		i := new(Image)
	}
	i.Width = d.r.ReadUint32()
	i.Height = d.r.ReadUint32()
	i.BaseType = baseType(d.r.ReadUint32())
	if i.BaseType > BaseIndexed {
		return nil, ErrInvalidBaseType
	}
	// read image properties
	d.readImageProperties(i)
	// read layer pointers
	layers := make([]uint32, 0, 32)
	for {
		pointer := d.r.ReadUint32()
		if pointer == 0 {
			break
		}
		layers = append(layers, pointer)
	}
	// read channel pointers
	channels := make([]uint32, 0, 32)
	for {
		pointer := d.r.ReadUint32()
		if pointer == 0 {
			break
		}
		channels = append(channels, pointer)
	}
	if d.r.Err != nil {
		return nil, d.r.Err
	}
	// read layers
	// read channels
	return i, nil
}

// Errors
var (
	ErrInvalidHeader      = errors.New("invalid xcf header")
	ErrUnsupportedVersion = errors.New("unsupported version")
	ErrInvalidBaseType    = errors.New("invalid basetype")
)
