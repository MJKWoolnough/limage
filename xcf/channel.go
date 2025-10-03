package xcf

import "image"

func (d *decoder) ReadChannel() *image.Gray {
	width := d.ReadUint32()
	height := d.ReadUint32()

	d.SkipString() // name
	d.skipProperties()

	var hptr uint64

	if d.mode < 2 {
		hptr = uint64(d.ReadUint32())
	} else {
		hptr = d.ReadUint64()
	}

	d.Goto(hptr)

	if im := d.ReadImage(width, height, 2); im != nil {
		return im.(*image.Gray) // gray
	}

	return nil
}

func (d *decoder) skipProperties() {
	for {
		typ := d.ReadUint32()
		plength := d.ReadUint32()

		switch typ {
		// general properties
		case propEnd:
			if plength != 0 {
				d.SetError(ErrInvalidProperties)
			}

			return
		case propOpacity:
			if d.ReadUint32() > 255 {
				d.SetError(ErrInvalidOpacity)
			}
		default:
			d.Skip(plength)
		}

	}
}

func (e *encoder) WriteChannel(c *image.Gray) {
	b := c.Bounds()

	e.WriteUint32(uint32(b.Dx()))
	e.WriteUint32(uint32(b.Dy()))
	e.WriteString("")

	e.WriteUint32(0) // No properties
	e.WriteUint32(0)

	e.WriteUint32(uint32(e.pos) + 4) // hptr
	e.WriteImage(c, (*encoder).grayToBuf, 1)
}
