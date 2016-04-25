package xcf

import "github.com/MJKWoolnough/byteio"

type rle struct {
	r     byteio.StickyReader
	mode  byte
	data  byte
	count uint16
}

func (r *rle) ReadByte() byte {
	if r.count == 0 {
		m := r.r.ReadUint8()
		if m < 127 {
			r.mode = 0
			r.count = uint16(m) + 1
			r.data = r.r.ReadUint8()
		} else if m == 127 {
			r.mode = 0
			r.count = r.r.ReadUint16()
			r.data = r.r.ReadUint8()
		} else if m == 128 {
			r.mode = 1
			r.count = r.r.ReadUint16()
		} else {
			r.mode = 1
			r.count = 256 - uint16(m)
		}
	}
	r.count--
	if r.mode == 0 {
		return r.data
	}
	return r.r.ReadUint8()
}
