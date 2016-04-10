package ora

import "encoding/xml"

type imageContent struct {
	Width  int          `xml:"w,attr,required"`
	Height int          `xml:"h,attr,required"`
	Name   string       `xml:"name,attr"`
	Stack  stackContent `xml:"stack"`
}

type stackContent struct {
	layerCommonAttributes
	Stack []struct {
		XMLName xml.Name
		layerCommonAttributes
		stackContent  `xml:"stack"`
		layerContent  `xml:"layer"`
		filterContent `xml:"filter"`
		textContent   `xml:"text"`
	} `xml:",any"`
}

type layerContent struct {
	Source      string          `xml:"src,attr"`
	CompositeOp string          `xml:"composite-op,attr"`
	Opacity     float32         `xml:"opacity,attr"`
	Filters     []filterContent `xml:"filter"`
}

type filterContent struct {
	Type   string        `xml:"type,attr"`
	Output string        `xml:"output,attr"`
	Params paramsContent `xml:"params"`
	Stack  stackContent  `xml:"stack"`
}

type textContent struct {
	Data string `xml:",chardata"`
}

type layerCommonAttributes struct {
	X    int    `xml:"x,attr"`
	Y    int    `xml:"y,attr"`
	Name string `xml:"name,attr"`
}

type paramsContent struct {
	Version int `xml:"version"`
	Params  []struct {
		Name string `xml:"name"`
		Data string `xml:",chardata"`
	}
}
