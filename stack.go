package ora

import "encoding/xml"

type image struct {
	Width  int          `xml:"w,attribute,required"`
	Height int          `xml:"h,attribute,required`
	Name   string       `xml:"name,attribute"`
	Stack  stackContent `xml:"stack"`
}

type stackContent struct {
	layerCommonAttributes
	Stack []struct {
		XMLName xml.Name
		layerCommonAttributes
		*stackContent  `xml:"stack"`
		*layerContent  `xml:"layer"`
		*filterContent `xml:"filter"`
		*textContent   `xml:"text"`
	} `xml:",any"`
}

type layerContent struct {
	Source      string          `xml:"src,attribute"`
	CompositeOp string          `xml:"composite-op,attribute"`
	Opacity     float32         `xml:"opacity,attribute"`
	Filters     []filterContent `xml:"filter"`
}

type filterContent struct {
	Type   string        `xml:"type,attribute"`
	Output string        `xml:"type,attribute"`
	Params paramsContent `xml:"params"`
	Stack  stackContent  `xml:"stack"`
}

type text struct {
	Data string `xml:",chardata"`
}

type layerCommonAttributes struct {
	X    int    `xml:"x"`
	Y    int    `xml:"y"`
	Name string `xml:"name"`
}

type paramsContent struct {
	Version int `xml:"version"`
	Params  []struct {
		Name string `xml:"name"`
		Data string `xml:",chardata"`
	}
}
