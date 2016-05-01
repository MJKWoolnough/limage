package xcf

import "os"

type Layer interface {
	IsGroup() bool
	AsGroup() *LayerGroup
	AsImage() *LayerImage
}

type LayerGroup struct {
}

func (LayerGroup) IsGroup() bool {
	return true
}

func (l *LayerGroup) AsGroup() *LayerGroup {
	return l
}

func (LayerGroup) AsImage() *LayerImage {
	return nil
}

type LayerImage struct {
}

func (LayerImage) IsGroup() bool {
	return true
}

func (LayerImage) AsGroup() *LayerGroup {
	return nil
}

func (l *LayerImage) AsImage() *LayerImage {
	return l
}

type layer struct {
	offsetX, offsetY                            int32
	width, height                               uint32
	name                                        string
	alpha                                       bool
	editMask, showMask, visible, locked, active bool
}

func (d *Decoder) readLayer() layer {
	var l layer
	l.width = d.r.ReadUint32()
	l.height = d.r.ReadUint32()
	typ := d.r.ReadUint32()
	l.name = d.r.ReadString()

Props:
	for {
		propID := property(d.r.ReadUint32())
		propLength := d.r.ReadUint32()
		switch propID {
		case propEnd:
			break Props
		case propActiveLayer:
			l.active = true
		case propFloatingSelection:
			f := d.r.ReadUint32()
			_ = f
		case propOpacity:
			o := d.readOpacity()
			_ = o
		case propApplyMask:
			a := d.readBool()
			_ = a
		case propEditMask:
			l.editMask = d.readBool()
			_ = e
		case propMode:
			m := d.readMode()
			_ = m
		case propLinked:
			l := d.readBool()
			_ = l
		case propLockAlpha:
			l := d.readBool()
			_ = l
		case propOffsets:
			l.offsetX = d.r.ReadInt32()
			l.offsetY = d.r.ReadInt32()
		case propShowMask:
			l.showMask = d.readBool()
			_ = s
		case propTattoo:
			t := d.readTattoo()
			_ = t
		case propParasites:
			p := d.readParasites(propLength)
			_ = p
		case propTextLayerFlags:
			t := d.readTextLayerFlags()
			_ = t
		case propLockContent:
			l.locked = d.readBool()
		case propVisible:
			l.visible = d.readBool()
		case propGroupItem:
			// g := d.readGroupItem()
			// no data, just set as item group
		case propItemPath:
			i := d.readItemPath(propLength)
			_ = i
		case propGroupItemFlags:
			g := d.r.ReadUint32()
			_ = g
		default:
			d.s.Seek(int64(propLength), os.SEEK_CUR)
		}
	}

	hptr := d.r.ReadUint32()
	mptr := d.r.ReadUint32()
	switch typ {
	case 0:
		//RGB
	case 1:
		//RGBA
	case 2:
		//Y
	case 3:
		//YA
	case 4:
		//I
	case 5:
		//IA
	default:
		d.r.Err = ErrInvalidState
		return
	}
}
