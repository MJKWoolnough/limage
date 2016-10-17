package xcf

import (
	"io"

	"github.com/MJKWoolnough/byteio"
	"github.com/MJKWoolnough/errors"
)

type writer struct {
	*byteio.StickyWriter
	*writerAtWriter
}

func newWriter(w io.WriterAt) writer {
	wr := &writerAtWriter{
		WriterAt: w,
	}
	return writer{
		StickyWriter:   &byteio.StickyWriter{Writer: byteio.BigEndianWriter{Writer: wr}},
		writerAtWriter: wr,
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

func (w writer) Write(p []byte) {
	w.StickyWriter.Write(p)
}

func (w writer) WriteString(str string) {
	w.WriteUint32(uint32(len(str)) + 1)
	w.Write([]byte(str))
	w.WriteUint8(0)
}

type pointerWriter struct {
	bw      byteio.StickyWriter
	toWrite uint32
	obw     *byteio.StickyWriter
}

func (p *pointerWriter) WritePointer(ptr uint32) {
	if p.toWrite > 0 {
		p.bw.WriteUint32(ptr)
		p.toWrite--
		if p.bw.Err != nil {
			p.obw.Err = p.bw.Err
		}
	}
}

func (w writer) ReservePointers(n uint32) *pointerWriter {
	p := &pointerWriter{
		bw: byteio.StickyWriter{
			Writer: byteio.BigEndianWriter{
				Writer: &writerAtWriter{
					WriterAt: w.writerAtWriter.WriterAt,
					pos:      w.pos,
				},
			},
		},
		toWrite: n,
		obw:     w.StickyWriter,
	}
	w.pos += int64(n) * 4
	w.WriteUint32(0)
	return p
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

// Errors
const (
	ErrTooBig errors.Error = "write too big"
)
