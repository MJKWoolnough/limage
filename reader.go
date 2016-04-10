package xcf

import (
	"errors"
	"io"

	"github.com/MJKWoolnough/byteio"
)

type reader struct {
	byteio.StickyReader
}

func newReader(r io.Reader) *reader {
	var read reader
	read.Reader = byteio.BigEndianReader{r}
	return &read
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
	if b[length+1] != 0 {
		r.err = ErrInvalidString
		return ""
	}
	return string(b[:length])
}

// Errors
var ErrInvalidString = errors.New("string is invalid")
