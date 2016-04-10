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
		case propCompression:
		case propGuides:
		case propResolution:
		case propTattoo:
		case propParasites:
		case propUnit:
		case propPaths:
		case propUserUnit:
		case propVectors:
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
		case propSelection:
		case propOpacity:
		case propVisible:
		case propLinked:
		case propShowMasked:
		case propColor:
		case propTattoo:
		case propParasites:
		case propLockContent:
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
		case propFloatingSelection:
		case propOpacity:
		case propApplyMask:
		case propEditMask:
		case propMode:
		case propLinked:
		case propLockAlpha:
		case propOffsets:
		case propShowMask:
		case propTattoo:
		case propParasites:
		case propTextLayerFlags:
		case propLockContent:
		case propVisible:
		case propGroupItem:
		case propItemPath:
		case propGroupItemFlags:
		default:
			d.s.Seek(int64(propLength), os.SEEK_CUR)
		}
	}
}
