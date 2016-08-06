package xcf

import (
	"bufio"
	"io"
	"unsafe"

	"github.com/MJKWoolnough/errors"
)

// copied from bufio
type bufioReader struct {
	buf          []byte
	rd           io.Reader
	r, w         int
	err          error
	lastByte     int
	lastRuneSize int
}

type readSeeker struct {
	Buffer     []byte
	ReadSeeker io.ReadSeeker
	Pos        int64
}

func (r *readSeeker) Read(p []byte) (n int, err error) {
	if r.Pos < int64(len(r.Buffer)) {
		n += copy(p, r.Buffer[r.Pos:])
		p = p[n:]
		r.Pos += int64(n)
	}
	if len(p) > 0 {
		m, err := r.ReadSeeker.Read(p)
		r.Pos += int64(m)
		return n + m, err
	}
	return n, nil
}

func (r *readSeeker) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0:
		r.Pos = offset
	case 1:
		if l := int64(len(r.Buffer)); r.Pos >= l {
			n, err := r.ReadSeeker.Seek(offset, 1)
			r.Pos = n + l
			return r.Pos, err
		}
		r.Pos += offset
	case 2:
		// should never be used
		return 0, errors.Error("unimplemented")
	}
	var err error
	if l := int64(len(r.Buffer)); r.Pos > l {
		_, err = r.Seek(r.Pos-l, 0)
	} else {
		_, err = r.Seek(0, 0)
	}
	return r.Pos, err
}

func bufioToReader(b *bufio.Reader) io.Reader {
	br := (*bufioReader)(unsafe.Pointer(b))
	if rs, ok := br.rd.(io.ReadSeeker); ok {
		rs.Seek(0, 0)
		return &readSeeker{
			Buffer:     br.buf,
			ReadSeeker: rs,
		}
	}
	return b
}
