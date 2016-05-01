package xcf

import (
	"errors"
	"io"

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
			Seeker: r,
		},
	}
}

func (r *reader) ReadString() string {
	length := r.ReadUint32()
	if length == 0 {
		return ""
	}
	b := make([]byte, length+1)
	_, err := io.ReadFull(r, b)
	if err != nil {
		r.Err = err
		return ""
	}
	if b[length] != 0 {
		r.Err = ErrInvalidString
		return ""
	}
	return string(b[:length])
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
var ErrInvalidString = errors.New("string is invalid")
