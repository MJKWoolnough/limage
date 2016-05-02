package xcf

import (
	"errors"
	"image/color"
)

type property uint8

const (
	propEnd               property = 0
	propColormap          property = 1
	propActiveLayer       property = 2
	propActiveChannel     property = 3
	propSelection         property = 4
	propFloatingSelection property = 5
	propOpacity           property = 6
	propMode              property = 7
	propVisible           property = 8
	propLinked            property = 9
	propLockAlpha         property = 10
	propApplyMask         property = 11
	propEditMask          property = 12
	propShowMask          property = 13
	propShowMasked        property = 14
	propOffsets           property = 15
	propColor             property = 16
	propCompression       property = 17
	propGuides            property = 18
	propResolution        property = 19
	propTattoo            property = 20
	propParasites         property = 21
	propUnit              property = 22
	propPaths             property = 23
	propUserUnit          property = 24
	propVectors           property = 25
	propTextLayerFlags    property = 26
	propSamplePoints      property = 27
	propLockContent       property = 28
	propGroupItem         property = 29
	propItemPath          property = 30
	propGroupItemFlags    property = 31
	propLockPosition      property = 32
	propFloatOpacity      property = 33
)

func (d *Decoder) readBool() bool {
	switch d.r.ReadUint32() {
	case 0:
		return false
	case 1:
		return true
	}
	d.r.Err = ErrInvalidState
	return false
}

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
	l := int64(length)
	for l > 0 && d.r.Err == nil {
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
		l -= d.r.Count
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

func (d *Decoder) readColorMap() color.Palette {
	num := d.r.ReadUint32()
	c := make(color.Palette, num)
	for i := uint32(0); i < num; i++ {
		c[i] = color.RGBA{
			d.r.ReadUint8(),
			d.r.ReadUint8(),
			d.r.ReadUint8(),
			255,
		}
	}
	return c
}

type compression uint8

const (
	compNone    compression = 0
	compRLE     compression = 1
	compZlib    compression = 2
	compFractal compression = 3
)

func (d *Decoder) readCompression() compression {
	switch c := d.r.ReadUint8(); c {
	case 0, 1, 2, 3:
		return compression(c)
	}
	d.r.Err = ErrInvalidState
	return 0
}

type orientation bool

func (o orientation) IsHorizontal() bool {
	return bool(!o)
}

func (o orientation) IsVertical() bool {
	return bool(o)
}

func (o orientation) String() string {
	if o {
		return "horizontal"
	}
	return "vertical"
}

type guide struct {
	coord       int32
	orientation orientation
}

func (d *Decoder) readGuides(size uint32) []guide {
	num := size / 5
	g := make([]guide, num)
	for i := uint32(0); i < num; i++ {
		coord := d.r.ReadInt32()
		o := d.r.ReadUint8()
		if o != 1 && o != 2 {
			if d.r.Err == nil {
				d.r.Err = ErrInvalidState
			}
			return g
		}
		g[num] = guide{
			coord:       coord,
			orientation: o == 2,
		}
	}
	return g
}

type unit uint8

const (
	unitInches unit = 1
	unitMM     unit = 2
	unitPoints unit = 3
	unitPicas  unit = 4
)

func (d *Decoder) readUnit() unit {
	u := unit(d.r.ReadUint32())
	switch u {
	case unitInches, unitMM, unitPoints, unitPicas:
		return u
	}
	if d.r.Err == nil {
		d.r.Err = ErrInvalidState
	}
	return 0
}

type point struct {
	anchorControl bool
	x, y          float64
}

type path struct {
	name   string
	linked bool
	closed bool
	tattoo uint32
	points []point
}

type paths struct {
	index uint32
	paths []path
}

func (d *Decoder) readPaths() paths {
	index := d.r.ReadUint32()
	num := d.r.ReadUint32()
	p := paths{
		index,
		make([]path, num),
	}
	for i := uint32(0); i < num; i++ {
		pt := &p.paths[i]
		pt.name = d.r.ReadString()
		linked := d.r.ReadUint32()
		switch linked {
		case 0:
			pt.linked = false
		case 1:
			pt.linked = true
		default:
			d.r.Err = ErrInvalidState
			return p
		}
		state := d.r.ReadUint8()
		switch d.r.ReadUint32() {
		case 0:
			if state != 4 {
				d.r.Err = ErrInvalidState
				return p
			}
			pt.closed = false
		case 1:
			if state != 2 {
				d.r.Err = ErrInvalidState
				return p
			}
			pt.closed = true
		default:
			d.r.Err = ErrInvalidState
			return p
		}
		np := d.r.ReadUint32()
		version := d.r.ReadUint32()
		if version == 2 || version == 3 {
			d.r.ReadUint32()
			if version == 3 {
				pt.tattoo = d.r.ReadUint32()
			}
		} else if version != 1 {
			if d.r.Err == nil {
				d.r.Err = ErrInvalidState
			}
			return p
		}
		pt.points = make([]point, np)
		for j := uint32(0); j < np; j++ {
			switch d.r.ReadInt32() {
			case 0:
				pt.points[j].anchorControl = false
			case 1:
				pt.points[j].anchorControl = true
			default:
				d.r.Err = ErrInvalidState
				return p
			}
			if version == 1 {
				pt.points[j].x = float64(d.r.ReadInt32())
				pt.points[j].y = float64(d.r.ReadInt32())
			} else {
				pt.points[j].x = float64(d.r.ReadFloat32())
				pt.points[j].y = float64(d.r.ReadFloat32())
			}
		}
	}
	return p
}

type userUnit struct {
	factor                           float32
	digits                           uint32
	id, symbol, abbrev, sname, pname string
}

func (d *Decoder) readUserUnit() userUnit {
	var u userUnit
	u.factor = d.r.ReadFloat32()
	u.digits = d.r.ReadUint32()
	u.id = d.r.ReadString()
	u.symbol = d.r.ReadString()
	u.abbrev = d.r.ReadString()
	u.sname = d.r.ReadString()
	u.pname = d.r.ReadString()
	return u
}

type vectorpoint struct {
	anchorControl                       bool
	x, y, pressure, xtilt, ytilt, wheel float32
}

type stroke struct {
	closed bool
	points []vectorpoint
}

type vectorpath struct {
	name            string
	tattoo          uint32
	visible, linked bool
	parasites       [][]parasite
	strokes         []stroke
}

type vectors struct {
	index uint32
	paths []vectorpath
}

func (d *Decoder) readVectors() vectors {
	var v vectors
	if d.r.ReadUint32() != 1 {
		if d.r.Err == nil {
			d.r.Err = ErrInvalidState
		}
		return v
	}
	v.index = d.r.ReadUint32()
	n := d.r.ReadUint32()
	v.paths = make([]vectorpath, n)
	for i := uint32(0); i < i; i++ {
		vp := &v.paths[i]
		vp.name = d.r.ReadString()
		vp.tattoo = d.r.ReadUint32()
		switch d.r.ReadUint32() {
		case 0:
			vp.visible = false
		case 1:
			vp.visible = true
		default:
			d.r.Err = ErrInvalidState
			return v
		}
		switch d.r.ReadUint32() {
		case 0:
			vp.linked = false
		case 1:
			vp.linked = true
		default:
			d.r.Err = ErrInvalidState
			return v
		}
		m := d.r.ReadUint32()
		k := d.r.ReadUint32()
		vp.parasites = make([][]parasite, m)
		for j := uint32(0); j < m; j++ {
			// TODO: the following could be incorrect, needs checking
			if property(d.r.ReadUint32()) != propParasites {
				if d.r.Err == nil {
					d.r.Err = ErrInvalidState
				}
			}
			vp.parasites[j] = d.readParasites(d.r.ReadUint32())
		}
		vp.strokes = make([]stroke, k)
		for j := uint32(0); j < k; j++ {
			if d.r.ReadUint32() != 1 { // should be 1, bezier curve
				if d.r.Err == nil {
					d.r.Err = ErrInvalidState
				}
				return v
			}
			switch d.r.ReadUint32() {
			case 0:
				vp.strokes[j].closed = false
			case 1:
				vp.strokes[j].closed = true
			default:
				d.r.Err = ErrInvalidState
				return v
			}
			nf := d.r.ReadUint32()
			if nf < 2 || nf > 6 {
				d.r.Err = ErrInvalidState
				return v
			}
			np := d.r.ReadUint32()
			vp.strokes[j].points = make([]vectorpoint, np)
			for k := uint32(0); k < np; k++ {
				p := &vp.strokes[j].points[k]
				switch d.r.ReadUint32() {
				case 0:
					p.anchorControl = false
				case 1:
					p.anchorControl = true
				default:
					d.r.Err = ErrInvalidState
					return v
				}
				p.x = d.r.ReadFloat32()
				p.y = d.r.ReadFloat32()
				if nf > 3 {
					p.pressure = d.r.ReadFloat32()
					if nf > 4 {
						p.xtilt = d.r.ReadFloat32()
						if nf > 5 {
							p.ytilt = d.r.ReadFloat32()
							if nf == 6 {
								p.wheel = d.r.ReadFloat32()
							}
						}
					}
				}
			}
		}
	}
	return v
}

func (d *Decoder) readItemPath(length uint32) []uint32 {
	pts := length >> 2
	pointers := make([]uint32, pts)
	for i := uint32(0); i < pts; i++ {
		pointers[i] = d.r.ReadUint32()
	}
	return pointers
}

func (d *Decoder) readMode() uint8 {
	m := d.r.ReadUint32()
	if m > 21 {
		d.r.Err = ErrInvalidState
		return 0
	}
	return uint8(m)
}

func (d *Decoder) readTextLayerFlags() uint8 {
	t := d.r.ReadUint32()
	if t > 3 {
		d.r.Err = ErrInvalidState
		return 0
	}
	return uint8(t)
}

// Errors
var (
	ErrInvalidState = errors.New("invalid state")
)
