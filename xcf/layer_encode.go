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

	writeProperties(e, im, offsetX, offsetY, groups, group, text)
	e.WriteUint32(0)
	writeLayer(e, img, mask)

	if group != nil {
		e.WriteLayers(group, offsetX, offsetY, groups, pw)
	}
}

func writeProperties(e *encoder, im limage.Layer, offsetX, offsetY int32, groups []uint32, group limage.Image, text limage.TextData) {
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

	e.WriteUint32(modeID(im.Mode))

	e.WriteUint32(0) // end of properties
}

func modeID(mode limage.Composite) uint32 {
	switch mode {
	case limage.CompositeNormal:
		return 0
	case limage.CompositeDissolve:
		return 1
	case limage.CompositeBehind:
		return 2
	case limage.CompositeMultiply:
		return 3
	case limage.CompositeScreen:
		return 4
	case limage.CompositeOverlay:
		return 5
	case limage.CompositeDifference:
		return 6
	case limage.CompositeAddition:
		return 7
	case limage.CompositeSubtract:
		return 8
	case limage.CompositeDarkenOnly:
		return 9
	case limage.CompositeLightenOnly:
		return 10
	case limage.CompositeHue:
		return 11
	case limage.CompositeSaturation:
		return 12
	case limage.CompositeColor:
		return 13
	case limage.CompositeValue:
		return 14
	case limage.CompositeDivide:
		return 15
	case limage.CompositeDodge:
		return 16
	case limage.CompositeBurn:
		return 17
	case limage.CompositeHardLight:
		return 18
	case limage.CompositeSoftLight:
		return 19
	case limage.CompositeGrainExtract:
		return 20
	case limage.CompositeGrainMerge:
		return 21
	case limage.CompositeLuminosity:
		return 22
	case limage.CompositeColorErase:
		return 22
	case limage.CompositeChroma:
		return 25
	case limage.CompositeLightness:
		return 27
	case limage.CompositeVividLight:
		return 48
	case limage.CompositePinLight:
		return 49
	case limage.CompositeLinearLight:
		return 50
	case limage.CompositeHardMix:
		return 51
	case limage.CompositeExclusion:
		return 52
	case limage.CompositeLinearBurn:
		return 53
	case limage.CompositeErase:
		return 58
	case limage.CompositeMerge:
		return 59
	case limage.CompositeSplit:
		return 60
	case limage.CompositePassThrough:
		return 61
	default:
		return 0 // Error instead?
	}
}

func writeLayer(e *encoder, img image.Image, mask *image.Gray) {
	ptrs := e.ReservePointers(2)

	ptrs.WritePointer(uint32(e.pos))

	e.WriteImage(img, e.colourFunc, e.colourChannels)

	if mask != nil {
		ptrs.WritePointer(uint32(e.pos))
		e.WriteChannel(mask)
	} else {
		ptrs.WritePointer(0)
	}
}
