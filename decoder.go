package xcf

import (
	"errors"
	"image"
	"image/color"
	"io"
	"os"
)

const (
	fileTypeID   = "gimp xcf "
	fileVersion0 = "file"
	fileVersion1 = "v001"
	fileVersion2 = "v002"
	fileVersion3 = "v003"
)

const (
	baseRGB     = 0
	baseGrey    = 1
	baseIndexed = 2
)

type decoder struct {
	reader
	width, height, baseType      uint32
	linked, lockContent, visible bool
	opacity                      uint8
	parasites                    parasites
	tattoo                       uint32
	palette                      color.Palette
	compression                  uint8
	guides                       []guide
	paths                        paths
	hres, vres                   float32
	samplePoints                 []samplePoint
	unit                         uint32
	userUnit                     struct {
		factor                           float32
		digits                           uint32
		id, symbol, abbrev, sname, pname string
	}
	vectors  vectors
	layers   []layer
	channels []channel
}

type guide struct {
	coord int32
	hv    bool
}

type samplePoint struct {
	x, y uint32
}

func Decode(r io.ReadSeeker) (image.Image, error) {
	d := decoder{reader: newReader(r)}

	// check header

	var header [15]byte
	d.Read(header[:])
	if d.Err != nil {
		return nil, d.Err // wrap?
	}
	if string(header[:9]) != fileTypeID {
		return nil, ErrInvalidFileTypeID
	}
	switch string(header[9:13]) {
	case fileVersion0, fileVersion1, fileVersion2, fileVersion3:
	default:
		return nil, ErrUnsupportedVersion
	}
	if header[14] != 0 {
		return nil, ErrInvalidHeader
	}
	d.width = d.ReadUint32()
	d.height = d.ReadUint32()
	d.baseType = d.ReadUint32()

	// read image properties
PropertyLoop:
	for {
		typ := d.ReadUint32()
		plength := d.ReadUint32()
		switch typ {
		case propEnd:
			if plength != 0 {
				return nil, ErrInvalidProperties
			}
			break PropertyLoop

		//general properties
		case propLinked:
			d.linked = d.ReadBoolProperty()
		case propLockContent:
			d.lockContent = d.ReadBoolProperty()
		case propOpacity:
			o := d.ReadUint32()
			if o > 255 {
				return nil, ErrInvalidOpacity
			}
			d.opacity = uint8(o)
		case propParasites:
			d.parasites = d.ReadParasites(plength)
		case propTattoo:
			d.tattoo = d.ReadUint32()
		case propVisible:
			d.visible = d.ReadBoolProperty()

		// image properties
		case propColorMap:
			if d.baseType != baseIndexed {
				d.Seek(int64(plength), os.SEEK_CUR) // skip
			}
			numColours := d.ReadUint32()
			d.palette = make(color.Palette, numColours)
			for i := uint32(0); i < numColours; i++ {
				r := d.ReadUint8()
				g := d.ReadUint8()
				b := d.ReadUint8()
				d.palette[i] = rgb{
					R: r,
					G: g,
					B: b,
				}
			}
		case propCompression:
			d.compression = d.ReadUint8()
			if d.compression > 1 {
				return nil, ErrUnknownCompression
			}
		case propGuides:
			ng := plength / 5
			if ng*5 != plength {
				return nil, ErrInvalidGuideLength
			}
			d.guides = make([]guide, ng)
			for n := range d.guides {
				d.guides[n].coord = d.ReadInt32()
				d.guides[n].hv = d.ReadBoolProperty()
			}
		case propPaths:
			d.paths = d.ReadPaths()
		case propResolution:
			d.hres = d.ReadFloat32()
			d.vres = d.ReadFloat32()
		case propSamplePoints:
			if plength&1 == 1 {
				return nil, ErrInvalidSampleLength
			}
			d.samplePoints = make([]samplePoint, plength>>1)
			for i := uint32(0); i < plength>>1; i++ {
				d.samplePoints[i].x = d.ReadUint32()
				d.samplePoints[i].y = d.ReadUint32()
			}
		case propUnit:
			d.unit = d.ReadUint32()
			if d.unit < 0 || d.unit > 4 {
				return nil, ErrInvalidUnit
			}
		case propUserUnit:
			d.userUnit.factor = d.ReadFloat32()
			d.userUnit.digits = d.ReadUint32()
			d.userUnit.id = d.ReadString()
			d.userUnit.symbol = d.ReadString()
			d.userUnit.abbrev = d.ReadString()
			d.userUnit.sname = d.ReadString()
			d.userUnit.pname = d.ReadString()
		case propVectors:
			d.vectors = d.ReadVectors()
		default:
			d.Seek(int64(plength), os.SEEK_CUR)
		}
	}

	if d.Err != nil {
		return nil, d.Err
	}
	layerptrs := make([]uint32, 0, 32)
	for {
		lptr := d.ReadUint32()
		if lptr == 0 {
			break
		}
		layerptrs = append(layerptrs, lptr)
	}
	channelptrs := make([]uint32, 0, 32)
	for {
		cptr := d.ReadUint32()
		if cptr == 0 {
			break
		}
		channelptrs = append(channelptrs, cptr)
	}

	d.layers = make([]layer, len(layerptrs))

	for i := range d.layers {
		d.Seek(int64(layerptrs[i]), os.SEEK_SET)
		d.layers[i] = d.ReadLayer()
	}

	d.channels = make([]channel, len(channelptrs))

	for i := range d.channels {
		d.Seek(int64(channelptrs[i]), os.SEEK_SET)
		d.channels[i] = d.ReadChannel()
	}

	return nil, nil
}

func (d *decoder) SetError(err error) {
	if d.Err == nil {
		d.Err = err
	}
}

// Errors
var (
	ErrInvalidFileTypeID   = errors.New("invalid file type identification")
	ErrUnsupportedVersion  = errors.New("unsupported file version")
	ErrInvalidHeader       = errors.New("invalid header")
	ErrInvalidProperties   = errors.New("invalid property list")
	ErrInvalidOpacity      = errors.New("opacity not in valid range")
	ErrInvalidGuideLength  = errors.New("invalid guide length")
	ErrInvalidUnit         = errors.New("invalid unit")
	ErrInvalidSampleLength = errors.New("invalid sample points length")
	ErrUnknownCompression  = errors.New("unknown compressio method")
)
