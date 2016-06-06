package xcf

import "errors"

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

func (d *decoder) ReadBoolProperty() bool {
	switch d.ReadUint32() {
	case 0:
		return false
	case 1:
		return true
	default:
		d.SetError(ErrInvalidBoolean)
		return false
	}
}

// Errors
var (
	ErrInvalidBoolean = errors.New("invalid boolean value")
)
