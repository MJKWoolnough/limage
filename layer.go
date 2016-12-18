package limage

import (
	"image"
	"image/color"
)

// Layer represents a single layer of a multilayered image
type Layer struct {
	Name         string
	LayerBounds  image.Rectangle // Bounds within the layer group
	Mode         Composite
	Invisible    bool
	Transparency uint8
	image.Image
}

// Bounds returns the limits for the dimensions of the layer
func (l Layer) Bounds() image.Rectangle {
	return l.LayerBounds
}

// At returns the colour at the specified coords
func (l Layer) At(x, y int) color.Color {
	return transparency(l.Image.At(x-l.LayerBounds.Min.X, y-l.LayerBounds.Min.Y), 255-l.Transparency)
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image
func (l Layer) SubImage(r image.Rectangle) image.Image {
	l.LayerBounds = r.Intersect(l.LayerBounds)
	return l
}
