package xcf

import (
	"image"

	"vimagination.zapto.org/errors"
	"vimagination.zapto.org/limage"
)

type layer struct {
	limage.Layer
	alpha    bool
	group    bool
	itemPath []rune
}

func (d *decoder) ReadLayer() layer {
	var (
		l         layer
		parasites parasites
	)
	l.LayerBounds.Max.X = int(d.ReadUint32())
	l.LayerBounds.Max.Y = int(d.ReadUint32())
	typ := d.ReadUint32()
	if typ>>1 != d.baseType {
		d.SetError(ErrInvalidLayerType)
		return l
	}
	l.alpha = typ&1 == 1
	l.Name = d.ReadString()

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
			d.SkipBoolProperty()
		case propLockContent:
			d.SkipBoolProperty()
		case propOpacity:
			o := d.ReadUint32()
			if o > 255 {
				d.SetError(ErrInvalidOpacity)
			}
			l.Transparency = 255 - uint8(o)
		case propParasites:
			parasites = d.ReadParasites(plength)
		case propTattoo:
			d.SkipUint32()
		case propVisible:
			l.Invisible = !d.ReadBoolProperty()

		//layer properties
		case propActiveLayer:
			// active layer
		case propApplyMask:
			d.SkipBoolProperty()
		case propEditMask:
			d.SkipBoolProperty()
		case propFloatingSelection:
			d.SkipUint32()
		case propGroupItem:
			l.group = true
		case propItemPath:
			if plength&3 != 0 {
				d.SetError(ErrInvalidItemPathLength)
			}
			l.itemPath = make([]rune, plength>>2)
			for i := uint32(0); i < plength>>2; i++ {
				l.itemPath[i] = rune(d.ReadUint32())
			}
		case propGroupItemFlags:
			d.SkipUint32()
		case propLockAlpha:
			d.SkipBoolProperty()
		case propMode:
			l.Mode = limage.Composite(d.ReadUint32())
			if d.baseType != 0 {
				switch l.Mode {
				case limage.CompositeNormal, limage.CompositeBehind:
				default:
					l.Mode = 0
				}
			}
		case propOffsets:
			offsetX := int(d.ReadInt32())
			offsetY := int(d.ReadInt32())
			l.LayerBounds = l.LayerBounds.Add(image.Pt(offsetX, offsetY))
		case propShowMask:
			d.SkipBoolProperty()
		case propTextLayerFlags:
			d.SkipUint32()
		case propFloatOpacity:
			l.Transparency = 255 - uint8(d.ReadFloat32()*255)
		default:
			d.Skip(plength)
		}
	}

	var hptr, mptr uint64
	if d.mode < 2 {
		hptr = uint64(d.ReadUint32())
		mptr = uint64(d.ReadUint32())
	} else {
		hptr = d.ReadUint64()
		mptr = d.ReadUint64()
	}

	d.Goto(hptr)
	// read hierarchy

	if !l.group { // skip reading image if its a group
		l.Image = d.ReadImage(uint32(l.LayerBounds.Dx()), uint32(l.LayerBounds.Dy()), typ)
		if l.Image == nil {
			return l
		}
	}
	if t := parasites.Get(textParasiteName); t != nil {
		textData, err := parseTextData(t)
		if err != nil {
			d.SetError(ErrInvalidLayerType)
			return l
		}
		l.Image = limage.Text{
			Image:    l.Image,
			TextData: textData,
		}
	}
	if mptr != 0 { // read layer mask
		d.Goto(mptr)
		var m limage.MaskedImage
		m.Image = l.Image
		m.Mask = d.ReadChannel()
		if m.Mask == nil {
			return l
		}
		if !l.LayerBounds.Eq(m.Mask.Bounds()) {
			d.SetError(ErrInconsistantData)
			return l
		}
		l.Image = m
	}
	return l
}

// Errors
const (
	ErrInvalidLayerType      errors.Error = "invalid layer type"
	ErrInvalidItemPathLength errors.Error = "invalid item path length"
	ErrInconsistantData      errors.Error = "inconsistant data read"
)
