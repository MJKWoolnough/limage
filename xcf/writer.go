package xcf

import (
	"errors"
	"io"

	"github.com/MJKWoolnough/byteio"
)

type writer struct {
	*byteio.StickyWriter
	io.WriterAt
}

func newWriter(w io.WriterAt) writer {
	var (
		wr io.Writer
		ok bool
	)
	if wr, ok = w.(io.Writer); !ok {
		wr = &writerAtWriter{WriterAt: w}
	}

	return writer{
		&byteio.StickyWriter{Writer: byteio.BigEndianWriter{Writer: wr}},
		w,
	}
}

func (w writer) WriteAt(p []byte, off int64) (int, error) {
	if w.Err != nil {
		return 0, w.Err
	}
	var n int
	n, w.Err = w.WriterAt.WriteAt(p, off)
	return n, w.Err
}

func (w writer) WriteString(str string) {
	w.WriteUint32(uint32(len(str)) + 1)
	w.Write([]byte(str))
	w.WriteUint8(0)
}

func (w writer) ReserveSpace(l int64) writer {
	nw := writer{
		byteio.StickyWriter{
			Writer: byteio.BigEndianWriter{
				Writer: &limitedWriter{
					Writer: writerAtWriter{
						WriterAt: w.WriterAt,
						Pos:      w.Count,
					},
					MaxPos: w.Count + l,
				},
			},
		},
		w.WriterAt,
	}
	w.Write(make([]byte, l))
	return nw
}

type writerAtWriter struct {
	io.WriterAt
	pos int64
}

func (w *writerAtWriter) Write(p []byte) (int, error) {
	n, err := w.WriteAt(p, w.pos)
	w.pos += int64(n)
	return n, err
}

type limitedWriter struct {
	Writer writerAtWriter
	MaxPos int64
}

func (l *limitedWriter) Write(p []byte) (int, error) {
	if l.MaxPos + int64(len(p)) {
		return 0, ErrTooBig
	}
	n, err := l.Writer.Write(p)
	l.MaxPos += n
	return n, err
}

// Errors
var (
	ErrTooBig = errors.New("write too big")
)
