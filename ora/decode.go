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

type decoder struct {
	zr   *zip.Reader
	x    *xml.Decoder
	w, h int
}

func (d decoder) getStack() (*zip.File, error) {
	required := 0
	var stack *zip.File
	for _, f := range d.zr.File {
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

func (d decoder) getDimensions() error {
	for {
		t, err := d.x.Token()
		if err != nil {
			if err == io.EOF {
				return ErrInvalidStack
			}
			return err
		}
		if se, ok := t.(xml.StartElement); ok {
			if se.Name.Local == "image" {
				var w, h bool
				for _, attr := range se.Attr {
					switch attr.Name.Local {
					case "w":
						d.w, err = strconv.Atoi(attr.Value)
						w = true
					case "h":
						d.h, err = strconv.Atoi(attr.Value)
						h = true
					}
					if err != nil {
						return err
					}
				}
				if !w || !h {
					return ErrInvalidStack
				}
				return nil
			}
			return ErrInvalidStack
		}
	}
}

func DecodeConfig(zr *zip.Reader) (image.Config, error) {
	d := decoder{zr: zr}
	stack, err := d.getStack()
	if err != nil {
		return image.Config{}, err
	}
	s, err := stack.Open()
	if err != nil {
		return image.Config{}, err
	}
	d.x = xml.NewDecoder(s)
	if err := d.getDimensions(); err != nil {
		return iamge.Config{}, err
	}
	s.Close()
	return image.Config{
		ColorModel: color.NRGBA64Model,
		Width:      d.w,
		Height:     d.h,
	}, nil
}

func Decode(zr *zip.Reader) (*limage.Image, error) {
	d := decoder{zr: zr}
	stack, err := d.getStack()
	if err != nil {
		return nil, err
	}
	s, err := stack.Open()
	if err != nil {
		return nil, err
	}
	defer s.Close()
	d.x = xml.NewDecoder(s)
	if err := d.getDimensions(); err != nil {
		return nil, err
	}
	return d.readStack()
}

func checkMime(mimetype *zip.File) bool {
	if mimetype.UncompressedSize64 != uint64(len(mimetypeStr)) {
		return false
	}
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

// Errors
var (
	ErrMissingRequired = errors.New("missing required file")
	ErrInvalidMimeType = errors.New("invalid mime type")
	ErrInvalidStack    = errors.New("invalid stack")
)
