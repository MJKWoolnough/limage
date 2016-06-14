package xcf

import (
	"image"
	"image/color"
)

type Image struct {
	Group
	Comment string
}

type Layer struct {
	OffsetX, OffsetY uint32
	Mode             uint32
	image.Image
}

type Group struct {
	Name   string
	Layers []Layer
}

type MaskedImage struct {
	Image image.Image
	Mask  image.Image
}

type Text struct {
	image.Image
	TextData []TextData
}

type TextData struct {
	ForeColor, BackColor                   color.Color
	Bold, Italic, Underline, Strikethrough bool
	Text                                   string
}
