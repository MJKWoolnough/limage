package xcf

import (
	"errors"
	"image"
	"io"
)

const (
	fileTypeID   = "gimp xcf "
	fileVersion0 = "file"
	fileVersion1 = "v001"
	fileVersion2 = "v002"
	fileVersion3 = "v003"
)

type decoder struct {
	reader
	width, height, baseType uint32
}

func Decoder(r io.ReadSeeker) (image.Image, error) {
	d := decoder{reader: newReader(r)}

	// check header

	var header [14]byte
	d.Read(header[:])
	if d.Err != nil {
		return nil, d.Err // wrap?
	}
	if string(header[:9]) != fileTypeID {
		return nil, ErrInvalidFileTypeID
	}
	switch string(header[9:13]) {
	case fileVersion0, fileVersion1, fileVersion2, fileVersion3:
	default:
		return nil, ErrUnsupportedVersion
	}
	if header[14] != 0 {
		return nil, ErrInvalidHeader
	}
	d.width = d.ReadUint32()
	d.height = d.ReadUint32()
	d.baseType = d.ReadUint32()

	return nil, nil
}

// Errors
var (
	ErrInvalidFileTypeID  = errors.New("invalid file type identification")
	ErrUnsupportedVersion = errors.New("unsupported file version")
	ErrInvalidHeader      = errors.New("invalid header")
)
