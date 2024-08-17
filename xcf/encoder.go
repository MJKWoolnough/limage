package xcf

import (
	"image"
	"image/color"
	"io"

	"vimagination.zapto.org/limage"
	"vimagination.zapto.org/limage/lcolor"
)

const chanLen = 64 * 64 * 1 // tile width (64) * tile height (64) * max channels (4) * bytewidth (1).

type colourBufFunc func(*encoder, color.Color)

type encoder struct {
	writer

	colourPalette  lcolor.AlphaPalette
	colourFunc     colourBufFunc
	colourType     uint8
	colourChannels uint8

	channelBuf [4][chanLen]byte
	colourBuf  [4]byte
}

// Encode encodes the given image as an XCF file to the given WriterAt.
func Encode(w io.WriterAt, im image.Image) error {
	switch imt := im.(type) {
	case *limage.Image:
		im = *imt
	case limage.Layer:
		im = limage.Image{imt}
	case *limage.Layer:
		im = limage.Image{*imt}
	}

	e := encoder{
		writer: newWriter(w),
	}

	e.Write(header)

	b := im.Bounds()

	e.WriteUint32(uint32(b.Dx()))
	e.WriteUint32(uint32(b.Dy()))

	switch cm := im.ColorModel(); cm {
	case color.GrayModel, color.Gray16Model, lcolor.GrayAlphaModel:
		e.colourType = 1
		e.colourFunc = (*encoder).grayAlphaToBuf
		e.colourChannels = 2
	default:
		switch m := cm.(type) {
		case color.Palette:
			e.colourPalette = lcolor.AlphaPalette(m)
			e.colourType = 2
			e.colourFunc = (*encoder).paletteAlphaToBuf
			e.colourChannels = 2
		case lcolor.AlphaPalette:
			e.colourPalette = m
			e.colourType = 2
			e.colourFunc = (*encoder).paletteAlphaToBuf
			e.colourChannels = 2
		default:
			e.colourType = 0
			e.colourFunc = (*encoder).rgbAlphaToBuf
			e.colourChannels = 4
		}
	}
	e.WriteUint32(uint32(e.colourType))

	// write property list

	if e.colourPalette != nil {
		e.WriteUint32(propColorMap)
		e.WriteUint32(3*uint32(len(e.colourPalette)) + 4)
		e.WriteUint32(uint32(len(e.colourPalette)))

		for _, colour := range e.colourPalette {
			rgb := lcolor.RGBModel.Convert(colour).(lcolor.RGB)

			e.WriteUint8(rgb.R)
			e.WriteUint8(rgb.G)
			e.WriteUint8(rgb.B)
		}
	}

	e.WriteUint32(propCompression)
	e.WriteUint32(1)
	e.WriteUint8(1) // rle

	e.WriteUint32(0)
	e.WriteUint32(0)

	switch im := im.(type) {
	case limage.Image:
		pw := e.ReservePointerList(layerCount(im))

		e.WriteUint32(0) // no channels
		e.WriteLayers(im, 0, 0, make([]uint32, 0, 32), pw)
	default:
		pw := e.ReservePointerList(1)

		e.WriteUint32(0) // no channels
		e.WriteLayer(limage.Layer{Image: im}, 0, 0, []uint32{}, pw)
	}

	return e.Err
}

func layerCount(g limage.Image) uint32 {
	count := uint32(len(g))

	for _, l := range g {
		switch g := l.Image.(type) {
		case limage.Image:
			count += layerCount(g)
		case *limage.Image:
			count += layerCount(*g)
		}
	}

	return count
}

var header = []byte{'g', 'i', 'm', 'p', ' ', 'x', 'c', 'f', ' ', 'v', '0', '0', '3', 0}
