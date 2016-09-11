package xcf

import (
	"errors"
	"image"
	"io"

	"github.com/MJKWoolnough/limage"
)

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
	_ = im
	return errors.New("unimplemented")
}
