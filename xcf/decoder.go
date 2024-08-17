// Package xcf implements an image encoder and decoder for GIMPs XCF format
package xcf // import "vimagination.zapto.org/limage/xcf"

import (
	"errors"
	"image"
	"image/color"
	"io"
	"sync"

	"vimagination.zapto.org/limage"
	"vimagination.zapto.org/limage/internal"
	"vimagination.zapto.org/limage/lcolor"
)

func decodeConfig(r io.Reader) (image.Config, error) {
	ra := internal.GetReaderAt(r)

	if ra == nil {
		return image.Config{}, ErrNeedReaderAt
	}

	return DecodeConfig(ra)
}

func decode(r io.Reader) (image.Image, error) {
	ra := internal.GetReaderAt(r)

	if ra == nil {
		return nil, ErrNeedReaderAt
	}

	return Decode(ra)
}

func init() {
	image.RegisterFormat("xcf", fileTypeID, decode, decodeConfig)
}

const (
	fileTypeID    = "gimp xcf "
	fileVersion0  = "file"
	fileVersion1  = "v001"
	fileVersion2  = "v002"
	fileVersion3  = "v003"
	fileVersion4  = "v004"
	fileVersion5  = "v005"
	fileVersion6  = "v006"
	fileVersion7  = "v007"
	fileVersion8  = "v008"
	fileVersion9  = "v009"
	fileVersion10 = "v010"
	fileVersion11 = "v011"
	fileVersion12 = "v012"
	fileVersion13 = "v013"
)

const (
	// baseRGB     = 0
	// baseGrey    = 1
	baseIndexed = 2
)

type decoder struct {
	reader
	compression uint8
	decompress  bool
	baseType    uint32
	palette     lcolor.AlphaPalette
	precision   uint32
	mode        uint32
}

// DecodeConfig retrieves the color model and dimensions of the XCF image.
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

	var newMode bool

	switch string(header[9:13]) {
	case fileVersion0, fileVersion1, fileVersion2, fileVersion3:
	case fileVersion4, fileVersion5, fileVersion6, fileVersion7, fileVersion8, fileVersion9, fileVersion10, fileVersion11, fileVersion12, fileVersion13:
		newMode = true
	default:
		return c, ErrUnsupportedVersion
	}

	if header[13] != 0 {
		return c, ErrInvalidHeader
	}

	c.Width = int(dr.ReadUint32())
	c.Height = int(dr.ReadUint32())
	baseType := dr.ReadUint32()

	if newMode {
		dr.ReadUint32()
	}

	switch baseType {
	case 0:
		c.ColorModel = color.NRGBAModel
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

			// general properties
			case propLinked:
				dr.SkipBoolProperty()
			case propLockContent:
				dr.SkipBoolProperty()
			case propOpacity:
				if o := dr.ReadUint32(); o > 255 {
					return c, ErrInvalidOpacity
				}
			case propParasites:
				dr.SkipParasites(plength)
			case propTattoo:
				dr.SkipUint32()
			case propVisible:
				dr.SkipBoolProperty()
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
					dr.SkipUint32()
					dr.SkipBoolProperty()
				}
			case propPaths:
				dr.SkipPaths()
			case propResolution:
				dr.SkipFloat32()
				dr.SkipFloat32()
			case propSamplePoints:
				if plength&1 == 1 {
					return c, ErrInvalidSampleLength
				}

				for i := uint32(0); i < plength>>1; i++ {
					dr.SkipUint32()
					dr.SkipUint32()
				}
			case propUnit:
				if unit := dr.ReadUint32(); unit > 4 {
					return c, ErrInvalidUnit
				}
			case propUserUnit:
				dr.SkipFloat32()
				dr.SkipUint32()
				dr.SkipString()
				dr.SkipString()
				dr.SkipString()
				dr.SkipString()
				dr.SkipString()
			case propVectors:
				dr.SkipVectors()
			default:
				dr.Skip(plength)
			}
		}
	}

	return c, dr.Err
}

