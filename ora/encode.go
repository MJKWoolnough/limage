package ora

import (
	"archive/zip"
	"encoding/xml"
	"image"
	"image/png"
	"io"
	"strconv"

	"github.com/MJKWoolnough/limage"
)

const mimetypeStr = "image/openraster"

type children []interface{}

type stack struct {
	XMLName     struct{} `xml:"stack"`
	X           uint     `xml:"x,attrib,omitempty"`
	Y           uint     `xml:"y,attrib,omitempty"`
	Name        string   `xml:"name,attrib,omitempty"`
	Opacity     float64  `xml:"opacity,attrib"`
	Visibility  string   `xml:"visibility,attrib,omitempty"`
	CompositeOp string   `xml:"composite-op,attrib,omitempty"`
	children
}

type layer struct {
	XMLName     struct{} `xml:"layer"`
	X           uint     `xml:"x,attrib,omitempty"`
	Y           uint     `xml:"y,attrib,omitempty"`
	Name        string   `xml:"name,attrib,omitempty"`
	Opacity     float64  `xml:"opacity,attrib"`
	Visibility  string   `xml:"visibility,attrib,omitempty"`
	CompositeOp string   `xml:"composite-op,attrib,omitempty"`
	Source      string   `xml:"src,attrib"`
}

func Encode(w io.Writer, m image.Image) error {
	var lim limage.Image
	switch im := m.(type) {
	case limage.Image:
		lim = im
	case *limage.Image:
		lim = *im
	case limage.Layer:
		lim = limage.Image{im}
	case *limage.Layer:
		lim = limage.Image{*im}
	default:
		lim = limage.Image{
			limage.Layer{
				LayerBounds: m.Bounds(),
				Image:       m,
			},
		}
	}

	zw := zip.NewWriter(w)
	defer e.Close()

	// Write MIME
	fw, err := zw.Create("mimetype")
	if err != nil {
		return err
	}
	_, err = fw.Write([]byte(mimetypeStr))
	if err != nil {
		return err
	}

	// Write Layer images
	stack, err := e.WriteLayers(lim)
	if err != nil {
		return err
	}

	// Write Stack
	b := lim.Bounds()
	e.xml.Image.Width = b.Dx()
	e.xml.Image.Height = b.Dy()
	fw, err = zw.Create("stack.xml")
	if err != nil {
		return err
	}
	err = xml.NewEncoder(fw).EncodeElement(stack, xml.StartElement{
		Name: xml.Name{
			Local: "image",
		},
		Attr: []xml.Attr{
			{
				Name:  xml.Name{Local: "version"},
				Value: "0.0.3",
			},
			{
				Name:  xml.Name{Local: "w"},
				Value: strconv.Itoa(b.Dx()),
			},
			{
				Name:  xml.Name{Local: "h"},
				Value: strconv.Itoa(b.Dy()),
			},
		},
	})
	if err != nil {
		return err
	}

	// Write Merged Image
	fw, err = zw.Create("mergedimage.png")
	if err != nil {
		return err
	}
	err = png.Encode(fw, m)
	if err != nil {
		return err
	}

	// Write Thumbnail
	fw, err = zw.Create("Thumbnails/thumbnail.png")
	if err != nil {
		return err
	}

	if w, h := b.Dx(), b.Dy(); w > 256 || h > 256 {
		var scale float64
		if w > h {
			scale = float64(w) / 256
		} else {
			scale = float64(h) / 256
		}
		m = thumbnail{Image: m, scale: scale}
	}

	err = png.Encode(fw, m)
	if err != nil {
		return err
	}

	return nil
}
