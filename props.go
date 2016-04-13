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
			r := d.readResolution()
		case propTattoo:
			t := d.readTattoo()
		case propParasites:
			p := d.readParasites()
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
			a := d.readActiveChannel()
		case propSelection:
			s := d.readSelection()
		case propOpacity:
			o := d.readOpacity()
		case propVisible:
			v := d.readVisible()
		case propLinked:
			l := d.readLinked()
		case propShowMasked:
			s := d.readShowMasked()
		case propColor:
			c := d.readColor()
		case propTattoo:
			t := d.readTattoo()
		case propParasites:
			p := d.readParasites()
		case propLockContent:
			l := d.readLockContent()
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
			a := d.readActiveLayer()
		case propFloatingSelection:
			f := d.readFloatingSelection()
		case propOpacity:
			o := d.readOpacity()
		case propApplyMask:
			a := d.readApplyMask()
		case propEditMask:
			e := d.readEditMask()
		case propMode:
			m := d.readPropMode()
		case propLinked:
			l := d.readLinked()
		case propLockAlpha:
			l := d.readLockAlpha()
		case propOffsets:
			o := d.readOffsets()
		case propShowMask:
			s := d.readShowMask()
		case propTattoo:
			t := d.readTattoo()
		case propParasites:
			p := d.readParasites()
		case propTextLayerFlags:
			t := d.readTextLayerFlags()
		case propLockContent:
			l := d.readLockContent()
		case propVisible:
			v := d.readVisible()
		case propGroupItem:
			g := d.readGroupItem()
		case propItemPath:
			i := d.readItemPath()
		case propGroupItemFlags:
			g := d.readGroupItemFlags()
		default:
			d.s.Seek(int64(propLength), os.SEEK_CUR)
		}
	}
}
