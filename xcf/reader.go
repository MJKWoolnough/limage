package xcf

import (
	"errors"
	"io"
	"unicode/utf8"

	"vimagination.zapto.org/byteio"
)

type reader struct {
	*byteio.StickyBigEndianReader
	rs *io.SectionReader
}

func newReader(r io.ReaderAt) reader {
	nr := reader{
		rs: io.NewSectionReader(r, 0, maxString),
	}

	nr.StickyBigEndianReader = &byteio.StickyBigEndianReader{Reader: nr.rs}

	return nr
}

const maxString = 16 * 1024 * 1024

func (r *reader) ReadString() string {
	length := r.ReadUint32()
	if length == 0 {
		return ""
	} else if length > maxString {
		r.SetError(ErrStringTooLong)

		return ""
	}

	b := make([]byte, length)

	if _, err := io.ReadFull(r, b); err != nil {
		r.Err = err

		return ""
	}

	if b[length-1] != 0 || !utf8.Valid(b[:length-1]) {
		r.SetError(ErrInvalidString)

		return ""
	}

	return string(b[:length-1])
}

func (r *reader) ReadByte() byte {
	return r.ReadUint8()
}

func (r *reader) Goto(n uint64) {
	r.rs.Seek(int64(n), io.SeekStart)
}

func (r *reader) SetError(err error) {
	if r.Err == nil {
		r.Err = err
	}
}

// Errors.
var (
	ErrInvalidString = errors.New("string is invalid")
	ErrStringTooLong = errors.New("string exceeds maximum length")
	ErrInvalidSeek   = errors.New("invalid seek")
)
