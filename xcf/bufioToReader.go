package xcf

import (
	"bufio"
	"io"
	"sync"
	"unsafe"
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

type readerAt struct {
	readMutex  sync.Mutex
	ReadSeeker io.ReadSeeker
}

func (r *readerAt) ReadAt(p []byte, offset int64) (int, error) {
	r.readMutex.Lock()
	r.ReadSeeker.Seek(offset, io.SeekStart)
	n, err := r.ReadSeeker.Read(p)
	r.readMutex.Unlock()
	return n, err
}

func bufioToReader(b *bufio.Reader) io.ReaderAt {
	br := (*bufioReader)(unsafe.Pointer(b))
	if ra, ok := br.rd.(io.ReaderAt); ok {
		return ra
	} else if rs, ok := br.rd.(io.ReadSeeker); ok {
		rs.Seek(0, 0)
		return &readerAt{
			ReadSeeker: rs,
		}
	}
	return nil
}
