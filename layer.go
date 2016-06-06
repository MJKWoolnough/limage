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
	mask                                        struct {
		name                         string
		linked, lockContent, visible bool
		opacity                      uint8
		parasites                    parasites
		tattoo                       uint32
		active, selection, show      bool
		color                        struct {
			r, g, b uint8
		}
		image image.Image
	}
}

func (d *decoder) ReadLayer() layer {
	var l layer
	l.width = d.ReadUint32()
	l.height = d.ReadUint32()
	typ := d.ReadUint()
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

	d.Seek(int64(hptr))
	// read hierarchy

	l.image = d.ReadImage(l.width, l.height, l.mode)

	if mptr != 0 { // read layer mask
		d.Seek(int64(mptr))
		width := d.ReadUint32()
		height := d.ReadUint32()
		if width != l.width || height != l.height {
			d.SetError(ErrInconsistantData)
			return l
		}
		l.mask.name = d.ReadString()

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
				l.mask.linked = d.ReadBoolProperty()
			case propLockContent:
				l.mask.lockContent = d.ReadBoolProperty()
			case propOpacity:
				o := d.ReadUint32()
				if o > 255 {
					d.SetError(ErrInvalidOpacity)
				}
				l.mask.opacity = uint8(o)
			case propParasites:
				l.mask.parasites = d.ReadParasites()
			case propTattoo:
				l.mask.tattoo = d.ReadUint32()
			case propVisible:
				l.mask.visible = d.ReadBoolProperty()

				//mask properties
			case propActiveChannel:
				l.mask.active = true
			case propColor:
				l.mask.color.r = d.ReadUint8()
				l.mask.color.g = d.ReadUint8()
				l.mask.color.b = d.ReadUint8()
			case propSelection:
				l.mask.selection = true
			case propShowMasked:
				l.mask.show = d.ReadBoolProperty()
			default:
				d.Seek(int64(plength), os.SEEK_CUR)
			}

		}

		hptr := d.ReadUint32()
		d.Seek(int64(hptr))

		l.mask.image = d.ReadImage(l.width, l.height, l.mode) // gray???
	}
	return l
}

// Errors
var (
	ErrInvalidLayerType      = errors.New("invalid layer type")
	ErrInvalidItemPathLength = errors.New("invalid item path length")
	ErrInconsistantData      = errors.New("inconsistant data read")
)
