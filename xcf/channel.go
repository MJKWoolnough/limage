package xcf

import "image"

func (d *decoder) ReadChannel() *image.Gray {
	width := d.ReadUint32()
	height := d.ReadUint32()

	d.ReadString() // name

PropertyLoop:
	for {
		typ := d.ReadUint32()
		plength := d.ReadUint32()
		switch typ {
		// general properties
		case propEnd:
			if plength != 0 {
				d.SetError(ErrInvalidProperties)
			}
			break PropertyLoop
		case propLinked:
			d.ReadBoolProperty()
		case propLockContent:
			d.ReadBoolProperty()
		case propOpacity:
			if d.ReadUint32() > 255 {
				d.SetError(ErrInvalidOpacity)
			}
		case propParasites:
			d.ReadParasites(plength)
		case propTattoo:
			d.ReadUint32()
		case propVisible:
			d.ReadBoolProperty()

			//channel properties
		case propActiveChannel:
			// active channel
		case propColor:
			d.ReadUint8() // r
			d.ReadUint8() // g
			d.ReadUint8() // b
		case propSelection:
			// selected
		case propShowMasked:
			d.ReadBoolProperty()
		default:
			d.Skip(plength)
		}

	}

	hptr := d.ReadUint32()
	d.Goto(hptr)

	im := d.ReadImage(width, height, 2)
	if im != nil {
		return im.(*image.Gray) // gray
	}
	return nil
}

func (e *encoder) WriteChannel(c *image.Gray) {
	b := c.Bounds()
	e.WriteUint32(uint32(b.Dx()))
	e.WriteUint32(uint32(b.Dy()))
	e.WriteString("")

	e.WriteUint32(0)
	e.WriteUint32(0)
}
