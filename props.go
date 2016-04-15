package xcf

import "os"

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
		propID := d.r.ReadUint32()
		propLength := d.r.ReadUint32()
		switch propID {
		case propEnd:
			return
		case propColormap:
			c := d.readColorMap()
		case propCompression:
			c := d.readCompression()
		case propGuides:
			g := d.readGuides(propLength)
		case propResolution:
			h := d.r.ReadFloat32()
			v := d.r.ReadFloat32()
		case propTattoo:
			t := d.readTattoo()
		case propParasites:
			p := d.readParasites(propLength)
		case propUnit:
			u := d.readUnit()
		case propPaths:
			p := d.readPaths()
		case propUserUnit:
			u := d.readUserUnit()
		case propVectors:
			v := d.readVectors()
		default:
			d.s.Seek(int64(propLength), os.SEEK_CUR)
		}
	}
}

func (d *Decoder) readChannelProperties() {
	for {
		propID := d.r.ReadUint32()
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
		case propVisible:
			v := d.readBool()
		case propLinked:
			l := d.readBool()
		case propShowMasked:
			s := d.readBool()
		case propColor:
			r := d.r.ReadUint8()
			g := d.r.ReadUint8()
			b := d.r.ReadUint8()
		case propTattoo:
			t := d.readTattoo()
		case propParasites:
			p := d.readParasites(propLength)
		case propLockContent:
			l := d.readBool()
		default:
			d.s.Seek(int64(propLength), os.SEEK_CUR)
		}
	}
}

func (d *Decoder) readLayerProperties() {
	for {
		propID := d.r.ReadUint32()
		propLength := d.r.ReadUint32()
		switch propID {
		case propEnd:
			return
		case propActiveLayer:
			// a := d.readActiveLayer()
			// no data, just set as active layer
		case propFloatingSelection:
			f := d.r.ReadUint32()
		case propOpacity:
			o := d.readOpacity()
		case propApplyMask:
			a := d.readBool()
		case propEditMask:
			e := d.readBool()
		case propMode:
			m := d.readPropMode()
		case propLinked:
			l := d.readBool()
		case propLockAlpha:
			l := d.readBool()
		case propOffsets:
			x := d.r.ReadInt32()
			y := d.r.ReadInt32()
		case propShowMask:
			s := d.readBool()
		case propTattoo:
			t := d.readTattoo()
		case propParasites:
			p := d.readParasites(propLength)
		case propTextLayerFlags:
			t := d.readTextLayerFlags()
		case propLockContent:
			l := d.readBool()
		case propVisible:
			v := d.readBool()
		case propGroupItem:
			// g := d.readGroupItem()
			// no data, just set as item group
		case propItemPath:
			i := d.readItemPath(propLength)
		case propGroupItemFlags:
			g := d.r.ReadUint32() | 1
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
