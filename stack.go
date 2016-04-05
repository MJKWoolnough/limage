package ora

type image struct {
	Width  int          `xml:"w,attribute,required"`
	Height int          `xml:"h,attribute,required`
	Name   string       `xml:"name,attribute"`
	Stack  stackContent `xml:"stack"`
}

type stackContent struct {
	layerCommonAttributes
	Stacks  []stackContent  `xml:"stack"`
	Layers  []layerContent  `xml:"layer"`
	Filters []filterContent `xml:"filter"`
	Texts   []textContent   `xml:"text"`
}

type layerContent struct {
	layerCommonAttributes
	Source      string          `xml:"src,attribute"`
	CompositeOp string          `xml:"composite-op,attribute"`
	Opacity     float32         `xml:"opacity,attribute"`
	Filters     []filterContent `xml:"filter"`
}

type filterContent struct {
	layerCommonAttributes
	Type   string        `xml:"type,attribute"`
	Output string        `xml:"type,attribute"`
	Params paramsContent `xml:"params"`
	Stack  stackContent  `xml:"stack"`
}

type text struct {
	layerCommonAttributes
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
