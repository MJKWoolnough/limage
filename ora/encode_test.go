package ora

import (
	"archive/zip"
	"image"
	"image/color"
	"testing"

	"github.com/MJKWoolnough/limage"
	"github.com/MJKWoolnough/memio"
)

type colourSquare struct {
	image.Rectangle
	color.Color
}

func (c colourSquare) ColorModel() color.Model {
	return color.ModelFunc(func(color.Color) color.Color {
		return c.Color
	})
}

func (c colourSquare) Bounds() image.Rectangle {
	return c.Rectangle
}

func (c colourSquare) At(int, int) color.Color {
	return c.Color
}

func TestEncode(t *testing.T) {
	tests := []limage.Image{
		{
			limage.Layer{
				Name: "Layer Group",
				Image: limage.Image{
					limage.Layer{
						Name: "Blue",
						Image: singleColourImage{
							Colour: color.RGBA64{B: 65535, A: 65535},
							Width:  30,
							Height: 30,
						},
						LayerBounds: image.Rect(0, 0, 30, 30),
					},
					limage.Layer{
						Name: "Red",
						Image: singleColourImage{Colour: color.RGBA64{R: 65535, A: 65535},
							Width:  30,
							Height: 30,
						},
						LayerBounds: image.Rect(20, 20, 50, 50),
					},
				},
				LayerBounds: image.Rect(0, 0, 50, 50),
			},
			limage.Layer{
				Name: "Background",
				Image: singleColourImage{
					Colour: color.RGBA64{A: 65535},
					Width:  50,
					Height: 50,
				},
				LayerBounds: image.Rect(0, 0, 50, 50),
			},
		},
		{
			limage.Layer{
				Name: "Layer",
				Image: singleColourImage{
					Colour: color.RGBA64{R: 65535, A: 65535},
					Width:  30,
					Height: 30,
				},
				LayerBounds: image.Rect(10, 10, 40, 40),
			},
			limage.Layer{
				Name: "Background",
				Image: singleColourImage{
					Colour: color.RGBA64{A: 65535},
					Width:  50,
					Height: 50,
				},
				LayerBounds: image.Rect(0, 0, 50, 50),
			},
		},
		{
			limage.Layer{
				Name: "Background",
				Image: singleColourImage{
					Colour: color.RGBA64{A: 65535},
					Width:  50,
					Height: 50,
				},
				LayerBounds: image.Rect(0, 0, 50, 50),
			},
		},
		{
			limage.Layer{
				Name: "Background",
				Image: singleColourImage{
					Colour: color.RGBA64{R: 65535, A: 65535},
					Width:  50,
					Height: 50,
				},
				LayerBounds: image.Rect(0, 0, 50, 50),
			},
		},
		{
			limage.Layer{
				Name: "Background",
				Image: singleColourImage{
					Colour: color.RGBA64{R: 65535, G: 65535, B: 65535, A: 65535},
					Width:  50,
					Height: 50,
				},
				LayerBounds: image.Rect(0, 0, 50, 50),
			},
		},
	}

	buf := make([]byte, 683)

	for n, test := range tests {
		buf = buf[:0]
		if err := Encode(memio.Create(&buf), test); err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
			continue
		}
		f, err := zip.NewReader(memio.Open(buf), int64(len(buf)))
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
			continue
		}
		l, err := Decode(f)
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
			continue
		}
		if err := compareLayers(l, test); err != nil {
			t.Errorf("test %d: %s", n+1, err)
		}
	}
}
