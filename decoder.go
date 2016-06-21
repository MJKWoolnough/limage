package xcf

import (
	"errors"
	"image"
	"image/color"
	"io"
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
	*Image
	reader
	baseType                     uint32
	linked, lockContent, visible bool
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
	vectors vectors
	//channels []channel
}

type guide struct {
	coord int32
	hv    bool
}

type samplePoint struct {
	x, y uint32
}

func Decode(r io.ReadSeeker) (image.Image, error) {
	d := decoder{Image: new(Image), reader: newReader(r)}

	// check header

	var header [14]byte
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
	if header[13] != 0 {
		return nil, ErrInvalidHeader
	}

	d.Width = int(d.ReadUint32())
	d.Height = int(d.ReadUint32())
	d.baseType = d.ReadUint32()
	switch d.baseType {
	case 0:
		d.Config.ColorModel = color.RGBAModel
	case 1:
		d.Config.ColorModel = color.ModelFunc(grayAlphaColourModel)
	}

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
			d.Opacity = uint8(o)
		case propParasites:
			d.parasites = d.ReadParasites(plength)
			if c := d.parasites.Get(commentParasiteName); c != nil && len(c.data) > 0 {
				d.Comment = string(c.data)
			}
		case propTattoo:
			d.tattoo = d.ReadUint32()
		case propVisible:
			d.visible = d.ReadBoolProperty()

		// image properties
		case propColorMap:
			if d.baseType != baseIndexed {
				d.Skip(plength) // skip
			}
			numColours := d.ReadUint32()
			d.palette = make(color.Palette, numColours)
			for i := uint32(0); i < numColours; i++ {
				r := d.ReadUint8()
				g := d.ReadUint8()
				b := d.ReadUint8()
				d.palette[i] = RGB{
					R: r,
					G: g,
					B: b,
				}
			}
			d.Image.Config.ColorModel = d.palette
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
			d.Skip(plength)
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

	/*
		channelptrs := make([]uint32, 0, 32)
		for {
			cptr := d.ReadUint32()
			if cptr == 0 {
				break
			}
			channelptrs = append(channelptrs, cptr)
		}

		if d.Err != nil {
			return nil, d.Err
		}

	*/

	type groupOffset struct {
		*Group
		OffsetX, OffsetY int
	}

	groups := make(map[string]groupOffset)
	groups[""] = groupOffset{Group: &d.Group}
	var n rune
	for _, lptr := range layerptrs {
		d.Goto(lptr)
		l := d.ReadLayer()
		if d.Err != nil {
			return nil, d.Err
		}
		if len(l.itemPath) == 0 {
			l.itemPath = []rune{n}
			n++
		}
		g := groups[string(l.itemPath[:len(l.itemPath)-1])]
		if g.Group == nil {
			return nil, ErrInvalidGroup
		}
		if l.group {
			gp := new(Group)
			gp.Width = int(l.width)
			gp.Height = int(l.height)
			gp.Config.ColorModel = d.Config.ColorModel
			l.Image = gp
			groups[string(l.itemPath)] = groupOffset{
				Group:   gp,
				OffsetX: l.OffsetX,
				OffsetY: l.OffsetY,
			}

		} else {
			if t := l.parasites.Get(textParasiteName); t != nil {
				textData, err := parseTextData(t)
				if err != nil {
					return nil, err
				}
				l.Image = &Text{
					Image:    l.Image,
					TextData: textData,
				}
			}
			if l.mask.image != nil {
				l.Image = &MaskedImage{
					Image: l.Image,
					Mask:  l.mask.image.(*image.Gray),
				}
			}
		}
		l.OffsetX -= g.OffsetX
		l.OffsetY -= g.OffsetY
		g.Layers = append(g.Layers, l.Layer)
	}

	/*
		d.channels = make([]channel, len(channelptrs))

		for i := range d.channels {
			d.Goto(channelptrs[i])
			d.channels[i] = d.ReadChannel()
			if d.Err != nil {
				return nil, d.Err
			}
		}
	*/

	return d.Image, nil
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
	ErrInvalidGroup        = errors.New("invalid or unknown group specified for layer")
	ErrUnknownCompression  = errors.New("unknown compression method")
)
