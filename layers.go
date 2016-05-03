package xcf

import "os"

type Layer interface {
	IsGroup() bool
	AsGroup() *LayerGroup
	AsImage() *LayerImage
}

type LayerGroup struct {
	layer
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
	layer
	alpha bool
}

func (LayerImage) IsGroup() bool {
	return false
}

func (LayerImage) AsGroup() *LayerGroup {
	return nil
}

func (l *LayerImage) AsImage() *LayerImage {
	return l
}

type layer struct {
	OffsetX, OffsetY                            int32
	Width, Height                               uint32
	Name                                        string
	editMask, showMask, visible, locked, active bool
}

func (d *Decoder) readLayer() Layer {
	var l layer
	l.Width = d.r.ReadUint32()
	l.Height = d.r.ReadUint32()
	typ := d.r.ReadUint32()
	l.Name = d.r.ReadString()

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
			l.OffsetX = d.r.ReadInt32()
			l.OffsetY = d.r.ReadInt32()
		case propShowMask:
			l.showMask = d.readBool()
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
	_, _ = hptr, mptr
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
		return nil
	}
	return &LayerImage{
		layer: l,
	}
}
