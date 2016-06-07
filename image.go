package xcf

import (
	"image"
	"image/color"
	"math"
	"os"
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
	}

	lptr := d.ReadUint32()

	/*
		for {
			if d.ReadUint32() == 0 { // dummy level
				break
			}
		}
	*/

	d.Seek(int64(lptr))

	w := d.ReadUint32()
	h := d.ReadUint32()

	if w != width || h != height {
		d.SetError(ErrInconsistantData)
		return l
	}

	tiles := make([]uint32, int(math.Ceil(float64(w)/64)*math.Ceil(float64(h)/64)))

	for i := range tiles {
		tiles[i] = d.ReadUint32()
	}

	if d.ReadUint32() != 0 {
		d.SetError(ErrInconsistantData)
		return l
	}

	var (
		im       image.Image
		imReader interface {
			ReadColour(int, int, colourReader)
		}
	)

	r := image.Rect(0, 0, width, height)

	switch mode {
	case 0: // rgb
		rgb := newRGB(r)
		im = rgb
		imReader = rgbImageReader{rgb}
	case 1: // rgba
		rgba := image.NewRGBA(r)
		im = rgba
		imReader = rgbaImageReader{rgba}
	case 2: // gray
		g := image.NewGray(r)
		im = g
		imReader = greyImageReader{g}
	case 3: // gray + alpha
		ga := newGrayAlpha(r)
		im = ga
		imReader = greyAlphaImageReader{ga}
	case 4: // indexed
		in := image.NewPaletted(r, d.palette)
		im = in
		imReader = indexedImageReader{in}
	case 5: // indexed + alpha
		in := newPalettedAlpha(r, d.palette)
		im = in
		imReader = indexedAlphaImageReader{in}
	}

	for y := uint32(0); y < height; y += 64 {
		for x := uint32(0); x < width; x += 64 {
			d.Seek(int64(tiles[0]), os.SEEK_SET)
			var cr colourReader
			if d.compression == 0 { // no compression
				cr = &d.reader
			} else { // rle
				cr = &rle{reader: d.reader}
			}
			for j := y; j < y+64 && j < height; j++ {
				for i := x; i < x+64 && i < width; i++ {
					imReader.ReadColour(int(i), int(j), cr)
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

func (rgba rgbaImageReader) ReadColour(x, y int, cr colourReader) {
	r := cr.ReadByte()
	g := cr.ReadByte()
	b := cr.ReadByte()
	a := cr.ReadByte()
	rgba.SetNRGBA(x, y, color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	})
}

type grayImageReader struct {
	*image.Gray
}

func (g grayImageReader) ReadColour(x, y int, cr colourReader) {
	yc := cr.ReadByte()
	g.SetGray(x, y, color.Gray{yc})
}

type indexedImageReader struct {
	*image.Paletted
}

func (p indexedImageReader) ReadColour(x, y, int, cr colourReader) {
	i := cr.ReadByte()
	p.SetColorIndex(x, y, i)
}
