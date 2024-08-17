package xcf

import (
	"errors"
	"image"

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

		// layer properties
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
			if d.baseType != 0 {
				switch d.ReadUint32() {
				case 2: // Behind
					l.Mode = limage.CompositeBehind
				default:
					l.Mode = limage.CompositeNormal
				}
			} else {
				switch d.ReadUint32() {
				case 0, 28:
					l.Mode = limage.CompositeNormal
				case 1:
					l.Mode = limage.CompositeDissolve
				case 2, 29:
					l.Mode = limage.CompositeBehind
				case 3, 30:
					l.Mode = limage.CompositeMultiply
				case 4, 31:
					l.Mode = limage.CompositeScreen
				case 5, 23:
					l.Mode = limage.CompositeOverlay
				case 6, 32:
					l.Mode = limage.CompositeDifference
				case 7, 33:
					l.Mode = limage.CompositeAddition
				case 8, 34:
					l.Mode = limage.CompositeSubtract
				case 9, 35, 54:
					l.Mode = limage.CompositeDarkenOnly
				case 10, 36, 55:
					l.Mode = limage.CompositeLightenOnly
				case 11, 24, 37:
					l.Mode = limage.CompositeHue
				case 12, 38:
					l.Mode = limage.CompositeSaturation
				case 13, 26, 39:
					l.Mode = limage.CompositeColor
				case 14, 40:
					l.Mode = limage.CompositeValue
				case 15, 41:
					l.Mode = limage.CompositeDivide
				case 16, 42:
					l.Mode = limage.CompositeDodge
				case 17, 43:
					l.Mode = limage.CompositeBurn
				case 18, 44:
					l.Mode = limage.CompositeHardLight
				case 19, 45:
					l.Mode = limage.CompositeSoftLight
				case 20, 46:
					l.Mode = limage.CompositeGrainExtract
				case 21, 47:
					l.Mode = limage.CompositeGrainMerge
				case 22, 57:
					l.Mode = limage.CompositeColorErase
				case 25:
					l.Mode = limage.CompositeChroma
				case 27:
					l.Mode = limage.CompositeLightness
				case 48:
					l.Mode = limage.CompositeVividLight
				case 49:
					l.Mode = limage.CompositePinLight
				case 50:
					l.Mode = limage.CompositeLinearLight
				case 51:
					l.Mode = limage.CompositeHardMix
				case 52:
					l.Mode = limage.CompositeExclusion
				case 53:
					l.Mode = limage.CompositeLinearBurn
				case 56:
					l.Mode = limage.CompositeLuminance
				case 58:
					l.Mode = limage.CompositeErase
				case 59:
					l.Mode = limage.CompositeMerge
				case 60:
					l.Mode = limage.CompositeSplit
				case 61:
					l.Mode = limage.CompositePassThrough
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
		if l.Image = d.ReadImage(uint32(l.LayerBounds.Dx()), uint32(l.LayerBounds.Dy()), typ); l.Image == nil {
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

		if m.Mask = d.ReadChannel(); m.Mask == nil {
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

// Errors.
var (
	ErrInvalidLayerType      = errors.New("invalid layer type")
	ErrInvalidItemPathLength = errors.New("invalid item path length")
	ErrInconsistantData      = errors.New("inconsistant data read")
)
