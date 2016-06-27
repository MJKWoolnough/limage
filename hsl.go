package xcf

import "image/color"

type hsl struct {
	H, S, L, A uint16
}

func rgbToHSL(cl color.Color) hsl {
	return hsl{}
}

func (h hsl) RGBA() (uint32, uint32, uint32, uint32) {
	return 0, 0, 0, 0
}

type hsv struct {
	H, S, V, A uint16
}

func rgbToHSV(cl color.Color) hsv {
	return hsv{}
}

func (h hsv) RGBA() (uint32, uint32, uint32, uint32) {
	return 0, 0, 0, 0
}
