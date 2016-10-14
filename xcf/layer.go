package xcf

import (
	"github.com/MJKWoolnough/errors"
	"github.com/MJKWoolnough/limage"
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
	width := d.ReadUint32()
	height := d.ReadUint32()
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
			d.ReadBoolProperty()
		case propLockContent:
			d.ReadBoolProperty()
		case propOpacity:
			o := d.ReadUint32()
			if o > 255 {
				d.SetError(ErrInvalidOpacity)
			}
			l.Transparency = 255 - uint8(o)
		case propParasites:
			parasites = d.ReadParasites(plength)
		case propTattoo:
			d.ReadUint32()
		case propVisible:
			l.Invisible = !d.ReadBoolProperty()

		//layer properties
		case propActiveLayer:
			// active layer
		case propApplyMask:
			d.ReadBoolProperty()
		case propEditMask:
			d.ReadBoolProperty()
		case propFloatingSelection:
			d.ReadUint32()
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
			d.ReadUint32()
		case propLockAlpha:
			d.ReadBoolProperty()
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
			l.OffsetX = int(d.ReadInt32())
			l.OffsetY = int(d.ReadInt32())
		case propShowMask:
			d.ReadBoolProperty()
		case propTextLayerFlags:
			d.ReadUint32()
		case propFloatOpacity:
			l.Transparency = 255 - uint8(d.ReadFloat32()*256)
		default:
			d.Skip(plength)
		}
	}

	hptr := d.ReadUint32()
	mptr := d.ReadUint32()

	d.Goto(hptr)
	// read hierarchy

	if l.group { // skip reading image if its a group
		return l
	}

	l.Image = d.ReadImage(width, height, typ)

	if mptr != 0 { // read layer mask
		d.Goto(mptr)
		var m limage.MaskedImage
		m.Image = l.Image
		m.Mask = d.ReadChannel()
		b := m.Mask.Bounds()
		if uint32(b.Dx()) != width || uint32(b.Dy()) != height {
			d.SetError(ErrInconsistantData)
			return l
		}
		l.Image = m
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
	return l
}

func (e *encoder) WriteLayers(layers limage.Image, groups []int32, w writer) {
	for n, layer := range layers {
		nGroups := append(groups, int32(n))
		w.WriteUint32(e.WriteLayer(layer, nGroups, w))
	}
}

func (e *encoder) WriteLayer(im limage.Layer, groups []int32, w writer) uint32 {
	var ptr uint32

	// write layer

	var g *limage.Image
	switch i := im.Image.(type) {
	case limage.Image:
		g = &i
	case *limage.Image:
		g = i
	default:
		return ptr
	}
	e.WriteLayers(*g, groups, w)
	return ptr
}

// Errors
const (
	ErrInvalidLayerType      errors.Error = "invalid layer type"
	ErrInvalidItemPathLength errors.Error = "invalid item path length"
	ErrInconsistantData      errors.Error = "inconsistant data read"
)
