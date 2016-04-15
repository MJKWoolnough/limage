package xcf

import "image/color"

func (d *Decoder) readColor() color.Color {
	return color.RGBA{
		d.r.ReadUint8(),
		d.r.ReadUint8(),
		d.r.ReadUint8(),
	}
}

type showMasked bool

func (d *Decoder) readShowMasked() showMasked {
	switch d.r.ReadUint32() {
	case 0:
		return false
	case 1:
		return true
	}
	d.r.Err = ErrInvalidState
	return false
}
