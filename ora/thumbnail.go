package ora

import (
	"image"
	"image/color"
)

type thumbnail struct {
	image.Image
	scale float64
}

func (t thumbnail) Bounds() image.Rectangle {
	b := t.Image.Bounds()
	b.Min.X = int(float64(b.Min.X) / t.scale)
	b.Min.Y = int(float64(b.Min.Y) / t.scale)
	b.Max.X = int(float64(b.Max.X) / t.scale)
	b.Max.Y = int(float64(b.Max.Y) / t.scale)
	return b
}

func (t thumbnail) At(x, y int) color.Color {
	return t.Image.At(int(float64(x)/t.scale), int(float64(x)/t.scale))
}
