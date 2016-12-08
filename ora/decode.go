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
	zr     *zip.Reader
	x      *xml.Decoder
	limits image.Point
}

func (d *decoder) getStack() (stack *zip.File, err error) {
	err = ErrMissingStack
	for _, f := range d.zr.File {
		switch f.Name {
		case "stack.xml":
			stack = f
			err = nil
		case "mimetype":
			if !checkMime(f) {
				return nil, ErrInvalidMimeType
			}
		}
	}
	return stack, err
}

func (d *decoder) getDimensions() error {
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
						d.limits.X, err = strconv.Atoi(attr.Value)
						w = true
					case "h":
						d.limits.Y, err = strconv.Atoi(attr.Value)
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
		return image.Config{}, err
	}
	s.Close()
	return image.Config{
		ColorModel: color.NRGBA64Model,
		Width:      d.w,
		Height:     d.h,
	}, nil
}

func Decode(zr *zip.Reader) (limage.Image, error) {
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
	for { // skip to first stack tag
		t, err := d.x.Token()
		if err != nil {
			if err == io.EOF {
				return nil, io.ErrUnexpectedEOF
			}
			return nil, err
		}
		if se, ok := t.(*xml.StartElement); ok {
			if se.Name.Local == "stack" {
				break
			}
			d.skipTag()
		}
	}
	return d.readStack(image.Point{})
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
	ErrMissingStack    = errors.New("missing stack file")
	ErrInvalidMimeType = errors.New("invalid mime type")
	ErrInvalidStack    = errors.New("invalid stack")
)
