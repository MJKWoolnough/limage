package xcf

import (
	"image"
	"image/color"

	"vimagination.zapto.org/byteio"
	"vimagination.zapto.org/limage/lcolor"
	"vimagination.zapto.org/memio"
)

type compressedImage struct {
	tiles        [][]byte
	sep          [4]byte
	width        int
	tile         int
	decompressed [64 * 64 * 4]byte
}

func (c *compressedImage) decompressTile(x, y, bpp int) int {
	tile := (y/64)*((c.width+63)/64) + (x / 64)
	if tile != c.tile {
		var data memio.Buffer
		r := rle{Reader: &byteio.StickyBigEndianReader{Reader: &data}}
		for i := 0; i < bpp; i++ {
			data = c.tiles[tile][c.sep[i]:c.sep[i+1]]
			r.Read(c.decompressed[64*64*i : 64*64*(i+1)])
		}
	}
	if x < c.width & ^63 {
		return 64*(y%64) + (x % 64)
	}
	return (c.width&63)*(y%64) + (x % 64)
}

type CompressedRGB struct {
	compressedImage
	Rect image.Rectangle
}

func (CompressedRGB) ColorModel() color.Model { return lcolor.RGBModel }

func (c *CompressedRGB) Bounds() image.Rectangle { return c.Rect }

func (c *CompressedRGB) At(x, y int) color.Color { return c.RGBAt(x, y) }

func (c *CompressedRGB) RGBAt(x, y int) lcolor.RGB {
	if !(image.Point{x, y}).In(c.Rect) {
		return lcolor.RGB{}
	}
	p := c.decompressTile(x, y, 3)
	return lcolor.RGB{
		c.decompressed[p],
		c.decompressed[p+64*64],
		c.decompressed[p+64*64*2],
	}
}

type CompressedNRGBA struct {
	compressedImage
	Rect image.Rectangle
}

func (CompressedNRGBA) ColorModel() color.Model { return color.NRGBAModel }

func (c *CompressedNRGBA) Bounds() image.Rectangle { return c.Rect }

func (c *CompressedNRGBA) At(x, y int) color.Color { return c.NRGBAAt(x, y) }

func (c *CompressedNRGBA) NRGBAAt(x, y int) color.NRGBA {
	if !(image.Point{x, y}).In(c.Rect) {
		return color.NRGBA{}
	}
	p := c.decompressTile(x, y, 4)
	return color.NRGBA{
		c.decompressed[p],
		c.decompressed[p+64*64],
		c.decompressed[p+64*64*2],
		c.decompressed[p+64*64*3],
	}
}

type CompressedGray struct {
	compressedImage
	Rect image.Rectangle
}

func (CompressedGray) ColorModel() color.Model { return color.GrayModel }

func (c *CompressedGray) Bounds() image.Rectangle { return c.Rect }

func (c *CompressedGray) At(x, y int) color.Color { return c.GrayAt(x, y) }

func (c *CompressedGray) GrayAt(x, y int) color.Gray {
	if !(image.Point{x, y}).In(c.Rect) {
		return color.Gray{}
	}
	p := c.decompressTile(x, y, 1)
	return color.Gray{
		c.decompressed[p],
	}
}

type CompressedGrayAlpha struct {
	compressedImage
	Rect image.Rectangle
}

func (CompressedGrayAlpha) ColorModel() color.Model { return lcolor.GrayAlphaModel }

func (c *CompressedGrayAlpha) Bounds() image.Rectangle { return c.Rect }

func (c *CompressedGrayAlpha) At(x, y int) color.Color { return c.GrayAlphaAt(x, y) }

func (c *CompressedGrayAlpha) GrayAlphaAt(x, y int) lcolor.GrayAlpha {
	if !(image.Point{x, y}).In(c.Rect) {
		return lcolor.GrayAlpha{}
	}
	p := c.decompressTile(x, y, 2)
	return lcolor.GrayAlpha{
		c.decompressed[p],
		c.decompressed[p+64*64],
	}
}

type CompressedPaletted struct {
	compressedImage
	Rect    image.Rectangle
	Palette color.Palette
}

func (c *CompressedPaletted) ColorModel() color.Model { return c.Palette }

func (c *CompressedPaletted) Bounds() image.Rectangle { return c.Rect }

func (c *CompressedPaletted) At(x, y int) color.Color {
	if c.Palette == nil {
		return nil
	}
	if !(image.Point{x, y}).In(c.Rect) {
		return color.Gray{}
	}
	p := c.decompressTile(x, y, 1)
	i := c.decompressed[p]
	r, g, b, _ := c.Palette[i].RGBA()
	return color.NRGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: 255,
	}
}

type CompressedPalettedAlpha struct {
	compressedImage
	Rect    image.Rectangle
	Palette lcolor.AlphaPalette
}

func (c *CompressedPalettedAlpha) ColorModel() color.Model { return c.Palette }

func (c *CompressedPalettedAlpha) Bounds() image.Rectangle { return c.Rect }

func (c *CompressedPalettedAlpha) At(x, y int) color.Color {
	if c.Palette == nil {
		return nil
	}
	if !(image.Point{x, y}).In(c.Rect) {
		return color.Gray{}
	}
	p := c.decompressTile(x, y, 2)
	r, g, b, _ := c.Palette[c.decompressed[p]].RGBA()
	return color.NRGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: c.decompressed[64*64+p],
	}
}
