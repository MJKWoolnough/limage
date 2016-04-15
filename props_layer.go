package xcf

type applyMask bool

func (d *Decoder) readApplyMask() applyMask {
	switch d.r.ReadUint32() {
	case 0:
		return false
	case 1:
		return true
	}
	d.r.Err = ErrInvalidState
	return false
}

type editMask bool

func (d *Decoder) readEditMask() editMask {
	switch d.r.ReadUint32() {
	case 0:
		return false
	case 1:
		return true
	}
	d.r.Err = ErrInvalidState
	return false
}

func (d *Decoder) readFloatingSelection() uint32 {
	return d.r.ReadUint32()
}

func (d *Decoder) readItemPath(length uint32) []uint32 {
	pts := length >> 2
	pointers := make([]uint32, pts)
	for i := uint32(0); i < pts; i++ {
		pointers[i] = d.r.ReadUint32()
	}
	return pointers
}

func (d *Decoder) readGroupItemFlags() uint32 {
	return d.r.ReadUint32() | 1
}

type lockAlpha bool

func (d *Decoder) readLockAlpha() lockAlpha {
	switch d.r.ReadUint32() {
	case 0:
		return false
	case 1:
		return true
	}
	d.r.Err = ErrInvalidState
	return false
}

type mode uint8

func (d *Decoder) readMode() mode {
	m := d.r.ReadUint32()
	if m > 21 {
		d.r.Err = ErrInvalidState
		return 0
	}
	return mode(m)
}

type offsets struct {
	x, y int32
}

func (d *Decoder) readOffsets() offsets {
	return offsets{
		d.r.ReadInt32(),
		d.r.ReadInt32(),
	}
}

type showMask bool

func (d *Decoder) readShowMask() showMask {
	switch d.r.ReadUint32() {
	case 0:
		return false
	case 1:
		return true
	}
	d.r.Err = ErrInvalidState
	return false
}

type textLayerFlags uint8

func (d *Decoder) readTextLayerFlags() textLayerFlags {
	t := d.r.ReadUint32()
	if t > 3 {
		d.r.Err = ErrInvalidState
		return 0
	}
	return textLayerFlags(t)
}
