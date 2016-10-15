package xcf

import (
	"bufio"
	"image"
	"image/color"
	"io"
	"sync"

	"github.com/MJKWoolnough/errors"
	"github.com/MJKWoolnough/limage"
	"github.com/MJKWoolnough/limage/lcolor"
)

func getReaderAt(r io.Reader) io.ReaderAt {
	if bb, ok := r.(*bufio.Reader); ok {
		return bufioToReader(bb)
	}
	return nil
}

func decodeConfig(r io.Reader) (image.Config, error) {
	ra := getReaderAt(r)
	if ra == nil {
		return image.Config{}, errors.Error("need a io.ReaderAt")
	}
	return DecodeConfig(ra)
}

func decode(r io.Reader) (image.Image, error) {
	ra := getReaderAt(r)
	if ra == nil {
		return nil, errors.Error("need a io.ReaderAt")
	}
	return Decode(ra)
}

func init() {
	image.RegisterFormat("xcf", fileTypeID, decode, decodeConfig)
}

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
	Width, Height int
	baseType      uint32
	palette       lcolor.AlphaPalette
	compression   uint8
}

// DecodeConfig retrieves the color model and dimensions of the XCF image
func DecodeConfig(r io.ReaderAt) (image.Config, error) {
	var c image.Config

	dr := newReader(r)

	// check header

	var header [14]byte
	dr.Read(header[:])
	if dr.Err != nil {
		return c, dr.Err
	}
	if string(header[:9]) != fileTypeID {
		return c, ErrInvalidFileTypeID
	}
	switch string(header[9:13]) {
	case fileVersion0, fileVersion1, fileVersion2, fileVersion3:
	default:
		return c, ErrUnsupportedVersion
	}
	if header[13] != 0 {
		return c, ErrInvalidHeader
	}

	c.Width = int(dr.ReadUint32())
	c.Height = int(dr.ReadUint32())
	baseType := dr.ReadUint32()
	switch baseType {
	case 0:
		c.ColorModel = color.RGBAModel
	case 1:
		c.ColorModel = lcolor.GrayAlphaModel
	case 2:
	PropertyLoop:
		for {
			typ := dr.ReadUint32()
			plength := dr.ReadUint32()
			switch typ {
			case propEnd:
				if plength != 0 {
					return c, ErrInvalidProperties
				}
				break PropertyLoop

			// the one we care about
			case propColorMap:
				if baseType != baseIndexed {
					dr.Skip(plength) // skip
				}
				palette := make(lcolor.AlphaPalette, dr.ReadUint32())
				for n := range palette {
					r := dr.ReadUint8()
					g := dr.ReadUint8()
					b := dr.ReadUint8()
					palette[n] = lcolor.RGB{
						R: r,
						G: g,
						B: b,
					}
				}
				c.ColorModel = palette
				break PropertyLoop

			//general properties
			case propLinked:
				dr.ReadBoolProperty()
			case propLockContent:
				dr.ReadBoolProperty()
			case propOpacity:
				if o := dr.ReadUint32(); o > 255 {
					return c, ErrInvalidOpacity
				}
			case propParasites:
				dr.ReadParasites(plength)
			case propTattoo:
				dr.ReadUint32()
			case propVisible:
				dr.ReadBoolProperty()
			case propCompression:
				if dr.ReadUint8() > 1 {
					return c, ErrUnknownCompression
				}
			case propGuides:
				ng := plength / 5
				if ng*5 != plength {
					return c, ErrInvalidGuideLength
				}
				for n := uint32(0); n < ng; n++ {
					dr.ReadInt32()
					dr.ReadBoolProperty()
				}
			case propPaths:
				dr.ReadPaths()
			case propResolution:
				dr.ReadFloat32()
				dr.ReadFloat32()
			case propSamplePoints:
				if plength&1 == 1 {
					return c, ErrInvalidSampleLength
				}
				for i := uint32(0); i < plength>>1; i++ {
					dr.ReadUint32()
					dr.ReadUint32()
				}
			case propUnit:
				if unit := dr.ReadUint32(); unit < 0 || unit > 4 {
					return c, ErrInvalidUnit
				}
			case propUserUnit:
				dr.ReadFloat32()
				dr.ReadUint32()
				dr.ReadString()
				dr.ReadString()
				dr.ReadString()
				dr.ReadString()
				dr.ReadString()
			case propVectors:
				dr.ReadVectors()
			default:
				dr.Skip(plength)
			}
		}
	}

	return c, dr.Err
}

