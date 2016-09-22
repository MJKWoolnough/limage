package xcf

import (
	"image"
	"image/color"
	"io"

	"github.com/MJKWoolnough/limage"
	"github.com/MJKWoolnough/limage/lcolor"
)

type encoder struct {
	writer
	channels map[string]uint32

	colorModel color.Model
	colorType  uint8
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

	if p, ok := e.colorModel.(lcolor.AlphaPalette); ok {
		e.WriteUint32(propColorMap)
		e.WriteUint32(3*uint32(len(p)) + 4)
		e.WriteUint32(uint32(len(p)))
		for _, colour := range p {
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
		e.colorModel = lcolor.GrayAlphaModel
		e.colorType = 1
	default:
		switch m := cm.(type) {
		case color.Palette:
			e.colorModel = lcolor.AlphaPalette(m)
			e.colorType = 2
		case lcolor.AlphaPalette:
			e.colorModel = m
			e.colorType = 2
		default:
			e.colorModel = color.RGBAModel
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

func (e *encoder) WriteChannel(d []byte) uint32 {
	if ptr, ok := e.channels[string(d)]; ok {
		return ptr
	}
	ptr := uint32(e.Count)
	e.channels[string(d)] = ptr
	r := rlencoder{Writer: e.StickyWriter}
	r.Write(d)
	r.Flush()
	return ptr
}
