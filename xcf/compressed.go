package xcf

import (
	"image"
	"image/color"
	"io"

	"vimagination.zapto.org/byteio"
	"vimagination.zapto.org/limage/lcolor"
	"vimagination.zapto.org/memio"
)

type compressedImage struct {
	tiles        [][][]byte
	width        int
	tile         int
	decompressed [64 * 64 * 4]byte
}

func (c *compressedImage) decompressTile(x, y int) int {
	tile := (y/64)*((c.width+63)/64) + (x / 64)
	if tile != c.tile {
		var (
			n    int
			data memio.Buffer
		)
		r := rle{Reader: &byteio.StickyBigEndianReader{Reader: &data}}
		for n, data = range c.tiles[tile] {
			r.Read(c.decompressed[64*64*n : 64*64*(n+1)])
			if r.Reader.Err == io.EOF {
				r.Reader.Err = nil
			}
		}
		c.tile = tile
	}
	if x < c.width & ^63 {
		return 64*(y%64) + (x % 64)
	}
	return (c.width&63)*(y%64) + (x % 64)
}

// CompressedRGB is an image.Image for which the data remains in a compressed
// form until read.
type CompressedRGB struct {
	compressedImage
	Rect image.Rectangle
}

// ColorModel returns the RGB Color Model
func (CompressedRGB) ColorModel() color.Model { return lcolor.RGBModel }

// Bounds returns a Rect containg the boundary data for the image
func (c *CompressedRGB) Bounds() image.Rectangle { return c.Rect }

// At returns colour at the specified coords
func (c *CompressedRGB) At(x, y int) color.Color { return c.RGBAt(x, y) }

// RGBAt returns RGB colour at the specified coords
func (c *CompressedRGB) RGBAt(x, y int) lcolor.RGB {
	if !(image.Point{x, y}).In(c.Rect) {
		return lcolor.RGB{}
	}
	p := c.decompressTile(x, y)
	return lcolor.RGB{
		c.decompressed[p],
		c.decompressed[p+64*64],
		c.decompressed[p+64*64*2],
	}
}

// CompressedNRGB is an image.Image for which the data remains in a compressed
// form until read.
type CompressedNRGBA struct {
	compressedImage
	Rect image.Rectangle
}

// ColorModel returns the NRGBA Color Model
func (CompressedNRGBA) ColorModel() color.Model { return color.NRGBAModel }

// Bounds returns a Rect containg the boundary data for the image
func (c *CompressedNRGBA) Bounds() image.Rectangle { return c.Rect }

// At returns colour at the specified coords
func (c *CompressedNRGBA) At(x, y int) color.Color { return c.NRGBAAt(x, y) }

// NRGBAAt returns NRGBA colour at the specified coords
func (c *CompressedNRGBA) NRGBAAt(x, y int) color.NRGBA {
	if !(image.Point{x, y}).In(c.Rect) {
		return color.NRGBA{}
	}
	p := c.decompressTile(x, y)
	return color.NRGBA{
		c.decompressed[p],
		c.decompressed[p+64*64],
		c.decompressed[p+64*64*2],
		c.decompressed[p+64*64*3],
	}
}

// CompressedGray is an image.Image for which the data remains in a compressed
// form until read.
type CompressedGray struct {
	compressedImage
	Rect image.Rectangle
}

// ColorModel returns the Gray Color Model
func (CompressedGray) ColorModel() color.Model { return color.GrayModel }

// Bounds returns a Rect containg the boundary data for the image
func (c *CompressedGray) Bounds() image.Rectangle { return c.Rect }

// At returns colour at the specified coords
func (c *CompressedGray) At(x, y int) color.Color { return c.GrayAt(x, y) }

// GrayAt returns Gray colour at the specified coords
func (c *CompressedGray) GrayAt(x, y int) color.Gray {
	if !(image.Point{x, y}).In(c.Rect) {
		return color.Gray{}
	}
	p := c.decompressTile(x, y)
	return color.Gray{
		c.decompressed[p],
	}
}

// CompressedGrayAlpha is an image.Image for which the data remains in a
// compressed form until read.
type CompressedGrayAlpha struct {
	compressedImage
	Rect image.Rectangle
}

// ColorModel returns the Gray Alpha Color Model
func (CompressedGrayAlpha) ColorModel() color.Model { return lcolor.GrayAlphaModel }

// Bounds returns a Rect containg the boundary data for the image
func (c *CompressedGrayAlpha) Bounds() image.Rectangle { return c.Rect }

// At returns colour at the specified coords
func (c *CompressedGrayAlpha) At(x, y int) color.Color { return c.GrayAlphaAt(x, y) }

// GrayAlphaAt returns Gray+Alpha colour at the specified coords
func (c *CompressedGrayAlpha) GrayAlphaAt(x, y int) lcolor.GrayAlpha {
	if !(image.Point{x, y}).In(c.Rect) {
		return lcolor.GrayAlpha{}
	}
	p := c.decompressTile(x, y)
	return lcolor.GrayAlpha{
		c.decompressed[p],
		c.decompressed[p+64*64],
	}
}

// CompressedPaletted is an image.Image for which the data remains in a
// compressed form until read.
type CompressedPaletted struct {
	compressedImage
	Rect    image.Rectangle
	Palette color.Palette
}

// ColorModel returns the Pallette of the image
func (c *CompressedPaletted) ColorModel() color.Model { return c.Palette }

// Bounds returns a Rect containg the boundary data for the image
func (c *CompressedPaletted) Bounds() image.Rectangle { return c.Rect }

// At returns colour at the specified coords
func (c *CompressedPaletted) At(x, y int) color.Color {
	if c.Palette == nil {
		return nil
	}
	if !(image.Point{x, y}).In(c.Rect) {
		return color.Gray{}
	}
	p := c.decompressTile(x, y)
	i := c.decompressed[p]
	r, g, b, _ := c.Palette[i].RGBA()
	return color.NRGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: 255,
	}
}

// CompressedPalettedAlpha is an image.Image for which the data remains in a
// compressed form until read.
type CompressedPalettedAlpha struct {
	compressedImage
	Rect    image.Rectangle
	Palette lcolor.AlphaPalette
}

// ColorModel returns the Pallette of the image
func (c *CompressedPalettedAlpha) ColorModel() color.Model { return c.Palette }

// Bounds returns a Rect containg the boundary data for the image
func (c *CompressedPalettedAlpha) Bounds() image.Rectangle { return c.Rect }

// At returns colour at the specified coords
func (c *CompressedPalettedAlpha) At(x, y int) color.Color {
	if c.Palette == nil {
		return nil
	}
	if !(image.Point{x, y}).In(c.Rect) {
		return color.Gray{}
	}
	p := c.decompressTile(x, y)
	r, g, b, _ := c.Palette[c.decompressed[p]].RGBA()
	return color.NRGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: c.decompressed[64*64+p],
	}
}
