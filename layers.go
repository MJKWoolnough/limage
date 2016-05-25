package xcf

import (
	"image"
	"image/color"
	"os"
)

type Layer interface {
	IsGroup() bool
	IsImage() bool
	IsText() bool
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

type LayerImage struct {
	layer
	alpha bool
	image.Image
}

func (LayerImage) IsImage() bool {
	return true
}

func (l *LayerImage) AsImage() *LayerImage {
	return l
}

type LayerText struct {
	layer
	TextData TextData
}

func (l *LayerText) AsText() *LayerText {
	return l
}

type layer struct {
	OffsetX, OffsetY                            int32
	Width, Height                               uint32
	Name                                        string
	Mode                                        uint8
	Opacity                                     color.Alpha
	editMask, showMask, visible, locked, active bool
}

func (layer) IsGroup() bool {
	return false
}

func (layer) IsImage() bool {
	return false
}

func (layer) IsText() bool {
	return false
}

func (layer) AsGroup() *LayerGroup {
	return nil
}

func (layer) AsImage() *LayerImage {
	return nil
}

func (layer) AsText() *LayerText {
	return nil
}

func (d *Decoder) readLayer() Layer {
	var (
		l         layer
		group     bool
		text      bool
		parasites []parasite
	)
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
			l.Opacity = d.readOpacity()
		case propApplyMask:
			a := d.readBool()
			_ = a
		case propEditMask:
			l.editMask = d.readBool()
		case propMode:
			l.Mode = d.readMode()
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
			parasites = d.readParasites(propLength)
		case propTextLayerFlags:
			t := d.readTextLayerFlags()
			_ = t
			text = true
		case propLockContent:
			l.locked = d.readBool()
		case propVisible:
			l.visible = d.readBool()
		case propGroupItem:
			group = true
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
	if group {
		return &LayerGroup{
			layer: l,
		}
	} else if text {
		var (
			t   TextData
			err error
		)
		for _, p := range parasites {
			if p.name == "gimp-text-layer" {
				t, err = parseTextParasite(p.data)
				if err != nil {
					d.r.Err = err
				}
				break
			}
		}
		return &LayerText{
			layer:    l,
			TextData: t,
		}
	}
	r := image.Rect(0, 0, int(l.Width), int(l.Height))
	alpha := typ&1 == 1
	var im image.Image
	if d.props.baseType != typ>>1 {
		//incorrect type
	}
	switch typ >> 1 {
	case 0:
		//RGB
		im = image.NewRGBA(r)
	case 1:
		//Y
		im = image.NewGray(r)
	case 2:
		//Indexed
		if alpha {
			im = NewPalettedAlpha(r, d.props.colours)
		} else {
			im = image.NewPaletted(r, d.props.colours)
		}
	default:
		d.r.Err = ErrInvalidState
		return nil
	}

	return &LayerImage{
		layer: l,
		Image: im,
	}
}
