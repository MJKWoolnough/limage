package xcf

import "errors"

type parasite struct {
	name  string
	flags uint32
	data  []byte
}

type parasites []parasite

func (p parasites) Get(name string) *parasite {
	for n := range p {
		if p[n].name == name {
			return &p[n]
		}
	}
	return nil
}

func (d *decoder) ReadParasites(l uint32) parasites {
	ps := make(parasites, 0, 32)
	for l > 0 {
		var p parasite
		p.name = d.ReadString()
		p.flags = d.ReadUint32()
		pplength := d.ReadUint32()
		read := 4 + uint32(len(p.name)) + 1 // length (uint32) + string([]byte) + \0 (byte)
		read += 4                           // flags
		read += 4                           // pplength
		read += pplength                    // len(data)
		if read > l {
			d.SetError(ErrInvalidParasites)
			return nil
		}
		l -= read
		p.data = make([]byte, pplength)
		d.Read(p.data)
		ps = append(ps, p)
	}
	return ps
}

func (d *decoder) ReadParasite() parasite {
	var p parasite
	p.name = d.ReadString()
	p.flags = d.ReadUint32()
	pplength := d.ReadUint32()
	p.data = make([]byte, pplength)
	d.Read(p.data)
	return p
}

// Errors
var (
	ErrInvalidParasites = errors.New("invalid parasites layout")
)
