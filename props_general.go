package xcf

import "image/color"

func (d *Decoder) readOpacity() color.Alpha {
	a := d.r.ReadUint32()
	if a > 255 {
		d.r.Err = ErrInvalidState
	}
	return color.Alpha{uint8(a)}
}

type parasite struct {
	name  string
	flags uint32
	data  []byte
}

func (d *Decoder) readParasites(length uint32) []parasite {
	o := d.r.Count
	ps := make([]parasite, 0)
	for length > 0 {
		d.r.Count = 0
		name := d.r.ReadString()
		flags := d.r.ReadUint32()
		pplength := d.r.ReadUint32()
		data := make([]byte, pplength)
		d.r.Read(data)
		ps = append(ps, parasite{
			name:  name,
			flags: flags,
			data:  data,
		})
		length -= uint32(d.r.Count)
		o += d.r.Count
	}
	d.r.Count = o
	return ps
}

type tattoo uint32

func (d *Decoder) readTattoo() tattoo {
	t := d.r.ReadUint32()
	if t == 0 {
		if d.r.Err == nil {
			d.r.Err = ErrInvalidState
		}
		return 0
	}
	return tattoo(t)
}
