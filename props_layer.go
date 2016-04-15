package xcf

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
