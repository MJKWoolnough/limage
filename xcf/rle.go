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

type rleItem struct {
	char byte
	num  uint16
}

type rlencoder struct {
	Writer *byteio.StickyWriter
	queue  []rleItem
}

func (r *rlencoder) Write(p []byte) (int, error) {
	if r.Writer.Err != nil || len(p) == 0 {
		return 0, r.Writer.Err
	}
	lastChar := p[0]
	run := uint(1)
	if l := len(r.queue); l > 0 {
		if lastChar == r.queue[l-1].char {
			run = uint(r.queue[l-1].num) + 1
			r.queue = r.queue[:l-1]
		}
	}
	for _, c := range p[1:] {
		if c == lastChar {
			run++
		} else {
			r.set(lastChar, run)
			lastChar = c
			run = 1
		}
	}
	r.set(lastChar, run)
	return len(p), r.Writer.Err
}

func (r *rlencoder) set(char byte, run uint) {
	tr := uint16(run)
	r.queue = append(r.queue, rleItem{char, tr})
	if tr > minRunLength {
		r.Flush()
	}
	run -= uint(tr)
	for run > 0 {
		r.queue = append(r.queue, rleItem{char, 0xffff})
		run >>= 16
		r.Flush()
	}
}

func (r *rlencoder) Flush() error {
	switch l := len(r.queue); l {
	case 0:
	case 1:
		r.writeRun(r.queue[0].char, r.queue[0].num)
	case 2:
		if r.queue[1].num >= minRunLength {
			r.writeRun(r.queue[0].char, r.queue[0].num)
			r.writeRun(r.queue[1].char, r.queue[1].num)
		} else {
			r.writeData(r.queue)
		}
	default:
		if r.queue[l-1].num >= minRunLength {
			r.writeData(r.queue[:l-1])
			r.writeRun(r.queue[l-1].char, r.queue[l-1].num)
		} else {
			r.writeData(r.queue)
		}
	}
	r.queue = r.queue[:0]
	return r.Writer.Err
}

func (r *rlencoder) writeRun(char byte, num uint16) {
	if num < 128 {
		r.Writer.WriteUint8(uint8(num - 1))
	} else {
		r.Writer.WriteUint8(127)
		r.Writer.WriteUint16(num)
	}
	r.Writer.WriteUint8(char)
}

func (r *rlencoder) writeData(data []rleItem) {
	l := uint(0)
	for _, i := range data {
		l += uint(i.num)
	}
	d := make([]byte, 0, l)
	for _, i := range data {
		for j := uint16(0); j < i.num; j++ {
			d = append(d, i.char)
		}
	}
	for len(d) > 0xffff {
		r.Writer.WriteUint8(128)
		r.Writer.WriteUint16(0xffff)
		r.Writer.Write(d[:0xffff])
		d = d[0xffff:]
	}
	if l := len(d); l == 0 {
		return
	} else if l < 128 {
		r.Writer.WriteUint8(uint8(256 - l))
	} else {
		r.Writer.WriteUint16(128)
		r.Writer.WriteUint16(uint16(l))
	}
	r.Writer.Write(d)
}
