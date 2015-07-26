package ora

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"image"
	_ "image/png"
	"io"
)

type Layer struct {
	name string
	file *zip.File
}

func (l *Layer) Image() (image.Image, error) {
	f, err := l.file.Open()
	if err != nil {
		return nil, err
	}
	i, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return i, nil
}

type ReadCloserAt interface {
	io.ReaderAt
	io.Closer
}

type ORA struct {
	io.Closer
	thumbnail, merged *zip.File
	structure         imageStructure
}

func OpenFile(f string) (*ORA, error) {
	r, err := zip.OpenReader(f)
	if err != nil {
		return nil, err
	}
	o, err := open(r.File)
	if err != nil {
		return nil, err
	}
	o.Closer = r
	return o, err
}

func Open(rca ReadCloserAt, size int64) (*ORA, error) {
	r, err := zip.NewReader(rca, size)
	if err != nil {
		return nil, err
	}
	o, err := open(r.File)
	if err != nil {
		return nil, err
	}
	o.Closer = rca
	return o, err
}

type Stack struct {
	Stack  []Stack `xml:"stack"`
	Name   string  `xml:"name,attr"`
	X      int     `xml:"x,attr"`
	Y      int     `xml:"y,attr"`
	Layers []struct {
		Name   string `xml:"name,attr"`
		Source string `xml:"src,attr"`
		X      int    `xml:"x,attr"`
		Y      int    `xml:"y,attr"`
	} `xml:"layer"`
}

func (s Stack) layer(name string) *Layer {
	return nil
}

type imageStructure struct {
	Width  int  `xml:"w,attr"`
	Height int  `xml:"h,attr"`
	Xres   uint `xml:"xres,attr"`
	Yres   uint `xml:"yres,attr"`
	Stack
}

func open(files []*zip.File) (*ORA, error) {
	o := new(ORA)
	for _, file := range files {
		switch file.Name {
		case "Thumbnails/thumbnail.png":
			o.thumbnail = file
		case "mergedimage.png":
			o.merged = file
		case "stack.xml":
			f, err := file.Open()
			if err != nil {
				return nil, err
			}
			var im imageStructure
			im.Xres = 72
			im.Yres = 72
			err = xml.NewDecoder(f).Decode(&im)
			if err != nil {
				return nil, err
			}
			f.Close()
			o.structure = im
		}
	}
	if o.thumbnail == nil {
		return nil, ErrInvalidLayout
	}
	return o, nil
}

func (o *ORA) Thumbnail() (image.Image, error) {
	f, err := o.thumbnail.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return i, err
}

func (o *ORA) Merged() (image.Image, error) {
	if o.merged == nil {
		return nil, ErrNoMerged
	}
	f, err := o.merged.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return i, err
}

func (o *ORA) Layer(name string) *Layer {
	return nil
}

// Errors
var (
	ErrInvalidLayout = errors.New("invalid layout")
	ErrNoMerged      = errors.New("ora does not include a merged image")
	ErrNoLayer       = errors.New("no longer with that name exists in ora")
)
