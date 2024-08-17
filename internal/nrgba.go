package internal

import "image/color"

func ColourToNRGBA(c color.Color) color.NRGBA64 {
	switch c := c.(type) {
	case color.NRGBA:
		var d color.NRGBA64

		d.R = uint16(c.R)
		d.R |= d.R << 8
		d.G = uint16(c.G)
		d.G |= d.G << 8
		d.B = uint16(c.B)
		d.B |= d.B << 8
		d.A = uint16(c.A)
		d.A |= d.A << 8

		return d
	case color.NRGBA64:
		return c
	}

	if n, ok := c.(interface {
		ToNRGBA() color.NRGBA64
	}); ok {
		return n.ToNRGBA()
	}

	return color.NRGBA64Model.Convert(c).(color.NRGBA64)
}
