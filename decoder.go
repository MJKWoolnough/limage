package xcf

import (
	"errors"
	"image/color"
	"image/draw"
	"io"
	"os"
)

type props struct {
	width, height uint32
	baseType      baseType
	colours       color.Palette
	compression   compression
	guides        []guide
	hres, vres    float32
	tattoo        tattoo
	parasites     []parasite
	unit          unit
	paths         paths
	userUnit      userUnit
	vectors       vectors
}

type Decoder struct {
	r reader
	props
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
	Layers        []Layer
	channels      []Channel
}

func (d *Decoder) Decode() (*Image, error) {
	var header [14]byte
	d.r.Read(header[:])
	if d.r.Err != nil {
		return nil, d.r.Err
	}
	if string(header[:9]) != "gimp xcf " || header[13] != 0 {
		return nil, ErrInvalidHeader
	}
	switch string(header[9:13]) {
	case "file", "v001", "v002", "v003":
	default:
		return nil, ErrUnsupportedVersion
	}
	d.props = props{}
	i := new(Image)
	i.Width = d.r.ReadUint32()
	i.Height = d.r.ReadUint32()
	d.width = i.Width
	d.height = i.Height
	d.baseType = baseType(d.r.ReadUint32())
	if d.baseType > BaseIndexed {
		return nil, ErrInvalidBaseType
	}
	// read image properties
Props:
	for {
		propID := property(d.r.ReadUint32())
		propLength := d.r.ReadUint32()
		switch propID {
		case propEnd:
			break Props
		case propColormap:
			d.colours = d.readColorMap()
		case propCompression:
			d.compression = d.readCompression()
		case propGuides:
			d.guides = d.readGuides(propLength)
		case propResolution:
			d.hres = d.r.ReadFloat32()
			d.vres = d.r.ReadFloat32()
		case propTattoo:
			d.tattoo = d.readTattoo()
		case propParasites:
			d.parasites = d.readParasites(propLength)
		case propUnit:
			d.unit = d.readUnit()
		case propPaths:
			d.paths = d.readPaths()
		case propUserUnit:
			d.userUnit = d.readUserUnit()
		case propVectors:
			d.vectors = d.readVectors()
		default:
			d.r.Seek(int64(propLength), os.SEEK_CUR)
		}
	}
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
	i.Layers = make([]Layer, len(layers))
	for n, ptr := range layers {
		d.r.Seek(int64(ptr), os.SEEK_SET)
		i.Layers[n] = d.readLayer()
	}
	// read channels
	i.channels = make([]Channel, len(channels))
	for n, ptr := range channels {
		d.r.Seek(int64(ptr), os.SEEK_SET)
		i.channels[n] = d.readChannel()
	}
	return i, d.r.Err
}

type hierarchy struct{}

func (d *Decoder) readHierarchy() hierarchy {
	width := d.r.ReadUint32()
	height := d.r.ReadUint32()
	bpp := d.r.ReadUint32()
	lptr := d.r.ReadUint32()
	_, _, _, _ = width, height, bpp, lptr
	for {
		if d.r.ReadUint32() == 0 {
			break
		}
	}
	return hierarchy{}
}

type level struct{}

func (d *Decoder) readLevel() level {
	width := d.r.ReadUint32()
	height := d.r.ReadUint32()
	_, _ = width, height
	for {
		if d.r.ReadUint32() == 0 {
			break
		}
	}
	return level{}
}

func (d *Decoder) readTile(i draw.Image, alpha bool) {
	b := i.Bounds()
	var r byteReader
	switch d.compression {
	case 0:
		r = &d.r
	case 1:
		r = &rle{r: &d.r.StickyReader}
	}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			i.Set(x, y, d.readColor(r, alpha))
		}
	}
}

func (d *Decoder) readColor(reader byteReader, alpha bool) color.NGBA {
	var (
		r, g, b uint8
		a       uint8 = 255
	)
	switch d.props.baseType {
	case BaseRGB:
		r = reader.ReadByte()
		g = reader.ReadByte()
		b = reader.ReadByte()
	case BaseGrayScale:
		r := reader.ReadByte()
		g, b = r, r
	case BaseIndexed:
		i := reader.ReadByte()
		if int(i) >= len(d.colours) {
			i = 0
		}
		c := d.colours[i]
		dr, dg, db, _ := c.RGBA()
		r = uint8(dr >> 8)
		g = uint8(dg >> 8)
		b = uint8(db >> 8)
	}
	if alpha {
		a = reader.ReadByte()
	}
	return color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}

type byteReader interface {
	ReadByte() byte
}

// Errors
var (
	ErrInvalidHeader      = errors.New("invalid xcf header")
	ErrUnsupportedVersion = errors.New("unsupported version")
	ErrInvalidBaseType    = errors.New("invalid basetype")
)
