package xcf

import (
	"errors"
	"image"
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
	image                                       image.Image
	mask                                        channel
}

func (d *decoder) ReadLayer() layer {
	var l layer
	l.width = d.ReadUint32()
	l.height = d.ReadUint32()
	typ := d.ReadUint32()
	if typ>>1 != d.baseType {
		d.SetError(ErrInvalidLayerType)
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
				d.SetError(ErrInvalidProperties)
			}
			break PropertyLoop
		case propLinked:
			l.linked = d.ReadBoolProperty()
		case propLockContent:
			l.lockContent = d.ReadBoolProperty()
		case propOpacity:
			o := d.ReadUint32()
			if o > 255 {
				d.SetError(ErrInvalidOpacity)
			}
			l.opacity = uint8(o)
		case propParasites:
			l.parasites = d.ReadParasites(plength)
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
			l.group = true
		case propItemPath:
			if plength&3 != 0 {
				d.SetError(ErrInvalidItemPathLength)
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

	hptr := d.ReadUint32()
	mptr := d.ReadUint32()

	d.Seek(int64(hptr), os.SEEK_SET)
	// read hierarchy

	l.image = d.ReadImage(l.width, l.height, l.mode)

	if mptr != 0 { // read layer mask
		d.Seek(int64(mptr), os.SEEK_SET)
		l.mask = d.ReadChannel()
		if l.mask.width != l.width || l.mask.height != l.height {
			d.SetError(ErrInconsistantData)
			return l
		}
	}
	return l
}

// Errors
var (
	ErrInvalidLayerType      = errors.New("invalid layer type")
	ErrInvalidItemPathLength = errors.New("invalid item path length")
	ErrInconsistantData      = errors.New("inconsistant data read")
)
