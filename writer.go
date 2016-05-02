package xcf

import (
	"io"

	"github.com/MJKWoolnough/byteio"
)

type writer struct {
	byteio.StickyWriter
}

func newWriter(w io.Writer) *writer {
	var write writer
	write.Writer = byteio.BigEndianWriter{w}
	return &write
}

func (w writer) WriteString(s string) {
	b := make([]byte, len(s)+1)
	copy(b, s)
	w.WriteUint32(uint32(len(b)))
	w.Write(b)
}
