package xcf

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
	"reflect"
	"testing"

	"vimagination.zapto.org/limage"
	"vimagination.zapto.org/memio"
)

func imageRandom(r image.Rectangle) image.Image {
	i := image.NewNRGBA(r)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			ci := rand.Uint32()
			i.Set(x, y, color.NRGBA{
				R: uint8(ci >> 24),
				G: uint8(ci >> 16),
				B: uint8(ci >> 8),
				A: uint8(ci),
			})
		}
	}
	return i
}

func TestCompressedImages(t *testing.T) {
	for n, test := range [...]limage.Image{
		{
			limage.Layer{
				Name:        "Layer 1",
				LayerBounds: image.Rect(0, 0, 64, 64),
				Image:       image.NewNRGBA(image.Rect(0, 0, 64, 64)),
			},
		},
		{
			limage.Layer{
				Name:        "Layer 1",
				LayerBounds: image.Rect(0, 0, 100, 100),
				Image:       image.NewNRGBA(image.Rect(0, 0, 100, 100)),
			},
		},
		{
			limage.Layer{
				Name:        "Layer 1",
				LayerBounds: image.Rect(0, 0, 10, 10),
				Image:       imageRandom(image.Rect(0, 0, 10, 10)),
			},
		},
		{
			limage.Layer{
				Name:        "Layer 1",
				LayerBounds: image.Rect(0, 0, 64, 64),
				Image:       imageRandom(image.Rect(0, 0, 64, 64)),
			},
		},
		{
			limage.Layer{
				Name:        "Layer 1",
				LayerBounds: image.Rect(0, 0, 100, 100),
				Image:       imageRandom(image.Rect(0, 0, 100, 100)),
			},
		},
	} {
		g := image.NewNRGBA(test.Bounds())
		draw.Draw(g, g.Rect, test, image.Point{}, draw.Over)
		f, _ := os.Create(fmt.Sprintf("%d-o.png", n))
		png.Encode(f, g)
		f.Close()
		var buf []byte
		Encode(memio.Create(&buf), &test)
		tl, err := DecodeCompressed(memio.Open(buf))
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
			continue
		}
		gt := image.NewNRGBA(tl.Bounds())
		draw.Draw(gt, g.Rect, tl, image.Point{}, draw.Over)
		f, _ = os.Create(fmt.Sprintf("%d-p.png", n))
		png.Encode(f, gt)
		f.Close()
		if !reflect.DeepEqual(g, gt) {
			t.Errorf("test %d: output does not match test", n+1)
		}
	}
}
