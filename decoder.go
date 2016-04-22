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
	tatoo         tatoo
	parasites     []parasite
	unit          unit
	paths         []path
	userUnit      userUnit
	vectors       vectors
}

type Decoder struct {
	r reader
	s io.Seeker
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

type Layer struct {
}

type Image struct {
	Width, Height uint32
	Layers        []Layer
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
	switch string(header[8:12]) {
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
	if i.BaseType > BaseIndexed {
		return nil, ErrInvalidBaseType
	}
	// read image properties
	d.readImageProperties()
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

type layer struct {
	offsetX, offsetY                            int32
	width, height                               uint32
	name                                        string
	alpha                                       bool
	editMask, showMask, visible, locked, active bool
}

func (d *Decoder) readLayer() layer {
	var l layer
	l.width = d.r.ReadUint32()
	l.height = d.r.ReadUint32()
	typ := d.r.ReadUint32()
	l.name = d.r.ReadString()

Props:
	for {
		propID := property(d.r.ReadUint32())
		propLength := d.r.ReadUint32()
		switch propID {
		case propEnd:
			break Props
		case propActiveLayer:
			l.active = true
		case propFloatingSelection:
			f := d.r.ReadUint32()
			_ = f
		case propOpacity:
			o := d.readOpacity()
			_ = o
		case propApplyMask:
			a := d.readBool()
			_ = a
		case propEditMask:
			l.editMask = d.readBool()
			_ = e
		case propMode:
			m := d.readMode()
			_ = m
		case propLinked:
			l := d.readBool()
			_ = l
		case propLockAlpha:
			l := d.readBool()
			_ = l
		case propOffsets:
			l.offsetX = d.r.ReadInt32()
			l.offsetY = d.r.ReadInt32()
		case propShowMask:
			l.showMask = d.readBool()
			_ = s
		case propTattoo:
			t := d.readTattoo()
			_ = t
		case propParasites:
			p := d.readParasites(propLength)
			_ = p
		case propTextLayerFlags:
			t := d.readTextLayerFlags()
			_ = t
		case propLockContent:
			l.locked = d.readBool()
		case propVisible:
			l.visible = d.readBool()
		case propGroupItem:
			// g := d.readGroupItem()
			// no data, just set as item group
		case propItemPath:
			i := d.readItemPath(propLength)
			_ = i
		case propGroupItemFlags:
			g := d.r.ReadUint32()
			_ = g
		default:
			d.s.Seek(int64(propLength), os.SEEK_CUR)
		}
	}

	hptr := d.r.ReadUint32()
	mptr := d.r.ReadUint32()
	switch typ {
	case 0:
		//RGB
	case 1:
		//RGBA
	case 2:
		//Y
	case 3:
		//YA
	case 4:
		//I
	case 5:
		//IA
	default:
		d.r.Err = ErrInvalidState
		return
	}
}

type channel struct{}

func (d *Decoder) readChannel() channel {
	width := d.r.ReadUint32()
	height := d.r.ReadUint32()
	name := d.r.ReadString()
	d.readChannelProperties()
	hptr := d.r.ReadUint32() //
}

type hierarchy struct{}

func (d *Decoder) readHierarchy() hierarchy {
	width := d.r.ReadUint32()
	height := d.r.ReadUint32()
	bpp := d.r.ReadUint32()
	lptr := d.r.ReadUint32()
	for {
		if d.r.ReadUint32() == 0 {
			break
		}
	}
}

type level struct{}

func (d *Decoder) readLevel() level {
	width := d.r.ReadUint32()
	height := d.r.ReadUint32()
	for {
		if d.r.ReadUint32() == 0 {
			break
		}
	}
}

func (d *Decoder) readTile(i draw.Image) {
	b := i.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {

		}
	}
}

// Errors
var (
	ErrInvalidHeader      = errors.New("invalid xcf header")
	ErrUnsupportedVersion = errors.New("unsupported version")
	ErrInvalidBaseType    = errors.New("invalid basetype")
)
