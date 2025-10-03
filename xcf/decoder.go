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

	mode, err := readHeader(dr)
	if err != nil {
		return image.Config{}, err
	}

	c.Width = int(dr.ReadUint32())
	c.Height = int(dr.ReadUint32())
	baseType := dr.ReadUint32()

	if mode == 2 {
		dr.ReadUint32()
	}

	switch baseType {
	case 0:
		c.ColorModel = color.NRGBAModel
	case 1:
		c.ColorModel = lcolor.GrayAlphaModel
	case 2:
		palette, _, err := readImageProperties(dr, 2)
		if err != nil {
			return c, err
		}

		c.ColorModel = palette
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

type groupOffset struct {
	Group            limage.Image
	OffsetX, OffsetY int
	Parent           *limage.Image
	Offset           int
}

func decodeImage(r io.ReaderAt, decompress bool) (limage.Image, error) {
	dr := newReader(r)

	mode, err := readHeader(dr)
	if err != nil {
		return nil, err
	}

	width := int(dr.ReadUint32())
	height := int(dr.ReadUint32())
	bounds := image.Rect(0, 0, width, height)
	baseType := dr.ReadUint32()

	var precision uint32

	if mode > 0 {
		precision = dr.ReadUint32()
	}

	palette, compression, err := readImageProperties(dr, baseType)
	if err != nil {
		return nil, err
	}

	layerptrs := readLayerPointers(dr, mode)

	if dr.Err != nil {
		return nil, dr.Err
	}

	layers := readLayers(dr, r, layerptrs, baseType, palette, compression, precision, mode, decompress)

	if dr.Err != nil {
		return nil, dr.Err
	}

	groups, err := makeGroups(layers, bounds)
	if err != nil {
		return nil, err
	}

	im := makeImage(groups)

	return im, nil
}

func readHeader(dr reader) (uint32, error) {
	var header [14]byte

	dr.Read(header[:])

	if dr.Err != nil {
		return 0, dr.Err // wrap?
	}

	if string(header[:9]) != fileTypeID {
		return 0, ErrInvalidFileTypeID
	}

	var mode uint32

	switch string(header[9:13]) {
	case fileVersion0, fileVersion1, fileVersion2, fileVersion3:
	case fileVersion4, fileVersion5, fileVersion6, fileVersion7, fileVersion8, fileVersion9, fileVersion10:
		mode = 1
	case fileVersion11, fileVersion12, fileVersion13:
		mode = 2
	default:
		return 0, ErrUnsupportedVersion
	}

	if header[13] != 0 {
		return 0, ErrInvalidHeader
	}

	return mode, nil
}

func readImageProperties(dr reader, baseType uint32) (lcolor.AlphaPalette, uint8, error) {
	var (
		palette     lcolor.AlphaPalette
		compression uint8
	)

PropertyLoop:
	for {
		typ := dr.ReadUint32()
		plength := dr.ReadUint32()

		switch typ {
		case propEnd:
			if plength != 0 {
				return nil, 0, ErrInvalidProperties
			}

			break PropertyLoop

		// general properties
		case propLinked:
			dr.ReadBoolProperty()
		case propLockContent:
			dr.ReadBoolProperty()
		case propOpacity:
			if o := dr.ReadUint32(); o > 255 {
				return nil, 0, ErrInvalidOpacity
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
				return nil, 0, ErrUnknownCompression
			}
		case propGuides:
			ng := plength / 5

			if ng*5 != plength {
				return nil, 0, ErrInvalidGuideLength
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
				return nil, 0, ErrInvalidSampleLength
			}

			for i := uint32(0); i < plength>>1; i++ {
				dr.SkipUint32()
				dr.SkipUint32()
			}
		case propUnit:
			if dr.ReadUint32() > 4 {
				return nil, 0, ErrInvalidUnit
			}
		case propUserUnit:
			dr.SkipFloat32() // factor
			dr.SkipUint32()  // number of decimal digits
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

	return palette, compression, nil
}

func readLayerPointers(dr reader, mode uint32) []uint64 {
	layerptrs := make([]uint64, 0, 32)

	for {
		var lptr uint64

		if mode < 2 {
			lptr = uint64(dr.ReadUint32())
		} else {
			lptr = dr.ReadUint64()
		}

		if lptr == 0 {
			return layerptrs
		}

		layerptrs = append(layerptrs, lptr)
	}
}

func readLayers(dr reader, r io.ReaderAt, layerptrs []uint64, baseType uint32, palette lcolor.AlphaPalette, compression uint8, precision, mode uint32, decompress bool) []layer {
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

	return layers
}

func makeGroups(layers []layer, bounds image.Rectangle) (map[string]*groupOffset, error) {
	var (
		groups = make(map[string]*groupOffset)
		n      rune
		alpha  = true
	)

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

	return groups, nil
}

func makeImage(groups map[string]*groupOffset) limage.Image {
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

	return im
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
