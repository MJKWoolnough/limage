package ora

import (
	"archive/zip"
	"fmt"
	"image"
	"image/png"
	"io"
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

type encoder struct {
	*zip.Writer
	xml struct {
		Image struct {
			Version string   `xml:"version,attr"`
			Width   uint     `xml:"w,attr"`
			Height  uint     `xml:"h,attr"`
			Stack   children `xml:"stack"`
		} `xml:"image"`
	}
}

func Encode(w io.Writer, m image.Image) error {
	zw := zip.NewWriter(w)
	defer zw.Close()
	fw, err := zw.Create("mimetype")
	if err != nil {
		return err
	}
	_, err = fw.Write([]byte(mimetypeStr))
	if err != nil {
		return err
	}
	fw, err = zw.Create("stack.xml")
	if err != nil {
		return err
	}
	b := m.Bounds()
	fmt.Fprintf(fw, "<?xml version='1.0' encoding='UTF-8'?>\n<image w=\"%d\" h=\"%d\"><stack><layer composite-op=\"svg:src-over\" name=\"Layer\" opacity=\"1.0\" src=\"data/layer.png\" visibility=\"visible\" x=\"0\" y=\"0\" /></stack></image>", b.Dx(), b.Dy())
	fw, err = zw.Create("data/layer.png")
	if err != nil {
		return err
	}
	err = png.Encode(fw, m)
	if err != nil {
		return err
	}
	fw, err = zw.Create("mergedimage.png")
	if err != nil {
		return err
	}
	err = png.Encode(fw, m)
	if err != nil {
		return err
	}

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
