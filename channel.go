package xcf

import (
	"image"
	"os"
)

type channel struct {
	width, height                uint32
	name                         string
	linked, lockContent, visible bool
	opacity                      uint8
	parasites                    parasites
	tattoo                       uint32
	active, selection, show      bool
	color                        rgb
	image                        image.Image
}

func (d *decoder) ReadChannel() channel {
	var c channel

	c.width = d.ReadUint32()
	c.height = d.ReadUint32()

	c.name = d.ReadString()

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
			c.linked = d.ReadBoolProperty()
		case propLockContent:
			c.lockContent = d.ReadBoolProperty()
		case propOpacity:
			o := d.ReadUint32()
			if o > 255 {
				d.SetError(ErrInvalidOpacity)
			}
			c.opacity = uint8(o)
		case propParasites:
			c.parasites = d.ReadParasites(plength)
		case propTattoo:
			c.tattoo = d.ReadUint32()
		case propVisible:
			c.visible = d.ReadBoolProperty()

			//channel properties
		case propActiveChannel:
			c.active = true
		case propColor:
			c.color.R = d.ReadUint8()
			c.color.G = d.ReadUint8()
			c.color.B = d.ReadUint8()
		case propSelection:
			c.selection = true
		case propShowMasked:
			c.show = d.ReadBoolProperty()
		default:
			d.Seek(int64(plength), os.SEEK_CUR)
		}

	}

	hptr := d.ReadUint32()
	d.Seek(int64(hptr), os.SEEK_SET)

	c.image = d.ReadImage(c.width, c.height, 2) // gray
	return c
}