// Decode reads an XCF layered image from the given ReaderAt.
func Decode(r io.ReaderAt) (limage.Image, error) {
	return decodeImage(r, true)
}

// DecodeCompressed reads an XCF layered image, as Decode, but defers decoding
// and decompressing, doing so upon an At method.
func DecodeCompressed(r io.ReaderAt) (limage.Image, error) {
	return decodeImage(r, false)
}

func decodeImage(r io.ReaderAt, decompress bool) (limage.Image, error) {
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

	var mode uint32

	switch string(header[9:13]) {
	case fileVersion0, fileVersion1, fileVersion2, fileVersion3:
	case fileVersion4, fileVersion5, fileVersion6, fileVersion7, fileVersion8, fileVersion9, fileVersion10:
		mode = 1
	case fileVersion11, fileVersion12, fileVersion13:
		mode = 2
	default:
		return nil, ErrUnsupportedVersion
	}

	if header[13] != 0 {
		return nil, ErrInvalidHeader
	}

	width := int(dr.ReadUint32())
	height := int(dr.ReadUint32())
	bounds := image.Rect(0, 0, width, height)
	baseType := dr.ReadUint32()

	var precision uint32

	if mode > 0 {
		precision = dr.ReadUint32()
	}

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

		// general properties
		case propLinked:
			dr.ReadBoolProperty()
		case propLockContent:
			dr.ReadBoolProperty()
		case propOpacity:
			if o := dr.ReadUint32(); o > 255 {
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
				palette[n] = lcolor.RGB{
					R: dr.ReadUint8(),
					G: dr.ReadUint8(),
					B: dr.ReadUint8(),
				}
			}
		case propCompression:
			if compression = dr.ReadUint8(); compression > 1 {
				return nil, ErrUnknownCompression
			}
		case propGuides:
			ng := plength / 5

			if ng*5 != plength {
				return nil, ErrInvalidGuideLength
			}

			for n := uint32(0); n < ng; n++ {
				dr.SkipUint32()
				dr.SkipBoolProperty()
			}
		case propPaths:
			dr.SkipPaths()
		case propResolution:
			dr.SkipFloat32() // x
			dr.SkipFloat32() // y
		case propSamplePoints:
			if plength&1 == 1 {
				return nil, ErrInvalidSampleLength
			}

			for i := uint32(0); i < plength>>1; i++ {
				dr.SkipUint32()
				dr.SkipUint32()
			}
		case propUnit:
			if dr.ReadUint32() > 4 {
				return nil, ErrInvalidUnit
			}
		case propUserUnit:
			dr.SkipFloat32() // factor
			dr.SkipUint32()  // number of decimal igits
			dr.SkipString()  // id
			dr.SkipString()  // symbol
			dr.SkipString()  // abbr.
			dr.SkipString()  // singular name
			dr.SkipString()  // plural name
		case propVectors:
			dr.SkipVectors()
		default:
			dr.Skip(plength)
		}
	}

	layerptrs := make([]uint64, 0, 32)

	for {
		var lptr uint64

		if mode < 2 {
			lptr = uint64(dr.ReadUint32())
		} else {
			lptr = dr.ReadUint64()
		}

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
		go func(n int, lptr uint64) {
			d := decoder{
				reader:      newReader(r),
				baseType:    baseType,
				palette:     palette,
				compression: compression,
				decompress:  decompress,
				precision:   precision,
				mode:        mode,
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
				OffsetX: l.LayerBounds.Min.X,
				OffsetY: l.LayerBounds.Min.Y,
				Parent:  &g.Group,
				Offset:  len(g.Group),
			}
		}

		l.LayerBounds = l.LayerBounds.Intersect(bounds).Sub(image.Pt(g.OffsetX, g.OffsetY))
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

// Errors.
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
	ErrMissingAlpha        = errors.New("non-bottom layer missing alpha channel")
	ErrNeedReaderAt        = errors.New("need a io.ReaderAt")
)
