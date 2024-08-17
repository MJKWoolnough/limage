package xcf

import (
	"image"
	"image/color"
	"io"
	"math"

	"vimagination.zapto.org/limage"
	"vimagination.zapto.org/limage/lcolor"
	"vimagination.zapto.org/memio"
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

	var lptr uint64

	if d.mode < 2 {
		lptr = uint64(d.ReadUint32())
	} else {
		lptr = d.ReadUint64()
	}

	d.Goto(lptr)

	w := d.ReadUint32()
	h := d.ReadUint32()

	if w != width || h != height {
		d.SetError(ErrInconsistantData)

		return nil
	}

	tiles := make([]uint64, int(math.Ceil(float64(w)/64)*math.Ceil(float64(h)/64)))

	if d.mode < 2 {
		for i := range tiles {
			tiles[i] = uint64(d.ReadUint32())
		}
	} else {
		for i := range tiles {
			tiles[i] = d.ReadUint64()
		}
	}

	if d.ReadUint32() != 0 {
		d.SetError(ErrInconsistantData)
		return nil
	}

	r := image.Rect(0, 0, int(width), int(height))

	var pixBuffer [64 * 64 * 4]byte

	if d.decompress || d.compression == 0 {
		var (
			im       image.Image
			imReader interface {
				ReadColour(int, int, []byte)
			}
		)

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

		var cr io.Reader

		if d.compression == 0 { // no compression
			cr = &d.reader
		} else { // rle
			cr = &rle{Reader: d.reader.StickyBigEndianReader}
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
				_, err := cr.Read(pixBuffer[:n*bpp])

				d.SetError(err)

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
	} else {
		ci := compressedImage{
			tiles: make([][][]byte, 0, len(tiles)),
			width: int(width),
			tile:  -1,
		}

		buf := make(memio.Buffer, 0, 64*64*4)

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
				ts := make([][]byte, 0, bpp)

				for i := uint32(0); i < bpp; i++ {
					d.SetError(d.readRLE(int(n), &buf))

					b := make([]byte, len(buf))

					copy(b, buf)

					buf = buf[:0]
					ts = append(ts, b)
				}

				ci.tiles = append(ci.tiles, ts)
			}
		}

		switch mode {
		case 0: // rgb
			return &CompressedRGB{ci, r}
		case 1: // rgba
			return &CompressedNRGBA{ci, r}
		case 2: // gray
			return &CompressedGray{ci, r}
		case 3: // gray + alpha
			return &CompressedGrayAlpha{ci, r}
		case 4: // indexed
			return &CompressedPaletted{ci, r, color.Palette(d.palette)}
		case 5: // indexed + alpha
			return &CompressedPalettedAlpha{ci, r, d.palette}
		default:
			return nil
		}
	}
}

/*
type colourReader interface {
	ReadByte() byte
}
*/

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
	ga.SetGrayAlpha(x, y, lcolor.GrayAlpha{Y: pixels[0], A: pixels[1]})
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
