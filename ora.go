// Package ora is an implementation of an OpenRaster decoder
package ora

import (
	"archive/zip"
	"errors"
	"image"
	_ "image/png" // PNG format is required for ORA, at least in the thumbnails
	"io"
)

// ReadCloserAt is an interface combining both io.ReaderAt and io.Closer
type ReadCloserAt interface {
	io.ReaderAt
	io.Closer
}

// ORA is a type which contains the description of the ORA file
type ORA struct {
	io.Closer
	structure *imageStack
	files     map[string]*zip.File
}

// OpenFile takes a filename and returns a new ORA or an error if one occured
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

// Open takes a ReadCloserAt and returns a new ORA or an error if one occured
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

// Bounds returns a rectangle describing the width and height of the image
func (o *ORA) Bounds() image.Rectangle {
	return image.Rectangle{image.Point{}, image.Point{int(o.structure.Width), int(o.structure.Height)}}
}

// Thumbnail returns the embedded thumbnail image.
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

// Merged returns an image containing all of the layers merged into a single
// image
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

// Layer represents a single layer of the ORA image
type Layer struct {
	layer
	offsetX, offsetY int
	f                *zip.File
}

// Layer returns a layer in the image corresponding with the given name.
// If there is no layer with the given name, nil is returned
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

// Image returns the layer as an editable image
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

// Offsets returns the layers offset with regards to the main canvas
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
