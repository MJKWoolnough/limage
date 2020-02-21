package xcf

import (
	"errors"
	"io"

	"vimagination.zapto.org/byteio"
	"vimagination.zapto.org/memio"
)

type rle struct {
	Reader     *byteio.StickyBigEndianReader
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
	if len(data) == 0 {
		return
	}
	var run, written int
	last := data[0]
	for n, b := range data {
		if b == last {
			run++
		} else {
			if run > minRunLength {
				w.WriteRLEData(data[written:n-run], run, last)
				written = n
			}
			run = 1
		}
		last = b
	}
	if run <= minRunLength && run < len(data) {
		run = 0
	}
	w.WriteRLEData(data[written:len(data)-run], run, last)
}

func (w *writer) WriteRLEData(data []byte, run int, last byte) {
	if len(data) > 0 {
		r := false
		if len(data) <= minRunLength {
			r = true
			for _, b := range data {
				if b != data[0] {
					r = false
					break
				}
			}
		}
		if r {
			w.WriteRLEData(nil, len(data), data[0])
		} else {
			if len(data) < 128 {
				w.WriteUint8(255 - uint8(len(data)-1))
			} else {
				w.WriteUint8(128)
				w.WriteUint16(uint16(len(data)))
			}
			w.Write(data)
		}
	}
	if run > 0 {
		if run < 128 {
			w.WriteUint8(uint8(run - 1))
		} else {
			w.WriteUint8(127)
			w.WriteUint16(uint16(run))
		}
		w.WriteUint8(last)
	}
}

func (d *decoder) readRLE(count int, buf *memio.Buffer) error {
	bw := byteio.BigEndianWriter{Writer: buf}
	for count > 0 {
		c := d.reader.ReadUint8()
		bw.WriteUint8(c)
		if c < 127 {
			count -= int(c) + 1
			bw.WriteUint8(d.ReadUint8())
		} else if c == 127 {
			p := d.ReadUint16()
			bw.WriteUint16(p)
			count -= int(p)
			bw.WriteUint8(d.ReadUint8())
		} else if c == 128 {
			p := d.ReadUint16()
			bw.WriteUint16(p)
			count -= int(p)
			io.CopyN(buf, d, int64(p))
		} else {
			count -= 256 - int(c)
			io.CopyN(buf, d, 256-int64(c))
		}
	}
	if count < 0 {
		return ErrInvalidRLE
	}
	return nil
}

// Errors
var (
	ErrInvalidRLE = errors.New("invalid RLE data")
)
