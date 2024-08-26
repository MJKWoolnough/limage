package xcf

import "errors"

type paths struct {
	aIndex uint32
	paths  []path
}

type path struct {
	name   string
	linked bool
	// state   byte
	closed bool
	// version uint32
	tattoo uint32
	points []pathPoint
}

type pathPoint struct {
	control bool
	x, y    float64
}

func (d *reader) ReadPaths() paths {
	var p paths

	p.aIndex = d.ReadUint32()

	n := d.ReadUint32()
	p.paths = make([]path, n)

	for i := uint32(0); i < n; i++ {
		p.paths[i].name = d.ReadString()
		p.paths[i].linked = d.ReadBoolProperty()
		state := d.ReadUint8()
		p.paths[i].closed = d.ReadBoolProperty()

		if p.paths[i].closed {
			if state != 4 {
				d.SetError(ErrInconsistantClosedState)

				return p
			}
		} else {
			if state != 2 {
				d.SetError(ErrInconsistantClosedState)

				return p
			}
		}

		np := d.ReadUint32()
		v := d.ReadUint32()

		if v < 1 || v > 3 {
			d.SetError(ErrUnknownPathsVersion)

			return p
		}

		if v == 2 || v == 3 {
			d.SkipUint32()
		}

		if v == 3 {
			p.paths[i].tattoo = d.ReadUint32()
		}

		p.paths[i].points = make([]pathPoint, np)

		for j := uint32(0); j < np; j++ {
			p.paths[i].points[j].control = d.ReadBoolProperty()

			if v == 1 {
				p.paths[i].points[j].x = float64(d.ReadInt32())
				p.paths[i].points[j].y = float64(d.ReadInt32())
			} else {
				p.paths[i].points[j].x = float64(d.ReadFloat32())
				p.paths[i].points[j].y = float64(d.ReadFloat32())
			}
		}
	}

	return p
}

// Errors.
var (
	ErrInconsistantClosedState = errors.New("inconsistent closed state")
	ErrUnknownPathsVersion     = errors.New("unknown paths version")
)
