package xcf

import (
	"io"

	"github.com/MJKWoolnough/byteio"
)

type rle struct {
	Reader     byteio.StickyReader
	repeatByte bool
	data       uint8
	count      uint16
}

func (r *rle) Read(p []byte) (int, error) {
	var n int
	for len(p) > 0 && r.Reader.Err == nil {
		if r.count == 0 {
			m := r.Reader.ReadUint8()
			if r.Reader.Err != nil {
				return n, r.Reader.Err
			}
			if m < 127 {
				r.repeatByte = true
				r.count = uint16(m) + 1
				r.data = r.Reader.ReadUint8()
			} else if m == 127 {
				r.repeatByte = true
				r.count = r.Reader.ReadUint16()
				r.data = r.Reader.ReadUint8()
			} else if m == 128 {
				r.repeatByte = false
				r.count = r.Reader.ReadUint16()
			} else {
				r.repeatByte = false
				r.count = 256 - uint16(m)
			}
			if r.Reader.Err != nil {
				if r.Reader.Err == io.EOF {
					r.Reader.Err = io.ErrUnexpectedEOF
				}
				return n, r.Reader.Err
			}
		}
		c := int(r.count)
		if len(p) < c {
			c = len(p)
		}
		if r.repeatByte {
			for i := 0; i < c; i++ {
				p[i] = r.data
			}
		} else {
			r.Reader.Read(p[:c])
		}
		r.count -= uint16(c)
		p = p[c:]
		n += c
	}
	return n, r.Reader.Err
}

const minRunLength = 3

func runLengthEncode(w *byteio.StickyWriter, pix []byte) {
	if len(pix) == 0 {
		return
	}
	var (
		lastByte      = pix[0]
		run      uint = 1
	)
	for i := 1; i < len(pix); i++ {
		if pix[i] == lastByte {
			run++
		} else {
			if run >= minRunLength {
				doWrites(w, lastByte, run, pix[:i-1])
				pix = pix[i:]
				i = 1
			}
			run = 1
			lastByte = pix[i]
		}
	}
	if run < minRunLength && uint(len(pix)) != run {
		run = 0
	}
	doWrites(w, lastByte, run, pix)
}

func doWrites(w *byteio.StickyWriter, lastByte byte, run uint, pix []byte) {
	if l := uint(len(pix)); l > run {
		writeData(w, pix[:l-run])
	}
	if run > 0 {
		writeRun(w, lastByte, run)
	}
}

func writeRun(w *byteio.StickyWriter, b byte, run uint) {
	if run <= 127 {
		w.WriteUint8(uint8(run) - 1)
	} else {
		w.WriteUint8(127)
		w.WriteUint16(uint16(run))
	}
	w.WriteUint8(b)
}

func writeData(w *byteio.StickyWriter, pix []byte) {
	if len(pix) <= 127 {
		w.WriteUint8(uint8(256 - len(pix)))
	} else {
		w.WriteUint8(128)
		w.WriteUint16(uint16(len(pix)))
	}
	w.Write(pix)
}
