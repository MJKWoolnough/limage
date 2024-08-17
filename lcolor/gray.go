package lcolor

import "image/color"

// GrayAlpha represents a Gray color with an Alpha channel.
type GrayAlpha struct {
	Y, A uint8
}

// RGBA implements the color.Color interface.
func (c GrayAlpha) RGBA() (r, g, b, a uint32) {
	y := uint32(c.Y)
	y |= y << 8
	a = uint32(c.A)
	y *= a
	y /= 0xff
	a |= a << 8

	return y, y, y, a
}

// ToNRGBA converts the HSV color into the RGB colorspace.
func (c GrayAlpha) ToNRGBA() color.NRGBA64 {
	y := uint16(c.Y)
	y |= y << 8
	a := uint16(c.A)
	a |= a << 8

	return color.NRGBA64{y, y, y, a}
}

func grayAlphaColourModel(c color.Color) color.Color {
	_, _, _, a := c.RGBA()

	return GrayAlpha{
		Y: color.GrayModel.Convert(c).(color.Gray).Y,
		A: uint8(a >> 8),
	}
}
