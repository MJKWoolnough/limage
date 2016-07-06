package lcolor

import "image/color"

// Color Models
var (
	GrayAlphaModel = color.ModelFunc(grayAlphaColourModel)
	RGBModel       = color.ModelFunc(rgbColourModel)
)
