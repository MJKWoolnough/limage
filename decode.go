package ora

import (
	"archive/zip"
	"errors"
	"image"
	"image/png"
	"io"
)

func Decode(r io.ReaderAt) (image.Image, error) {
	zr, err := zip.NewReader(r)
	if err != nil {
		return nil, err
	}
	required := 0
	var merged *zip.File
	for _, f := range zr.File {
		switch f.Name {
		case "mimetype", "stack.xml", "data", "Thumbnails/thumbnail.png":
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
		return
	}
	defer f.Close()
	return png.Decode(f)
}

// Errors
var (
	ErrMissingRequired = errors.New("missing required file")
)
