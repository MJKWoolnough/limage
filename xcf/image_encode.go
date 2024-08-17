package xcf

import (
	"image"
	"image/color"

	"vimagination.zapto.org/limage/lcolor"
)

func (e *encoder) WriteImage(im image.Image, colourFunc colourBufFunc, colourChannels uint8) {
	bounds := im.Bounds()

	dx := int64(bounds.Dx())
	dy := int64(bounds.Dy())

	// Hierarchy

	e.WriteUint32(uint32(dx))
	e.WriteUint32(uint32(dy))
	e.WriteUint32(uint32(colourChannels))

	e.WriteUint32(uint32(e.pos) + 8) // currPos + this pointer (4) + zero pointer (4)
	e.WriteUint32(0)

	// Level

	e.WriteUint32(uint32(dx))
	e.WriteUint32(uint32(dy))

	nx := dx >> 6 // each tile is 64 wide
	ny := dy >> 6 // each tile is 64 high

	if dx&63 > 0 { // last tile not as wide
		nx++
	}

	if dy&63 > 0 { // last tile not as high
		ny++
	}

	w := e.ReservePointerList(uint32(nx * ny))

	// Tiles

	for y := bounds.Min.Y; y < bounds.Max.Y; y += 64 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 64 {
			l := uint16(0)

			for j := y; j < y+64 && j < bounds.Max.Y; j++ {
				for i := x; i < x+64 && i < bounds.Max.X; i++ {
					colourFunc(e, im.At(i, j))

					for n := uint8(0); n < colourChannels; n++ {
						e.channelBuf[n][l] = e.colourBuf[n]
					}

					l++
				}
			}

			w.WritePointer(uint32(e.pos))

			for n := uint8(0); n < colourChannels; n++ {
				e.WriteRLE(e.channelBuf[n][:l])
			}
		}
	}
}

func (e *encoder) rgbAlphaToBuf(c color.Color) {
	rgba := color.NRGBAModel.Convert(c).(color.NRGBA)
	e.colourBuf[0] = rgba.R
	e.colourBuf[1] = rgba.G
	e.colourBuf[2] = rgba.B
	e.colourBuf[3] = rgba.A
}

func (e *encoder) grayAlphaToBuf(c color.Color) {
	ga := lcolor.GrayAlphaModel.Convert(c).(lcolor.GrayAlpha)
	e.colourBuf[0] = ga.Y
	e.colourBuf[1] = ga.A
}

func (e *encoder) grayToBuf(c color.Color) {
	e.colourBuf[0] = color.GrayModel.Convert(c).(color.Gray).Y
}

func (e *encoder) paletteAlphaToBuf(c color.Color) {
	r, g, b, a := c.RGBA()
	i := e.colourPalette.Index(lcolor.RGB{R: uint8(r), G: uint8(g), B: uint8(b)})
	e.colourBuf[0] = uint8(i)
	e.colourBuf[1] = uint8(a)
}
