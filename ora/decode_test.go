package ora

import (
	"archive/zip"
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

var buf [14485]byte

func openFile(str string) (*zip.Reader, error) {
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
	return zip.NewReader(bytes.NewReader(buf[:n]), int64(n))
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
		zr, err := openFile(test.File)
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
			continue
		}
		c, err := DecodeConfig(zr)
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
			File: blackRedBlueFile,
			Image: limage.Image{
				limage.Layer{
					Name: "Layer Group",
					Image: limage.Image{
						limage.Layer{
							Name: "Blue",
							Image: singleColourImage{
								Colour: color.RGBA64{B: 65535, A: 65535},
								Width:  30,
								Height: 30,
							},
							LayerBounds: image.Rect(0, 0, 30, 30),
						},
						limage.Layer{
							Name: "Red",
							Image: singleColourImage{Colour: color.RGBA64{R: 65535, A: 65535},
								Width:  30,
								Height: 30,
							},
							LayerBounds: image.Rect(20, 20, 50, 50),
						},
					},
					LayerBounds: image.Rect(0, 0, 50, 50),
				},
				limage.Layer{
					Name: "Background",
					Image: singleColourImage{
						Colour: color.RGBA64{A: 65535},
						Width:  50,
						Height: 50,
					},
					LayerBounds: image.Rect(0, 0, 50, 50),
				},
			},
		},
		{
			File: blackRedFile,
			Image: limage.Image{
				limage.Layer{
					Name: "Layer",
					Image: singleColourImage{
						Colour: color.NRGBA{R: 255, A: 255},
						Width:  30,
						Height: 30,
					},
					LayerBounds: image.Rect(10, 10, 40, 40),
				},
				limage.Layer{
					Name: "Background",
					Image: singleColourImage{
						Colour: color.RGBA{A: 255},
						Width:  50,
						Height: 50,
					},
					LayerBounds: image.Rect(0, 0, 50, 50),
				},
			},
		},
		{
			File: blackFile,
			Image: limage.Image{
				limage.Layer{
					Name: "Background",
					Image: singleColourImage{
						Colour: color.RGBA{A: 255},
						Width:  50,
						Height: 50,
					},
					LayerBounds: image.Rect(0, 0, 50, 50),
				},
			},
		},
		{
			File: redFile,
			Image: limage.Image{
				limage.Layer{
					Name: "Background",
					Image: singleColourImage{
						Colour: color.RGBA{R: 255, A: 255},
						Width:  50,
						Height: 50,
					},
					LayerBounds: image.Rect(0, 0, 50, 50),
				},
			},
		},
		{
			File: whiteFile,
			Image: limage.Image{
				limage.Layer{
					Name: "Background",
					Image: singleColourImage{
						Colour: color.RGBA{R: 255, G: 255, B: 255, A: 255},
						Width:  50,
						Height: 50,
					},
					LayerBounds: image.Rect(0, 0, 50, 50),
				},
			},
		},
	}

	for n, test := range tests {
		zr, err := openFile(test.File)
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
			continue
		}
		i, err := Decode(zr)
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
		ia := la.Image
		ib := lb.Image
		la.Image = nil
		lb.Image = nil
		if !reflect.DeepEqual(la, lb) {
			return fmt.Errorf("layer properties mismatched, expecting %#v, got %#v", lb, la)
		}
		if mib, ok := ib.(limage.MaskedImage); ok {
			if mia, ok := ia.(limage.MaskedImage); ok {
				if err := compareImages(mia.Mask, mib.Mask); err != nil {
					return err
				}
				ia = mia.Image
				ib = mib.Image
			} else {
				return fmt.Errorf("expecting MaskedImage, got %T", ia)
			}
		}
		if layb, ok := ib.(limage.Image); ok {
			if laya, ok := ia.(limage.Image); ok {
				if err := compareLayers(laya, layb); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("expecting Layer Group, got %T", ia)
			}
		} else if tb, ok := ib.(limage.Text); ok {
			if ta, ok := ia.(limage.Text); ok {
				ta.Image = nil
				tb.Image = nil
				if !reflect.DeepEqual(ta, tb) {
					return fmt.Errorf("expecting text layer %#v, got %#v", tb, ta)
				}
			} else {
				return fmt.Errorf("expecting Text Layer, got %T", ia)
			}
		} else if err := compareImages(ia, ib); err != nil {
			return err
		}
	}
	return nil
}

func compareImages(ia, ib image.Image) error {
	bnds := ia.Bounds()
	bndsb := ib.Bounds()
	if bnds.Min.X != bndsb.Min.X || bnds.Min.Y != bndsb.Min.Y || bnds.Max.X != bndsb.Max.X || bnds.Max.Y != bndsb.Max.Y {
		return fmt.Errorf("bounds mismatch, expecting %v, got %v", bndsb, bnds)
	}
	for j := 0; j < bnds.Dy(); j++ {
		for i := 0; i < bnds.Dx(); i++ {
			ca := ia.At(i, j)
			cb := ib.At(i, j)
			if !reflect.DeepEqual(ca, cb) {
				return fmt.Errorf("pixel mismatch: expecting %#v, got %#v", cb, ca)
			}
		}
	}
	return nil
}
