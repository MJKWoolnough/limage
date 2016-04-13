package xcf

import "os"

const (
	propEnd               = 0
	propColormap          = 1
	propActiveLayer       = 2
	propActiveChannel     = 3
	propSelection         = 4
	propFloatingSelection = 5
	propOpacity           = 6
	propMode              = 7
	propVisible           = 8
	propLinked            = 9
	propLockAlpha         = 10
	propApplyMask         = 11
	propEditMask          = 12
	propShowMask          = 13
	propShowMasked        = 14
	propOffsets           = 15
	propColor             = 16
	propCompression       = 17
	propGuides            = 18
	propResolution        = 19
	propTattoo            = 20
	propParasites         = 21
	propUnit              = 22
	propPaths             = 23
	propUserUnit          = 24
	propVectors           = 25
	propTextLayerFlags    = 26
	propSamplePoints      = 27
	propLockContent       = 28
	propGroupItem         = 29
	propItemPath          = 30
	propGroupItemFlags    = 31
	propLockPosition      = 32
	propFloatOpacity      = 33
)

const (
	compNone    = 0
	compRLE     = 1
	compZlib    = 2
	compFractal = 3
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
			g := d.readGuides()
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
