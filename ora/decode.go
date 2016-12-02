package ora

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"image"
	"image/color"
	"io"
	"strconv"

	"github.com/MJKWoolnough/limage"
)

func getStack(zr *zip.Reader) (*zip.File, error) {
	required := 0
	var stack *zip.File
	for _, f := range zr.File {
		switch f.Name {
		case "stack.xml":
			required++
			stack = f
		case "Thumbnails/thumbnail.png", "mergedimage.png":
			required++
		case "mimetype":
			if !checkMime(f) {
				return nil, ErrInvalidMimeType
			}
			required++
		}
	}
	if required < 4 {
		return nil, ErrMissingRequired
	}
	return stack, nil
}

func DecodeConfig(zr *zip.Reader) (image.Config, error) {
	stack, err := getStack(zr)
	if err != nil {
		return image.Config{}, err
	}
	s, err := stack.Open()
	if err != nil {
		return image.Config{}, err
	}
	x := xml.NewDecoder(s)
	var width, height int
	for {
		t, err := x.Token()
		if err != nil {
			if err == io.EOF {
				return image.Config{}, ErrInvalidStack
			}
			return image.Config{}, err
		}
		if se, ok := t.(xml.StartElement); ok {
			if se.Name.Local == "image" {
				var w, h bool
				for _, attr := range se.Attr {
					switch attr.Name.Local {
					case "w":
						width, err = strconv.Atoi(attr.Value)
						w = true
					case "h":
						height, err = strconv.Atoi(attr.Value)
						h = true
					}
					if err != nil {
						return image.Config{}, err
					}
				}
				if !w || !h {
					return image.Config{}, ErrInvalidStack
				}
				break
			}
			return image.Config{}, ErrInvalidStack
		}
	}
	s.Close()
	return image.Config{
		ColorModel: color.NRGBA64Model,
		Width:      width,
		Height:     height,
	}, nil
}

func Decode(zr *zip.Reader) (*limage.Image, error) {
	stack, err := getStack(zr)
	if err != nil {
		return nil, err
	}
	s, err := stack.Open()
	if err != nil {
		return nil, err
	}
	x := xml.NewDecoder(s)
	_ = x
	s.Close()
	return nil, nil
}

func checkMime(mimetype *zip.File) bool {
	if mimetype.UncompressedSize64 != uint64(len(mimetypeStr)) {
		return false
	} else {
		mr, err := mimetype.Open()
		if err != nil {
			return false
		}
		var mime [len(mimetypeStr)]byte
		_, err = io.ReadFull(mr, mime[:])
		mr.Close()
		if err != nil {
			return false
		}
		return string(mime[:]) == mimetypeStr
	}
}

// Errors
var (
	ErrMissingRequired = errors.New("missing required file")
	ErrInvalidMimeType = errors.New("invalid mime type")
	ErrInvalidStack    = errors.New("invalid stack")
)
