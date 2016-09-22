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

func Encode(w io.WriterAt, i image.Image) error {

	var im *limage.Image

	if m, ok := i.(*limage.Image); ok {
		im = m
	} else {
		b := i.Bounds()
		im = &limage.Image{
			Group: limage.Group{
				Config: image.Config{
					ColorModel: i.ColorModel(),
					Width:      b.Dx(),
					Height:     b.Dy(),
				},
				Layers: []limage.Layer{{Image: i}},
			},
		}
	}

	e := encoder{
		writer:   newWriter(w),
		channels: make(map[string]uint32),
	}

	e.WriteHeader(im.Config)

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

	e.WriteUint32(0)

	e.WriteGroup(&im.Group, make([]int32, 0, 32), e.ReserveSpace(layerCount(&im.Group)<<2))

	e.WriteUint32(0)

	return e.Err
}

var header = []byte{'g', 'i', 'm', 'p', ' ', 'x', 'c', 'f', 'v', '0', '0', '3', 0}

func (e *encoder) WriteHeader(c image.Config) {
	e.Write(header)
	e.WriteUint32(uint32(c.Width))
	e.WriteUint32(uint32(c.Height))
	switch c.ColorModel {
	case color.GrayModel, color.Gray16Model, lcolor.GrayAlphaModel:
		e.colorModel = lcolor.GrayAlphaModel
		e.colorType = 1
	default:
		switch m := c.ColorModel.(type) {
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

func (e *encoder) WriteGroup(g *limage.Group, groups []int32, w writer) uint32 {
	ptr := e.WriterLayer(g)
	for n, layer := range g.Layers {
		nGroups := append(groups, n)
		switch l := layer.Image.(type) {
		case limage.Group:
			w.WriteInt32(e.WriteGroup(&l, nGroups, w))
		case *limage.Group:
			w.WriteInt32(e.WriteGroup(l, nGroups, w))
		default:
			w.WriteInt32(e.WriterLayer(layer.Image, groups))
		}
	}
	return ptr
}

func (e *encoder) WriterLayer(l image.Image, groups []int32) uint32 {

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
