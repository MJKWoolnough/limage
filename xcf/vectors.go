package xcf

import "vimagination.zapto.org/errors"

type vectors struct {
	aIndex uint32
	paths  []vpath
}

type vpath struct {
	name            string
	tattoo          uint32
	visible, linked bool
	parasites
	strokes []stroke
}

type stroke struct {
	closed bool
	points []point
}

type point struct {
	control                             bool
	x, y, pressure, xtilt, ytilt, wheel float32
}

func (d *reader) ReadVectors() vectors {
	v := d.ReadUint32()
	if v != 1 {
		d.SetError(ErrUnknownVectorVersion)
		return vectors{}
	}
	var vs vectors
	vs.aIndex = d.ReadUint32()
	n := d.ReadUint32()
	vs.paths = make([]vpath, n)
	for i := uint32(0); i < n; i++ {
		vs.paths[i].name = d.ReadString()
		vs.paths[i].tattoo = d.ReadUint32()
		vs.paths[i].visible = d.ReadBoolProperty()
		vs.paths[i].linked = d.ReadBoolProperty()
		m := d.ReadUint32()
		k := d.ReadUint32()
		vs.paths[n].parasites = make(parasites, m)
		vs.paths[n].strokes = make([]stroke, k)
		for j := uint32(0); j < m; j++ {
			vs.paths[i].parasites[j] = d.ReadParasite()
		}
		for j := uint32(0); j < k; j++ {
			b := d.ReadUint32()
			if b != 1 {
				d.SetError(ErrUnknownStrokeType)
				return vs
			}
			vs.paths[i].strokes[j].closed = d.ReadBoolProperty()
			nf := d.ReadUint32()
			if nf < 2 || nf > 6 {
				d.SetError(ErrInvalidFloatsNumber)
				return vs
			}
			np := d.ReadUint32()
			vs.paths[i].strokes[j].points = make([]point, np)
			for ii := uint32(0); ii < np; ii++ {
				vs.paths[i].strokes[j].points[ii].control = d.ReadBoolProperty()
				vs.paths[i].strokes[j].points[ii].x = d.ReadFloat32()
				vs.paths[i].strokes[j].points[ii].y = d.ReadFloat32()
				vs.paths[i].strokes[j].points[ii].pressure = 1
				vs.paths[i].strokes[j].points[ii].xtilt = 0.5
				vs.paths[i].strokes[j].points[ii].ytilt = 0.5
				vs.paths[i].strokes[j].points[ii].wheel = 0.5
				if nf >= 3 {
					vs.paths[i].strokes[j].points[ii].pressure = d.ReadFloat32()
					if nf >= 4 {
						vs.paths[i].strokes[j].points[ii].xtilt = d.ReadFloat32()
						if nf >= 5 {
							vs.paths[i].strokes[j].points[ii].ytilt = d.ReadFloat32()
							if nf == 6 {
								vs.paths[i].strokes[j].points[ii].wheel = d.ReadFloat32()
							}
						}
					}
				}
			}
		}
	}
	return vs
}

// Errors
const (
	ErrUnknownVectorVersion errors.Error = "unknown vector version"
	ErrUnknownStrokeType    errors.Error = "unknown stroke type"
	ErrInvalidFloatsNumber  errors.Error = "invalids number of floats"
)
