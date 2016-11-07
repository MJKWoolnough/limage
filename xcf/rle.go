package xcf

import (
	"io"

	"github.com/MJKWoolnough/byteio"
	"github.com/MJKWoolnough/errors"
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
	if r.count != 0 && r.Reader.Err == nil {
		return n, ErrInvalidRLE
	}
	return n, r.Reader.Err
}

const minRunLength = 3

func (w *writer) WriteRLE(data []byte) {
	var (
		last         byte
		run, written int
	)
	for n, b := range data {
		if b == last {
			run++
		} else {
			if run > minRunLength {
				nm := n - run - written
				w.WriteRLEData(data[written:], run, nm, last)
				written += nm + run
			}
			run = 0
		}
		last = b
	}
	nm := len(data) - run - written
	w.WriteRLEData(data[written:], run, nm, last)
}

func (w *writer) WriteRLEData(data []byte, run, l int, last byte) {
	if nm := n - run - written; l > 0 {
		if nm < 128 {
			w.WriteUint8(255 - uint8(nm-1))
		} else {
			w.WriteUint8(128)
			w.WriteUint16(uint16(nm))
		}
		w.Write(data[:written+n-run])
	}
	if run < 128 {
		w.WriteUint8(uint8(run - 1))
	} else {
		w.WriteUint8(127)
		w.WriteUint16(uint16(run))
	}
	w.WriteUint8(last)
}

const (
	ErrInvalidRLE errors.Error = "invalid RLE data"
)
