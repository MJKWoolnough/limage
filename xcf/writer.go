package xcf

import (
	"errors"
	"io"

	"vimagination.zapto.org/byteio"
)

type writer struct {
	*byteio.StickyBigEndianWriter
	*writerAtWriter
}

func newWriter(w io.WriterAt) writer {
	wr := &writerAtWriter{
		WriterAt: w,
	}
	return writer{
		StickyBigEndianWriter: &byteio.StickyBigEndianWriter{Writer: wr},
		writerAtWriter:        wr,
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
	w.StickyBigEndianWriter.Write(p)
}

func (w writer) WriteString(str string) {
	w.WriteUint32(uint32(len(str)) + 1)
	w.Write([]byte(str))
	w.WriteUint8(0)
}

type pointerWriter struct {
	bw      *byteio.StickyBigEndianWriter
	toWrite uint32
	obw     *byteio.StickyBigEndianWriter
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
		bw: &byteio.StickyBigEndianWriter{
			Writer: &writerAtWriter{
				WriterAt: w.writerAtWriter.WriterAt,
				pos:      w.pos,
			},
		},
		toWrite: n,
		obw:     w.StickyBigEndianWriter,
	}

	w.pos += int64(n) * 4

	return p
}

func (w writer) ReservePointerList(n uint32) *pointerWriter {
	pw := w.ReservePointers(n)

	w.WriteUint32(0)

	return pw
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

// Errors.
var (
	ErrTooBig = errors.New("write too big")
)
