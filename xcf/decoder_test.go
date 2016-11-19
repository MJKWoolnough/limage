package xcf

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"image"
	"image/color"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/MJKWoolnough/limage"
)

var buf [2098]byte

func openFile(str string) (io.ReaderAt, error) {
	gz, err := gzip.NewReader(strings.NewReader(str))
	if err != nil {
		return nil, err
	}
	n, err := gz.Read(buf[:])
	if err != nil {
		if err != io.EOF {
			return nil, err
		}
	}
	return bytes.NewReader(buf[:n]), nil
}

func TestConfigDecoder(t *testing.T) {
	tests := []struct {
		File   string
		Config image.Config
	}{
		{
			File: abcFile,
			Config: image.Config{
				ColorModel: color.NRGBAModel,
				Width:      36,
				Height:     13,
			},
		},
		{
			File: blackMaskFile,
			Config: image.Config{
				ColorModel: color.NRGBAModel,
				Width:      50,
				Height:     50,
			},
		},
		{
			File: blackRedBlueFile,
			Config: image.Config{
				ColorModel: color.NRGBAModel,
				Width:      50,
				Height:     50,
			},
		},
		{
			File: blackRedFile,
			Config: image.Config{
				ColorModel: color.NRGBAModel,
				Width:      50,
				Height:     50,
			},
		},
		{
			File: blackFile,
			Config: image.Config{
				ColorModel: color.NRGBAModel,
				Width:      50,
				Height:     50,
			},
		},
		{
			File: redFile,
			Config: image.Config{
				ColorModel: color.NRGBAModel,
				Width:      50,
				Height:     50,
			},
		},
		{
			File: whiteFile,
			Config: image.Config{
				ColorModel: color.NRGBAModel,
				Width:      50,
				Height:     50,
			},
		},
	}

	for n, test := range tests {
		f, err := openFile(test.File)
		if err != nil {
			t.Errorf("test %d: unexpected error opening file: %s", n+1, err)
			continue
		}
		c, err := DecodeConfig(f)
		if err != nil {
			t.Errorf("test %d: unexpected error decoding config: %s", n+1, err)
			continue
		}
		if !reflect.DeepEqual(test.Config, c) {
			t.Errorf("test %d: no config match", n+1)
		}
	}
}

type singleColourImage struct {
	Colour        color.Color
	Width, Height int
}

func (s singleColourImage) ColorModel() color.Model {
	return s
}

func (s singleColourImage) Convert(color.Color) color.Color {
	return s.Colour
}

func (s singleColourImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, s.Width, s.Height)
}

func (s singleColourImage) At(int, int) color.Color {
	return s.Colour
}

func TestDecoder(t *testing.T) {
	tests := []struct {
		File  string
		Image limage.Image
	}{
		{
			File: redFile,
			Image: limage.Image{
				limage.Layer{
					Name:  "Layer",
					Image: singleColourImage{Colour: color.RGBA{R: 255}},
				},
			},
		},
	}

	for n, test := range tests {
		f, err := openFile(test.File)
		if err != nil {
			t.Errorf("test %d: unexpected error opening file: %s", n+1, err)
			continue
		}
		i, err := Decode(f)
		if err != nil {
			t.Errorf("test %d: unexpected error decoding image: %s", n+1, err)
			continue
		}
		if err := compareLayers(i, test.Image); err != nil {
			t.Errorf("test %d: %s", n+1, err)
		}
	}
}

func compareLayers(a, b limage.Image) error {
	if len(a) != len(b) {
		return fmt.Errorf("incorrect number of layers, expecting %d, got %d", len(b), len(a))
	}
	for n, la := range a {
		lb := b[n]
		if la.Invisible != lb.Invisible {

		}
		if la.Mode != lb.Mode {

		}
		if la.Transparency != lb.Transparency {

		}
		if la.OffsetX != lb.OffsetX {

		}
		if la.OffsetY != lb.OffsetY {

		}
		if la.Name != lb.Name {

		}
		if ba, bb := la.Bounds(), b.Bounds(); ba != bb {

		}
		// compare images
	}
	return nil
}
