package xcf

import (
	"image"
	"image/color"
	"io"

	"github.com/MJKWoolnough/limage"
	"github.com/MJKWoolnough/limage/lcolor"
)

const chanLen = 64 * 64 * 1 // tile width (64) * tile height (64) * max channels (4) * bitwidth (1)

type encoder struct {
	writer

	colorPalette lcolor.AlphaPalette
	colorType    uint8

	channelBuf [chanLen * 4]byte // 4 channels max
}

func Encode(w io.WriterAt, im image.Image) error {

	if li, ok := im.(limage.Image); ok {
		im = &li
	}

	e := encoder{
		writer:   newWriter(w),
		channels: make(map[string]uint32),
	}

	e.WriteHeader(im)

	// write property list

	if e.colorPalette != nil {
		e.WriteUint32(propColorMap)
		e.WriteUint32(3*uint32(len(e.colorPalette)) + 4)
		e.WriteUint32(uint32(len(e.colorPalette)))
		for _, colour := range e.colorPalette {
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
	switch cm := im.ColorModel(); cm {
	case color.GrayModel, color.Gray16Model, lcolor.GrayAlphaModel:
		e.colorType = 1
	default:
		switch m := cm.(type) {
		case color.Palette:
			e.colorPalette = lcolor.AlphaPalette(m)
			e.colorType = 2
		case lcolor.AlphaPalette:
			e.colorPalette = m
			e.colorType = 2
		default:
			e.colorType = 0
		}
	}
	e.WriteUint32(uint32(e.colorType))
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

func (e *encoder) WriteRGBTile(im image.Image, bounds image.Rectangle, w writer) {
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 64 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 64 {
			red := e.channelBuf[0:0:chanLen]
			green := e.channelBuf[chanLen:chanLen : chanLen*2]
			blue := e.channelBuf[chanLen*2 : chanLen*2 : chanLen*3]
			alpha := e.channelBuf[chanLen*3 : chanLen*3 : chanLen*4]
			for j := y; j < y+64 && j < bounds.Max.Y; j++ {
				for i := x; i < x+64 && i < bounds.Max.X; i++ {
					r, g, b, a := im.At(i, j).RGBA()
					red = append(red, uint8(r))
					green = append(green, uint8(g))
					blue = append(blue, uint8(b))
					alpha = append(alpha, uint8(a))
				}
			}
			w.WriteUint32(e.WriteChannel(red))
			w.WriteUint32(e.WriteChannel(green))
			w.WriteUint32(e.WriteChannel(blue))
			w.WriteUint32(e.WriteChannel(alpha))
		}
	}
}

func (e *encoder) WriteGrayTile(im image.Image, bounds image.Rectangle, w writer) {
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 64 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 64 {
			gamma := e.channelBuf[0:0:chanLen]
			alpha := e.channelBuf[chanLen*3 : chanLen*3 : chanLen*4]
			for j := y; j < y+64 && j < bounds.Max.Y; j++ {
				for i := x; i < x+64 && i < bounds.Max.X; i++ {
					g, _, _, a := im.At(i, j).RGBA()
					gamma = append(alpha, uint8(g))
					alpha = append(alpha, uint8(a))
				}
			}
			w.WriteUint32(e.WriteChannel(gamma))
			w.WriteUint32(e.WriteChannel(alpha))
		}
	}
}

func (e *encoder) WritePalettedTile(im image.Image, bounds image.Rectangle, w writer) {
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 64 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 64 {
			index := e.channelBuf[0:0:chanLen]
			alpha := e.channelBuf[chanLen*3 : chanLen*3 : chanLen*4]
			for j := y; j < y+64 && j < bounds.Max.Y; j++ {
				for i := x; i < x+64 && i < bounds.Max.X; i++ {
					c := im.At(i, j)
					in := e.colorPalette.Index(c)
					_, _, _, a := c.RGBA()
					index = append(index, uint8(in))
					alpha = append(alpha, uint8(a))
				}
			}
			w.WriteUint32(e.WriteChannel(index))
			w.WriteUint32(e.WriteChannel(alpha))
		}
	}
}

func (e *encoder) WriteChannels(data ...[]byte) uint32 {
	ptr := uint32(e.Count)
	r := rlencoder{Writer: e.StickyWriter}
	for _, d := range data {
		r.Write(d)
		r.Flush()
	}
	return ptr
}
