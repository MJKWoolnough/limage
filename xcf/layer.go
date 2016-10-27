package xcf

import (
	"image"

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
	if l.Image == nil {
		return l
	}

	if mptr != 0 { // read layer mask
		d.Goto(mptr)
		var m limage.MaskedImage
		m.Image = l.Image
		m.Mask = d.ReadChannel()
		if m.Mask == nil {
			return l
		}
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

func (e *encoder) WriteLayers(layers limage.Image, groups []uint32, pw *pointerWriter) {
	for n, layer := range layers {
		nGroups := append(groups, uint32(n))
		e.WriteLayer(layer, nGroups, pw)
	}
}

func (e *encoder) WriteLayer(im limage.Layer, groups []uint32, pw *pointerWriter) {
	pw.WritePointer(uint32(e.pos))

	var (
		mask  *image.Gray
		img   image.Image
		text  limage.TextData
		group limage.Image
	)
	if mim, ok := im.Image.(limage.MaskedImage); ok {
		mask = mim.Mask
		img = mim.Image
	} else if mim, ok := im.Image.(*limage.MaskedImage); ok {
		mask = mim.Mask
		img = mim.Image
	} else {
		img = im.Image
	}

	switch i := im.Image.(type) {
	case limage.Text:
		text = i.TextData
	case *limage.Text:
		text = i.TextData
	case limage.Image:
		group = i
	case *limage.Image:
		group = *i
	}

	b := im.Bounds()
	dx, dy := uint32(b.Dx()), uint32(b.Dy())
	e.WriteUint32(dx)
	e.WriteUint32(dy)
	e.WriteUint32(uint32(e.colourType)<<1 | 1)
	e.WriteString(im.Name)

	e.WriteUint32(propOpacity)
	e.WriteUint32(4)
	e.WriteUint32(255 - uint32(im.Transparency))

	e.WriteUint32(propVisible)
	e.WriteUint32(4)
	if im.Invisible {
		e.WriteUint32(0)
	} else {
		e.WriteUint32(1)
	}

	e.WriteUint32(propOffsets)
	e.WriteUint32(8)
	e.WriteUint32(uint32(im.OffsetX))
	e.WriteUint32(uint32(im.OffsetY))

	if len(groups) > 1 {
		e.WriteUint32(propItemPath)
		e.WriteUint32(4 * uint32(len(groups)))
		for _, g := range groups {
			e.WriteUint32(g)
		}
	}

	if len(text) > 0 {
		e.WriteText(text, dx, dy)
	}

	if group != nil {
		e.WriteUint32(propGroupItem)
		e.WriteUint32(0)
	}

	e.WriteUint32(0) // end of properties
	e.WriteUint32(0)

	// write layer

	e.WriteImage(img, e.colourFunc, e.colourChannels)
	if mask != nil {
		e.WriteChannel(mask)
	}
	if group != nil {
		e.WriteLayers(group, groups, pw)
	}
}

// Errors
const (
	ErrInvalidLayerType      errors.Error = "invalid layer type"
	ErrInvalidItemPathLength errors.Error = "invalid item path length"
	ErrInconsistantData      errors.Error = "inconsistant data read"
)
