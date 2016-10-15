package xcf

import (
	"image"
	"image/color"
	"io"
	"math"

	"github.com/MJKWoolnough/limage"
	"github.com/MJKWoolnough/limage/lcolor"
)

func (d *decoder) ReadImage(width, height, mode uint32) image.Image {
	twidth := d.ReadUint32()
	theight := d.ReadUint32()

	if twidth != width || theight != height {
		d.SetError(ErrInconsistantData)
		return nil
	}

	bpp := d.ReadUint32()

	switch mode {
	case 2, 4:
		if bpp != 1 {
			d.SetError(ErrInconsistantData)
			return nil
		}
	case 3, 5:
		if bpp != 2 {
			d.SetError(ErrInconsistantData)
			return nil
		}
	case 0:
		if bpp != 3 {
			d.SetError(ErrInconsistantData)
			return nil
		}
	case 1:
		if bpp != 4 {
			d.SetError(ErrInconsistantData)
			return nil
		}
	}

	lptr := d.ReadUint32()

	d.Goto(lptr)

	w := d.ReadUint32()
	h := d.ReadUint32()

	if w != width || h != height {
		d.SetError(ErrInconsistantData)
		return nil
	}

	tiles := make([]uint32, int(math.Ceil(float64(w)/64)*math.Ceil(float64(h)/64)))

	for i := range tiles {
		tiles[i] = d.ReadUint32()
	}

	if d.ReadUint32() != 0 {
		d.SetError(ErrInconsistantData)
		return nil
	}

	var (
		im       image.Image
		imReader interface {
			ReadColour(int, int, []byte)
		}
	)

	r := image.Rect(0, 0, int(width), int(height))

	switch mode {
	case 0: // rgb
		rgb := limage.NewRGB(r)
		im = rgb
		imReader = rgbImageReader{rgb}
	case 1: // rgba
		rgba := image.NewNRGBA(r)
		im = rgba
		imReader = rgbaImageReader{rgba}
	case 2: // gray
		g := image.NewGray(r)
		im = g
		imReader = grayImageReader{g}
	case 3: // gray + alpha
		ga := limage.NewGrayAlpha(r)
		im = ga
		imReader = grayAlphaImageReader{ga}
	case 4: // indexed
		in := image.NewPaletted(r, color.Palette(d.palette))
		im = in
		imReader = indexedImageReader{in}
	case 5: // indexed + alpha
		in := limage.NewPalettedAlpha(r, d.palette)
		im = in
		imReader = palettedAlphaReader{in}
	}

	var pixBuffer [64 * 64 * 4]byte

	var cr io.Reader
	if d.compression == 0 { // no compression
		cr = &d.reader
	} else { // rle
		cr = &rle{Reader: d.reader.StickyReader}
	}

	pixel := make([]byte, bpp)
	channels := make([][]byte, bpp)

	for y := uint32(0); y < height; y += 64 {
		for x := uint32(0); x < width; x += 64 {
			d.Goto(tiles[0])
			tiles = tiles[1:]
			w := width - x
			if w > 64 {
				w = 64
			}
			h := height - y
			if h > 64 {
				h = 64
			}
			n := w * h
			cr.Read(pixBuffer[:n*bpp])
			for i := uint32(0); i < bpp; i++ {
				channels[i] = pixBuffer[n*i : n*(i+1)]
			}
			for j := uint32(0); j < h; j++ {
				for i := uint32(0); i < w; i++ {
					for k := uint32(0); k < bpp; k++ {
						pixel[k] = channels[k][0]
						channels[k] = channels[k][1:]
					}
					imReader.ReadColour(int(x+i), int(y+j), pixel)
				}
			}
		}
	}
	return im
}

type colourReader interface {
	ReadByte() byte
}

type rgbaImageReader struct {
	*image.NRGBA
}

func (rgba rgbaImageReader) ReadColour(x, y int, pixel []byte) {
	rgba.SetNRGBA(x, y, color.NRGBA{
		R: pixel[0],
		G: pixel[1],
		B: pixel[2],
		A: pixel[3],
	})
}

type grayImageReader struct {
	*image.Gray
}

func (g grayImageReader) ReadColour(x, y int, pixel []byte) {
	g.SetGray(x, y, color.Gray{pixel[0]})
}

type indexedImageReader struct {
	*image.Paletted
}

func (p indexedImageReader) ReadColour(x, y int, pixel []byte) {
	p.SetColorIndex(x, y, pixel[0])
}

type grayAlphaImageReader struct {
	*limage.GrayAlpha
}

func (ga grayAlphaImageReader) ReadColour(x, y int, pixels []byte) {
	ga.SetGrayAlpha(x, y, lcolor.GrayAlpha{pixels[0], pixels[1]})
}

type palettedAlphaReader struct {
	*limage.PalettedAlpha
}

func (p palettedAlphaReader) ReadColour(x, y int, pixels []byte) {
	p.SetIndexAlpha(x, y, lcolor.IndexedAlpha{
		I: pixels[0],
		A: pixels[1],
	})
}

type rgbImageReader struct {
	*limage.RGB
}

func (rg rgbImageReader) ReadColour(x, y int, pixels []byte) {
	rg.SetRGB(x, y, lcolor.RGB{R: pixels[0], G: pixels[1], B: pixels[2]})
}
