package xcf

import (
	"io"
	"unicode/utf8"

	"github.com/MJKWoolnough/byteio"
	"github.com/MJKWoolnough/errors"
)

type reader struct {
	byteio.StickyReader
	io.Seeker
	rs readSeeker
}

type readSeeker struct {
	io.ReaderAt
	pos int64
}

func (r *readSeeker) Read(p []byte) (int, error) {
	n, err := r.ReadAt(p, r.pos)
	r.pos += int64(n)
	return n, err
}

func (r *readSeeker) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		r.pos = offset
	case io.SeekCurrent:
		r.pos += offset
	default:
		return 0, ErrInvalidSeek
	}
	return r.pos, nil
}

func newReader(r io.ReaderAt) reader {
	nr := reader{
		rs: readSeeker{ReaderAt: r},
	}
	nr.StickyReader.Reader = byteio.BigEndianReader{&nr.rs}
	nr.Seeker = &nr.rs
	return nr
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

func (r *reader) Goto(n uint32) {
	if r.Err != nil {
		return
	}
	_, r.Err = r.Seeker.Seek(int64(n), io.SeekStart)
}

func (r *reader) Skip(n uint32) {
	if r.Err != nil {
		return
	}
	_, r.Err = r.Seeker.Seek(int64(n), io.SeekCurrent)
}

// Errors
var (
	ErrInvalidString errors.Error = "string is invalid"
	ErrStringTooLong errors.Error = "string exceeds maximum length"
	ErrInvalidSeek   errors.Error = "invalid seek"
)
