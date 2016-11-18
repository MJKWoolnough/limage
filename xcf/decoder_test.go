package xcf

import (
	"bytes"
	"compress/gzip"
	"image"
	"image/color"
	"io"
	"reflect"
	"strings"
	"testing"
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

func TestDecoder(t *testing.T) {
}
