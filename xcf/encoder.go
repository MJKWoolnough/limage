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
				Layers: []limage.Layer{
					{
						Image:   i,
						Visible: true,
						Opacity: 255,
					},
				},
			},
			Opacity: 255,
		}
	}

	e := encoder{
		writer:   newWriter(w),
		channels: make(map[string]uint64),
	}

	e.WriteHeader(im.Config)

	// write property list

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
		e.WriteUint32(1)
	default:
		switch c.ColorModel.(type) {
		case color.Palette, lcolor.AlphaPalette:
			e.colorModel = c.ColorModel
			e.WriteUint32(2)
		default:
			e.colorModel = color.RGBAModel
			e.WriteUint32(0)
		}
	}
}
