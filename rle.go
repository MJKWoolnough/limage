package xcf

import "github.com/MJKWoolnough/byteio"

type rle struct {
	Reader     byteio.StickyReader
	repeatByte bool
	data       uint8
	count      uint16
}

func (r *rle) Read(p []byte) (int, error) {
	for i := range p {
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
		r.count--
		if r.repeatByte {
			p[i] = r.data
		} else {
			p[i] = r.Reader.ReadUint8()
		}
	}
	return len(p), nil
}
