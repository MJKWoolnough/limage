package xcf

import (
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
							Colour: color.NRGBA{B: 255, A: 255},
							Width:  30,
							Height: 30,
						},
					},
					limage.Layer{
						Name: "Red",
						Image: singleColourImage{Colour: color.NRGBA{R: 255, A: 255},
							Width:  30,
							Height: 30,
						},
						OffsetX: 20,
						OffsetY: 20,
					},
				},
			},
			limage.Layer{
				Name: "Background",
				Image: singleColourImage{
					Colour: color.NRGBA{A: 255},
					Width:  50,
					Height: 50,
				},
			},
		},
		{
			limage.Layer{
				Name: "Layer",
				Image: singleColourImage{
					Colour: color.NRGBA{R: 255, A: 255},
					Width:  30,
					Height: 30,
				},
				OffsetX: 10,
				OffsetY: 10,
			},
			limage.Layer{
				Name: "Background",
				Image: singleColourImage{
					Colour: color.NRGBA{A: 255},
					Width:  50,
					Height: 50,
				},
			},
		},
		{
			limage.Layer{
				Name: "Background",
				Image: singleColourImage{
					Colour: color.NRGBA{A: 255},
					Width:  50,
					Height: 50,
				},
			},
		},
		{
			limage.Layer{
				Name: "Background",
				Image: singleColourImage{
					Colour: color.NRGBA{R: 255, A: 255},
					Width:  50,
					Height: 50,
				},
			},
		},
		{
			limage.Layer{
				Name: "Background",
				Image: singleColourImage{
					Colour: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
					Width:  50,
					Height: 50,
				},
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
		l, err := Decode(memio.Open(buf))
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
			continue
		}
		if err := compareLayers(l, test); err != nil {
			t.Errorf("test %d: %s", n+1, err)
		}
	}
}
