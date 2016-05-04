package xcf

import (
	"errors"
	"io"
	"unicode/utf8"

	"github.com/MJKWoolnough/byteio"
)

type reader struct {
	byteio.StickyReader
	io.Seeker
}

func newReader(r io.ReadSeeker) reader {
	return reader{
		StickyReader: byteio.StickyReader{
			Reader: byteio.BigEndianReader{r},
		},
		Seeker: r,
	}
}

const maxString = 16 * 1024 * 1024

func (r *reader) ReadString() string {
	length := r.ReadUint32()
	if length == 0 {
		return ""
	}
	if length > maxString {
		if r.Err == nil {
			r.Err = ErrStringTooLong
		}
		return ""
	}
	b := make([]byte, length)
	_, err := io.ReadFull(r, b)
	if err != nil {
		r.Err = err
		return ""
	}
	if b[length-1] != 0 || !utf8.Valid(b[:length-1]) {
		r.Err = ErrInvalidString
		return ""
	}
	return string(b[:length-1])
}

func (r *reader) ReadByte() byte {
	return r.ReadUint8()
}

func (r *reader) Seek(offset int64, whence int) (int64, error) {
	if r.Err != nil {
		return 0, r.Err
	}
	n, err := r.Seeker.Seek(offset, whence)
	if err != nil {
		r.Err = err
	}
	return n, err
}

// Errors
var (
	ErrInvalidString = errors.New("string is invalid")
	ErrStringTooLong = errors.New("string exceeds maximum length")
)
