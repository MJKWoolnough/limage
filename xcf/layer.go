package xcf

import (
	"github.com/MJKWoolnough/errors"
	"github.com/MJKWoolnough/limage"
)

type layer struct {
	limage.Layer
	width, height                               uint32
	linked, lockContent, alpha                  bool
	parasites                                   parasites
	tattoo                                      uint32
	apply, active, edit, group, lockAlpha, show bool
	selection                                   uint32
	itemPath                                    []rune
	groupItemFlags                              uint32
	textLayerFlags                              uint32
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
			l.linked = d.ReadBoolProperty()
		case propLockContent:
			l.lockContent = d.ReadBoolProperty()
		case propOpacity:
			o := d.ReadUint32()
			if o > 255 {
				d.SetError(ErrInvalidOpacity)
			}
			l.Transparency = 255 - uint8(o)
		case propParasites:
			l.parasites = d.ReadParasites(plength)
		case propTattoo:
			l.tattoo = d.ReadUint32()
		case propVisible:
			l.Invisible = !d.ReadBoolProperty()

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
			l.itemPath = make([]rune, plength>>2)
			for i := uint32(0); i < plength>>2; i++ {
				l.itemPath[i] = rune(d.ReadUint32())
			}
		case propGroupItemFlags:
			l.groupItemFlags = d.ReadUint32()
		case propLockAlpha:
			l.lockAlpha = d.ReadBoolProperty()
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
			l.show = d.ReadBoolProperty()
		case propTextLayerFlags:
			l.textLayerFlags = d.ReadUint32()
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

	l.Image = d.ReadImage(l.width, l.height, typ)

	if mptr != 0 { // read layer mask
		d.Goto(mptr)
		l.mask = d.ReadChannel()
		if l.mask.width != l.width || l.mask.height != l.height {
			d.SetError(ErrInconsistantData)
			return l
		}
	}
	return l
}

func (e *encoder) WriteLayers(layers []limage.Layer, groups []int32, w writer) {
	for n, layer := range layers {
		nGroups := append(groups, int32(n))
		w.WriteUint32(e.WriteLayer(layer, nGroups, w))
	}
}

func (e *encoder) WriteLayer(im limage.Layer, groups []int32, w writer) uint32 {
	var ptr uint32

	// write layer

	var g *limage.Group
	switch i := im.Image.(type) {
	case limage.Group:
		g = &i
	case *limage.Group:
		g = i
	default:
		return ptr
	}
	e.WriteLayers(g.Layers, groups, w)
	return ptr
}

// Errors
var (
	ErrInvalidLayerType      errors.Error = "invalid layer type"
	ErrInvalidItemPathLength errors.Error = "invalid item path length"
	ErrInconsistantData      errors.Error = "inconsistant data read"
)
