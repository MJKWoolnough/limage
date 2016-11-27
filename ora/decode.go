package ora

import (
	"archive/zip"
	"errors"
	"image"
	"io"

	"github.com/MJKWoolnough/limage"
)

func DecodeConfig(zr zip.Reader) (image.Config, error) {
	return image.Config{}, nil
}

func Decode(zr zip.Reader) (*limage.Image, error) {
	required := 0
	var stack *zip.File
	for _, f := range zr.File {
		switch f.Name {
		case "stack.xml":
			required++
			stack = f
		case "data", "Thumbnails/thumbnail.png", "mergedimage.png":
			required++
		case "mimetype":
			if !checkMime(f) {
				return nil, ErrInvalidMimeType
			}
			required++
		}
	}
	if required < 5 {
		return nil, ErrMissingRequired
	}
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
		var mime [16]byte
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
)
