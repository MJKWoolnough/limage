package ora

import (
	"archive/zip"
	"errors"
	"image"
	_ "image/png"
	"io"
)

type ReadCloserAt interface {
	io.ReaderAt
	io.Closer
}

type ORA struct {
	io.Closer
	structure *imageStack
	files     map[string]*zip.File
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

func open(files []*zip.File) (*ORA, error) {
	o := &ORA{
		files: make(map[string]*zip.File),
	}
	for _, file := range files {
		o.files[file.Name] = file
	}
	if file, ok := o.files["stack.xml"]; ok {
		f, err := file.Open()
		if err != nil {
			return nil, err
		}
		o.structure, err = processLayerStack(f)
		f.Close()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, ErrInvalidLayout
	}
	if _, ok := o.files["Thumbnails/thumbnail.png"]; !ok {
		return nil, ErrInvalidLayout
	}
	return o, nil
}

func (o *ORA) Bounds() image.Rectangle {
	return image.Rectangle{image.Point{}, image.Point{int(o.structure.Width), int(o.structure.Height)}}
}

func (o *ORA) Thumbnail() (image.Image, error) {
	f, err := o.files["thumbnail.png"].Open()
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
	merged, ok := o.files["merged.png"]
	if !ok {
		return nil, ErrNoMerged
	}
	f, err := merged.Open()
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

type Layer struct {
	layer
	offsetX, offsetY int
	f                *zip.File
}

func (o *ORA) Layer(name string) *Layer {
	l, x, y := o.structure.Get(name)
	ly, _ := l.(*layer)
	if l == nil {
		return nil
	}
	f, ok := o.files[ly.Src]
	if !ok {
		return nil
	}
	return &Layer{
		*ly,
		x, y,
		f,
	}
}

func (l *Layer) Image() (image.Image, error) {
	f, err := l.f.Open()
	if err != nil {
		return nil, err
	}
	i, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (l *Layer) Offsets() (int, int) {
	return l.offsetX, l.offsetY
}

// Errors
var (
	ErrInvalidLayout = errors.New("invalid layout")
	ErrNoMerged      = errors.New("ora does not include a merged image")
	ErrNoLayer       = errors.New("no layer with that name exists in ora")
	ErrLayerMissing  = errors.New("layer image missing")
)
