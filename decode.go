package ora

import (
	"archive/zip"
	"errors"
	"image"
	"image/png"
	"io"
)

func Decode(r io.ReaderAt, size int64) (image.Image, error) {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	required := 0
	var merged *zip.File
	for _, f := range zr.File {
		switch f.Name {
		case "stack.xml", "data", "Thumbnails/thumbnail.png":
			required++
		case "mimetype":
			if f.UncompressedSize64 != uint64(len(mimetypeStr)) {
				return nil, ErrInvalidMimeType
			} else {
				mr, err := f.Open()
				if err != nil {
					return nil, err
				}
				var mime [16]byte
				_, err = io.ReadFull(mr, mime[:])
				mr.Close()
				if err != nil {
					return nil, err
				}
				if string(mime[:]) != mimetypeStr {
					return nil, ErrInvalidMimeType
				}
			}
			required++
		case "mergedimage.png":
			merged = f
			required++
		}
	}
	if required < 5 {
		return nil, ErrMissingRequired
	}
	f, err := merged.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}

// Errors
var (
	ErrMissingRequired = errors.New("missing required file")
	ErrInvalidMimeType = errors.New("invalid mime type")
)
