package xcf

import (
	"image"

	"vimagination.zapto.org/limage"
)

func (e *encoder) WriteLayers(layers limage.Image, offsetX, offsetY int32, groups []uint32, pw *pointerWriter) {
	for n, layer := range layers {
		nGroups := append(groups, uint32(n))

		e.WriteLayer(layer, offsetX+int32(layer.LayerBounds.Min.X), offsetY+int32(layer.LayerBounds.Min.Y), nGroups, pw)
	}
}

func (e *encoder) WriteLayer(im limage.Layer, offsetX, offsetY int32, groups []uint32, pw *pointerWriter) {
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
	e.WriteInt32(offsetX)
	e.WriteInt32(offsetY)

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

	e.WriteUint32(propMode)
	e.WriteUint32(4)

	switch im.Mode {
	case limage.CompositeNormal:
		e.WriteUint32(0)
	case limage.CompositeDissolve:
		e.WriteUint32(1)
	case limage.CompositeBehind:
		e.WriteUint32(2)
	case limage.CompositeMultiply:
		e.WriteUint32(3)
	case limage.CompositeScreen:
		e.WriteUint32(4)
	case limage.CompositeOverlay:
		e.WriteUint32(5)
	case limage.CompositeDifference:
		e.WriteUint32(6)
	case limage.CompositeAddition:
		e.WriteUint32(7)
	case limage.CompositeSubtract:
		e.WriteUint32(8)
	case limage.CompositeDarkenOnly:
		e.WriteUint32(9)
	case limage.CompositeLightenOnly:
		e.WriteUint32(10)
	case limage.CompositeHue:
		e.WriteUint32(11)
	case limage.CompositeSaturation:
		e.WriteUint32(12)
	case limage.CompositeColor:
		e.WriteUint32(13)
	case limage.CompositeValue:
		e.WriteUint32(14)
	case limage.CompositeDivide:
		e.WriteUint32(15)
	case limage.CompositeDodge:
		e.WriteUint32(16)
	case limage.CompositeBurn:
		e.WriteUint32(17)
	case limage.CompositeHardLight:
		e.WriteUint32(18)
	case limage.CompositeSoftLight:
		e.WriteUint32(19)
	case limage.CompositeGrainExtract:
		e.WriteUint32(20)
	case limage.CompositeGrainMerge:
		e.WriteUint32(21)
	case limage.CompositeLuminosity:
		e.WriteUint32(22)
	case limage.CompositeColorErase:
		e.WriteUint32(22)
	case limage.CompositeChroma:
		e.WriteUint32(25)
	case limage.CompositeLightness:
		e.WriteUint32(27)
	case limage.CompositeVividLight:
		e.WriteUint32(48)
	case limage.CompositePinLight:
		e.WriteUint32(49)
	case limage.CompositeLinearLight:
		e.WriteUint32(50)
	case limage.CompositeHardMix:
		e.WriteUint32(51)
	case limage.CompositeExclusion:
		e.WriteUint32(52)
	case limage.CompositeLinearBurn:
		e.WriteUint32(53)
	case limage.CompositeErase:
		e.WriteUint32(58)
	case limage.CompositeMerge:
		e.WriteUint32(59)
	case limage.CompositeSplit:
		e.WriteUint32(60)
	case limage.CompositePassThrough:
		e.WriteUint32(61)
	default:
		e.WriteUint32(0) // Error instead?
	}

	e.WriteUint32(0) // end of properties
	e.WriteUint32(0)

	// write layer

	ptrs := e.ReservePointers(2)

	ptrs.WritePointer(uint32(e.pos))

	e.WriteImage(img, e.colourFunc, e.colourChannels)

	if mask != nil {
		ptrs.WritePointer(uint32(e.pos))
		e.WriteChannel(mask)
	} else {
		ptrs.WritePointer(0)
	}

	if group != nil {
		e.WriteLayers(group, offsetX, offsetY, groups, pw)
	}
}
