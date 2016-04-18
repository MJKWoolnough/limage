package xcf

import (
	"errors"
	"os"
)

type property uint8

const (
	propEnd               property = 0
	propColormap          property = 1
	propActiveLayer       property = 2
	propActiveChannel     property = 3
	propSelection         property = 4
	propFloatingSelection property = 5
	propOpacity           property = 6
	propMode              property = 7
	propVisible           property = 8
	propLinked            property = 9
	propLockAlpha         property = 10
	propApplyMask         property = 11
	propEditMask          property = 12
	propShowMask          property = 13
	propShowMasked        property = 14
	propOffsets           property = 15
	propColor             property = 16
	propCompression       property = 17
	propGuides            property = 18
	propResolution        property = 19
	propTattoo            property = 20
	propParasites         property = 21
	propUnit              property = 22
	propPaths             property = 23
	propUserUnit          property = 24
	propVectors           property = 25
	propTextLayerFlags    property = 26
	propSamplePoints      property = 27
	propLockContent       property = 28
	propGroupItem         property = 29
	propItemPath          property = 30
	propGroupItemFlags    property = 31
	propLockPosition      property = 32
	propFloatOpacity      property = 33
)

func (d *Decoder) readImageProperties(i *Image) {
	for {
		propID := property(d.r.ReadUint32())
		propLength := d.r.ReadUint32()
		switch propID {
		case propEnd:
			return
		case propColormap:
			d.colours = d.readColorMap()
		case propCompression:
			d.compression = d.readCompression()
		case propGuides:
			d.guides = d.readGuides(propLength)
		case propResolution:
			d.hres = d.r.ReadFloat32()
			d.vres = d.r.ReadFloat32()
		case propTattoo:
			d.tatoo = d.readTattoo()
		case propParasites:
			d.parasites = d.readParasites(propLength)
		case propUnit:
			d.unit = d.readUnit()
		case propPaths:
			d.paths = d.readPaths()
		case propUserUnit:
			d.userUnit = d.readUserUnit()
		case propVectors:
			d.vectors = d.readVectors()
		default:
			d.s.Seek(int64(propLength), os.SEEK_CUR)
		}
	}
}

func (d *Decoder) readChannelProperties() {
	for {
		propID := property(d.r.ReadUint32())
		propLength := d.r.ReadUint32()
		switch propID {
		case propEnd:
			return
		case propActiveChannel:
			//a := d.readActiveChannel()
			// no data, just set as active
		case propSelection:
			//s := d.readSelection()
			// no data, just set as selection
		case propOpacity:
			o := d.readOpacity()
			_ = o
		case propVisible:
			v := d.readBool()
			_ = v
		case propLinked:
			l := d.readBool()
			_ = l
		case propShowMasked:
			s := d.readBool()
			_ = s
		case propColor:
			r := d.r.ReadUint8()
			g := d.r.ReadUint8()
			b := d.r.ReadUint8()
			_, _, _ = r, g, b
		case propTattoo:
			t := d.readTattoo()
			_ = t
		case propParasites:
			p := d.readParasites(propLength)
			_ = p
		case propLockContent:
			l := d.readBool()
			_ = l
		default:
			d.s.Seek(int64(propLength), os.SEEK_CUR)
		}
	}
}

func (d *Decoder) readLayerProperties() {
	for {
		propID := property(d.r.ReadUint32())
		propLength := d.r.ReadUint32()
		switch propID {
		case propEnd:
			return
		case propActiveLayer:
			// a := d.readActiveLayer()
			// no data, just set as active layer
		case propFloatingSelection:
			f := d.r.ReadUint32()
			_ = f
		case propOpacity:
			o := d.readOpacity()
			_ = o
		case propApplyMask:
			a := d.readBool()
			_ = a
		case propEditMask:
			e := d.readBool()
			_ = e
		case propMode:
			m := d.readMode()
			_ = m
		case propLinked:
			l := d.readBool()
			_ = l
		case propLockAlpha:
			l := d.readBool()
			_ = l
		case propOffsets:
			x := d.r.ReadInt32()
			y := d.r.ReadInt32()
			_, _ = x, y
		case propShowMask:
			s := d.readBool()
			_ = s
		case propTattoo:
			t := d.readTattoo()
			_ = t
		case propParasites:
			p := d.readParasites(propLength)
			_ = p
		case propTextLayerFlags:
			t := d.readTextLayerFlags()
			_ = t
		case propLockContent:
			l := d.readBool()
			_ = l
		case propVisible:
			v := d.readBool()
			_ = v
		case propGroupItem:
			// g := d.readGroupItem()
			// no data, just set as item group
		case propItemPath:
			i := d.readItemPath(propLength)
			_ = i
		case propGroupItemFlags:
			g := d.r.ReadUint32() | 1
			_ = g
		default:
			d.s.Seek(int64(propLength), os.SEEK_CUR)
		}
	}
}

func (d *Decoder) readBool() bool {
	switch d.r.ReadUint32() {
	case 0:
		return false
	case 1:
		return true
	}
	d.r.Err = ErrInvalidState
	return false
}

// Errors
var (
	ErrInvalidState = errors.New("invalid state")
)
