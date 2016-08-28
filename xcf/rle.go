package xcf

import "github.com/MJKWoolnough/byteio"

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
			n := r.Reader.ReadUint8()
			if n < 127 {
				r.repeatByte = true
				r.count = uint16(n) + 1
				r.data = r.Reader.ReadUint8()
			} else if n == 127 {
				r.repeatByte = true
				r.count = r.Reader.ReadUint16()
				r.data = r.Reader.ReadUint8()
			} else if n == 128 {
				r.repeatByte = false
				r.count = r.Reader.ReadUint16()
			} else {
				r.repeatByte = false
				r.count = 256 - uint16(n)
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

const minRunLength = 1

func runLengthEncode(w byteio.StickyWriter, pix []byte) {
	var (
		lastByte      = pix[0]
		run      uint = 0
	)
	for i := 1; i < len(pix); i++ {
		if pix[i] == lastByte {
			run++
		} else {
			if run >= minRunLength {
				doWrites(w, lastByte, run, pix[:i-1])
			}
			pix = pix[:i]
			run = 1
			lastByte = pix[0]
		}
	}
	doWrites(w, lastByte, run, pix)
}

func doWrites(w byteio.StickyWriter, lastByte byte, run uint, pix []byte) {
	if l := uint(len(pix)); l-run > 0 {
		writeData(w, pix[:l-run])
	}
	if run > 0 {
		writeRun(w, lastByte, run)
	}
}

func writeRun(w byteio.StickyWriter, b byte, run uint) {
	if run <= 127 {
		w.WriteUint8(uint8(run) - 1)
	} else {
		w.WriteUint8(127)
		w.WriteUint16(uint16(run))
	}
	w.WriteUint8(b)
}

func writeData(w byteio.StickyWriter, pix []byte) {
	if len(pix) <= 127 {
		w.WriteUint8(uint8(256 - len(pix)))
	} else {
		w.WriteUint8(128)
		w.WriteUint16(uint16(len(pix)))
	}
	w.Write(pix)
}
