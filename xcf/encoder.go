package xcf

import (
	"image"
	"image/color"
	"io"

	"github.com/MJKWoolnough/limage"
	"github.com/MJKWoolnough/limage/lcolor"
)

const chanLen = 64 * 64 * 1 // tile width (64) * tile height (64) * max channels (4) * bitwidth (1)

type colourBufFunc func(*encoder, color.Color)

type encoder struct {
	writer

	colourPalette  lcolor.AlphaPalette
	colourType     uint8
	colourFunc     colourBufFunc
	colourChannels uint8

	channelBuf [chanLen*4 + 4]byte // 4 channels max + 4 for max colourBuf
}

func Encode(w io.WriterAt, im image.Image) error {

	if li, ok := im.(limage.Image); ok {
		im = &li
	}

	e := encoder{
		writer: newWriter(w),
	}

	e.WriteHeader(im)

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

	switch im := im.(type) {
	case *limage.Image:
		if im.Comment != "" {
			// write comment parasite
		}
		e.WriteUint32(0)
		e.WriteLayers(im.Layers, make([]int32, 0, 32), e.ReserveSpace(layerCount(&im.Group)<<2))
	default:
		e.WriteUint32(0)
		e.WriteLayer(limage.Layer{Image: im}, []int32{}, e.ReserveSpace(4))
	}

	e.WriteUint32(0)

	return e.Err
}

var header = []byte{'g', 'i', 'm', 'p', ' ', 'x', 'c', 'f', 'v', '0', '0', '3', 0}

func (e *encoder) WriteHeader(im image.Image) {
	e.Write(header)
	b := im.Bounds()
	e.WriteUint32(uint32(b.Dx()))
	e.WriteUint32(uint32(b.Dy()))
	var colourType uint32
	switch cm := im.ColorModel(); cm {
	case color.GrayModel, color.Gray16Model, lcolor.GrayAlphaModel:
		colourType = 1
		e.colourFunc = (*encoder).grayAlphaToBuf
		e.colourChannels = 2
	default:
		switch m := cm.(type) {
		case color.Palette:
			e.colourPalette = lcolor.AlphaPalette(m)
			colourType = 2
			e.colourFunc = (*encoder).paletteAlphaToBuf
			e.colourChannels = 2
		case lcolor.AlphaPalette:
			e.colourPalette = m
			colourType = 2
			e.colourFunc = (*encoder).paletteAlphaToBuf
			e.colourChannels = 2
		default:
			colourType = 0
			e.colourFunc = (*encoder).rgbAlphaToBuf
			e.colourChannels = 4
		}
	}
	e.WriteUint32(colourType)
}

func layerCount(g *limage.Group) int64 {
	count := int64(len(g.Layers))
	for _, l := range g.Layers {
		switch g := l.Image.(type) {
		case limage.Group:
			count += layerCount(&g)
		case *limage.Group:
			count += layerCount(g)
		}

	}
	return count
}

func (e *encoder) WriteLayers(layers []limage.Layer, groups []int32, w writer) {
	for n, layer := range layers {
		nGroups := append(groups, int32(n))
		w.WriteUint32(e.WriteLayer(layer, nGroups, w))
	}
}

func (e *encoder) WriteLayer(im limage.Layer, groups []int32, w writer) uint32 {
	var ptr uint32

	// write layer

	var g *limage.Group
	switch i := im.Image.(type) {
	case limage.Group:
		g = &i
	case *limage.Group:
		g = i
	default:
		return ptr
	}
	e.WriteLayers(g.Layers, groups, w)
	return ptr
}

func (e *encoder) WriteImage(im image.Image) {
	w := e.ReserveSpace(8) // 2 uint32 pointers

	var mask *image.Gray

	switch imm := im.(type) {
	case limage.MaskedImage:
		im = imm.Image
		mask = imm.Mask
	case *limage.MaskedImage:
		im = imm.Image
		mask = imm.Mask
	}

	w.WriteUint32(e.Count)

	// image hierarchy

	bounds := im.Bounds()

	e.WriteUint32(uint32(bounds.Dx()))
	e.WriteUint32(uint32(bounds.Dy()))
	e.WriteUint32(uint32(e.colourChannels))

	e.WriteTiles(im, e.colourFunc, e.colourChannels)

	if mask != nil {
		w.WriteUint32(e.Count)
		e.WriteTiles(mask, grayToBuf, 1)
	}

}

func (e *encoder) WriteTiles(im image.Image, colourFunc colourBufFunc, colourChannels uint8) {
	bounds := im.Bounds()

	dx := int64(bounds.Dx())
	dy := int64(bounds.Dy())

	nx := dx >> 6 // each tile is 64 wide
	ny := dy >> 6 // each tile is 64 high

	if dx&63 > 0 { // last tile not as wide
		nx++
	}
	if dy&63 > 0 { // last tile not as high
		ny++
	}

	w := e.ReserveSpace((nx * ny) << 2)
	channels := make([][]byte, colourChannels)
	r := rlencoder{Writer: e.StickyWriter}
	for i := 0; i < int(colourChannels); i++ {
		channels[i] = e.channelBuf[i*chanLen : i*chanLen : (i+1)*chanLen]
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 64 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 64 {
			for n := range channels {
				channels[n] = channels[n][:0]
			}
			for j := y; j < y+64 && j < bounds.Max.Y; j++ {
				for i := x; i < x+64 && i < bounds.Max.X; i++ {
					colourFunc(e, im.At(i, j))
					for n := range channels {
						channels[n] = append(channels[n], e.channelBuf[4*chanLen+n])
					}
				}
			}
			ptr := uint32(e.Count)
			for _, channel := range channels {
				r.Write(channel)
				r.Flush()
			}
			w.WriteUint32(ptr)
		}
	}
}

func (e *encoder) rgbAlphaToBuf(c color.Color) {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	e.channelBuf[4*chanLen] = rgba.R
	e.channelBuf[4*chanLen+1] = rgba.G
	e.channelBuf[4*chanLen+2] = rgba.B
	e.channelBuf[4*chanLen+3] = rgba.A
}

func (e *encoder) grayAlphaToBuf(c color.Color) {
	ga := lcolor.GrayAlphaModel.Convert(c).(lcolor.GrayAlpha)
	e.channelBuf[4*chanLen] = ga.Y
	e.channelBuf[4*chanLen+1] = ga.A
}

func (e *encoder) grayToBuf(c color.Color) {
	e.channelBuf[4*chanLen] = color.GrayModel.Convert(c).(color.Gray).Y
}

func (e *encoder) paletteAlphaToBuf(c color.Color) {
	r, g, b, a := c.RGBA()
	i := e.colourPalette.Index(lcolor.RGB{uint8(r), uint8(g), uint8(b)})
	e.channelBuf[4*chanLen] = uint8(i)
	e.channelBuf[4*chanLen+1] = uint8(a)
}