// Decode reads an XCF layered image from the given ReadSeeker
func Decode(r io.ReaderAt) (limage.Image, error) {
	dr := newReader(r)

	// check header

	var header [14]byte
	dr.Read(header[:])
	if dr.Err != nil {
		return nil, dr.Err // wrap?
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

	width := int(dr.ReadUint32())
	height := int(dr.ReadUint32())
	baseType := dr.ReadUint32()

	var (
		palette     lcolor.AlphaPalette
		compression uint8
	)

	// read image properties
PropertyLoop:
	for {
		typ := dr.ReadUint32()
		plength := dr.ReadUint32()
		switch typ {
		case propEnd:
			if plength != 0 {
				return nil, ErrInvalidProperties
			}
			break PropertyLoop

		//general properties
		case propLinked:
			dr.ReadBoolProperty()
		case propLockContent:
			dr.ReadBoolProperty()
		case propOpacity:
			o := dr.ReadUint32()
			if o > 255 {
				return nil, ErrInvalidOpacity
			}
		case propParasites:
			dr.ReadParasites(plength)
		case propTattoo:
			dr.ReadUint32()
		case propVisible:
			dr.ReadBoolProperty()

		// image properties
		case propColorMap:
			if baseType != baseIndexed {
				dr.Skip(plength) // skip
			}
			palette = make(lcolor.AlphaPalette, dr.ReadUint32())
			for n := range palette {
				r := dr.ReadUint8()
				g := dr.ReadUint8()
				b := dr.ReadUint8()
				palette[n] = lcolor.RGB{
					R: r,
					G: g,
					B: b,
				}
			}
		case propCompression:
			compression = dr.ReadUint8()
			if compression > 1 {
				return nil, ErrUnknownCompression
			}
		case propGuides:
			ng := plength / 5
			if ng*5 != plength {
				return nil, ErrInvalidGuideLength
			}
			for n := uint32(0); n < ng; n++ {
				dr.ReadInt32()
				dr.ReadBoolProperty()
			}
		case propPaths:
			dr.ReadPaths()
		case propResolution:
			dr.ReadFloat32() // x
			dr.ReadFloat32() // y
		case propSamplePoints:
			if plength&1 == 1 {
				return nil, ErrInvalidSampleLength
			}
			for i := uint32(0); i < plength>>1; i++ {
				dr.ReadUint32()
				dr.ReadUint32()
			}
		case propUnit:
			u := dr.ReadUint32()
			if u < 0 || u > 4 {
				return nil, ErrInvalidUnit
			}
		case propUserUnit:
			dr.ReadFloat32() // factor
			dr.ReadUint32()  // number of decimal igits
			dr.ReadString()  // id
			dr.ReadString()  // symbol
			dr.ReadString()  // abbr.
			dr.ReadString()  // singular name
			dr.ReadString()  // plural name
		case propVectors:
			dr.ReadVectors()
		default:
			dr.Skip(plength)
		}
	}

	layerptrs := make([]uint32, 0, 32)
	for {
		lptr := dr.ReadUint32()
		if lptr == 0 {
			break
		}
		layerptrs = append(layerptrs, lptr)
	}

	if dr.Err != nil {
		return nil, dr.Err
	}

	type groupOffset struct {
		Group            limage.Image
		OffsetX, OffsetY int
		Parent           *limage.Image
		Offset           int
	}

	var (
		groups = make(map[string]*groupOffset)
		n      rune
		alpha  = true
	)

	layers := make([]layer, len(layerptrs))

	var (
		errLock sync.Mutex
		wg      sync.WaitGroup
	)
	wg.Add(len(layerptrs))
	for n, lptr := range layerptrs {
		go func(n int, lptr uint32) {
			d := decoder{
				reader:      newReader(r),
				Width:       width,
				Height:      height,
				baseType:    baseType,
				palette:     palette,
				compression: compression,
			}
			d.Goto(lptr)
			layers[n] = d.ReadLayer()
			if d.Err != nil {
				errLock.Lock()
				dr.SetError(d.Err)
				errLock.Unlock()
			}
			wg.Done()
		}(n, lptr)
	}

	wg.Wait()

	if dr.Err != nil {
		return nil, dr.Err
	}

	groups[""] = &groupOffset{Group: make(limage.Image, 0, 32)}
	for _, l := range layers {
		if !alpha {
			return nil, ErrMissingAlpha
		}
		alpha = l.alpha
		if len(l.itemPath) == 0 {
			l.itemPath = []rune{n}
			n++
		}
		g := groups[string(l.itemPath[:len(l.itemPath)-1])]
		if g == nil {
			return nil, ErrInvalidGroup
		}
		if l.group {
			groups[string(l.itemPath)] = &groupOffset{
				Group:   make(limage.Image, 0, 32),
				OffsetX: l.OffsetX,
				OffsetY: l.OffsetY,
				Parent:  &g.Group,
				Offset:  len(g.Group),
			}
		}
		l.OffsetX -= g.OffsetX
		l.OffsetY -= g.OffsetY
		g.Group = append(g.Group, l.Layer)
	}

	var im limage.Image

	for _, g := range groups {
		ng := make(limage.Image, len(g.Group))
		copy(ng, g.Group)
		g.Group = ng
		if g.Parent == nil {
			im = ng
		} else {
			(*g.Parent)[g.Offset].Image = ng
		}
	}

	if len(im) > 0 {
		switch im[len(im)-1].Mode {
		case limage.CompositeNormal, limage.CompositeDissolve:
		default:
			im[len(im)-1].Mode = 0
		}
	}

	return im, nil
}

// Errors
const (
	ErrInvalidFileTypeID   errors.Error = "invalid file type identification"
	ErrUnsupportedVersion  errors.Error = "unsupported file version"
	ErrInvalidHeader       errors.Error = "invalid header"
	ErrInvalidProperties   errors.Error = "invalid property list"
	ErrInvalidOpacity      errors.Error = "opacity not in valid range"
	ErrInvalidGuideLength  errors.Error = "invalid guide length"
	ErrInvalidUnit         errors.Error = "invalid unit"
	ErrInvalidSampleLength errors.Error = "invalid sample points length"
	ErrInvalidGroup        errors.Error = "invalid or unknown group specified for layer"
	ErrUnknownCompression  errors.Error = "unknown compression method"
	ErrMissingAlpha        errors.Error = "non-bottom layer missing alpha channel"
)
