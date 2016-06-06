package xcf

import (
	"errors"
	"os"
)

type layer struct {
	width, height                               uint32
	name                                        string
	linked, lockContent, visible                bool
	opacity                                     uint8
	parasites                                   parasites
	tattoo                                      uint32
	apply, active, edit, group, lockAlpha, show bool
	selection                                   uint32
	itemPath                                    []uint32
	groupItemFlags                              uint32
	mode                                        uint32
	offsetX, offsetY                            int32
	textLayerFlags                              uint32
}

func (d *decoder) ReadLayer() layer {
	var l layer
	l.width = d.ReadUint32()
	l.height = d.ReadUint32()
	typ := d.ReadUint()
	if typ>>1 != d.baseType {
		d.Err = ErrInvalidLayerType
		return l
	}
	l.name = d.ReadString()

	// read properties
PropertyLoop:
	for {
		typ := d.ReadUint32()
		plength := d.ReadUint32()
		switch typ {
		// general properties
		case propEnd:
			if plength != 0 {
				d.Err = ErrInvalidProperties
			}
			break PropertyLoop
		case propLinked:
			l.linked = d.ReadBoolProperty()
		case propLockContent:
			l.lockContent = d.ReadBoolProperty()
		case propOpacity:
			o := d.ReadUint32()
			if o > 255 {
				d.Err = ErrInvalidOpacity
			}
			d.opacity = uint8(o)
		case propParasites:
			l.parasites = d.ReadParasites()
		case propTattoo:
			l.tattoo = d.ReadUint32()
		case propVisible:
			l.visible = d.ReadBoolProperty()

		//layer properties
		case propActiveLayer:
			l.active = true
		case propApplyMask:
			l.apply = d.ReadBoolProperty()
		case propEditMask:
			l.edit = d.ReadBoolProperty()
		case propFloatingSelection:
			l.selection = d.ReadUint32()
		case propGroupItem:
			l.group = d.ReadBoolProperty()
		case propItemPath:
			if plength&3 != 0 {
				d.Err = ErrInvalidItemPathLength
			}
			l.itemPath = make([]uint32, plength>>2)
			for i := uint32(0); i < plength>>2; i++ {
				l.itemPath[i] = d.ReadUint32()
			}
		case propGroupItemFlags:
			l.groupItemFlags = d.ReadUint32()
		case propLockAlpha:
			l.lockAlpha = d.ReadBoolProperty()
		case propMode:
			l.mode = d.ReadUint32()
		case propOffsets:
			l.offsetX = d.ReadInt32()
			l.offsetY = d.ReadInt32()
		case propShowMask:
			l.show = d.ReadBoolProperty()
		case propTextLayerFlags:
			l.textLayerFlags = d.ReadUint32()
		default:
			d.Seek(int64(plength), os.SEEK_CUR)
		}
	}
}

// Errors
var (
	ErrInvalidLayerType      = errors.New("invalid layer type")
	ErrInvalidItemPathLength = errors.New("invalid item path length")
)
