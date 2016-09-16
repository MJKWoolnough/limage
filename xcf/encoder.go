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
	channels map[string]uint64

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
		channels: make(map[string]uint64),
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

	// write channel list

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
