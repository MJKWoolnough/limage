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

// Encode encodes the given image as an ORA file to the given Writer
func Encode(w io.Writer, m image.Image) error {
	var lim limage.Layer
	switch im := m.(type) {
	case limage.Layer:
		lim.Image = im.Image
	case *limage.Layer:
		lim.Image = im.Image
	default:
		lim = limage.Layer{
			Image: m,
		}
	}
	b := m.Bounds()
	lim.LayerBounds.Max.X = b.Dx()
	lim.LayerBounds.Max.Y = b.Dy()

	zw := zip.NewWriter(w)
	defer zw.Close()

	// Write MIME
	fw, err := zw.CreateHeader(&zip.FileHeader{
		Name:   "mimetype",
		Method: zip.Store,
	})
	if err != nil {
		return err
	}
	if _, err = fw.Write([]byte(mimetypeStr)); err != nil {
		return err
	}

	// Write Layer images
	if _, err = writeLayers(zw, lim, 0); err != nil {
		return err
	}

	// Write Stack
	if fw, err = zw.Create("stack.xml"); err != nil {
		return err
	}
	if _, err = fw.Write([]byte(xml.Header)); err != nil {
		return err
	}
	e := xml.NewEncoder(fw)
	e.Indent("", "	")
	if err = e.EncodeToken(xml.StartElement{
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
				Value: strconv.Itoa(lim.LayerBounds.Max.X),
			},
			{
				Name:  xml.Name{Local: "h"},
				Value: strconv.Itoa(lim.LayerBounds.Max.Y),
			},
		},
	}); err != nil {
		return err
	}

	if _, err = writeStack(e, lim, 0); err != nil {
		return err
	}

	if err = e.EncodeToken(xml.EndElement{
		Name: xml.Name{
			Local: "image",
		},
	}); err != nil {
		return err
	}

	if err = e.Flush(); err != nil {
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

	if lim.LayerBounds.Max.X > 256 || lim.LayerBounds.Max.Y > 256 {
		var scale float64
		if lim.LayerBounds.Max.X > lim.LayerBounds.Max.Y {
			scale = float64(lim.LayerBounds.Max.X) / 256
		} else {
			scale = float64(lim.LayerBounds.Max.Y) / 256
		}
		m = thumbnail{Image: m, scale: scale}
	}

	err = png.Encode(fw, m)
	if err != nil {
		return err
	}

	return nil
}

func writeLayers(zw *zip.Writer, lim limage.Layer, layerNum int) (int, error) {
	var (
		err error
		f   io.Writer
	)
	switch im := lim.Image.(type) {
	case limage.Image:
		layerNum, err = writeGroup(zw, im, layerNum)
	case *limage.Image:
		layerNum, err = writeGroup(zw, *im, layerNum)
	// case limage.Text, *limage.Text: // text is not yet in the spec
	default:
		layerNum++
		f, err = zw.Create("data/" + strconv.Itoa(layerNum) + ".png")
		if err != nil {
			return 0, err
		}
		err = png.Encode(f, lim.Image)
		if err != nil {
			return 0, err
		}
	}
	return layerNum, err
}

func writeGroup(zw *zip.Writer, lim limage.Image, layerNum int) (int, error) {
	var err error
	for _, l := range lim {
		layerNum, err = writeLayers(zw, l, layerNum)
		if err != nil {
			return 0, err
		}
	}
	return layerNum, nil
}

func writeStack(e *xml.Encoder, lim limage.Layer, layerNum int) (int, error) {
	attrs := make([]xml.Attr, 0, 7)
	b := lim.Bounds().Min
	if lim.Name != "" {
		attrs = append(attrs, xml.Attr{
			Name:  xml.Name{Local: "name"},
			Value: lim.Name,
		})
	}
	if b.X > 0 {
		attrs = append(attrs, xml.Attr{
			Name:  xml.Name{Local: "x"},
			Value: strconv.Itoa(b.X),
		})
	}
	if b.Y > 0 {
		attrs = append(attrs, xml.Attr{
			Name:  xml.Name{Local: "y"},
			Value: strconv.Itoa(b.Y),
		})
	}
	if lim.Invisible {
		attrs = append(attrs, xml.Attr{
			Name:  xml.Name{Local: "visibility"},
			Value: "hidden",
		})
	}
	if lim.Transparency != 0 {
		attrs = append(attrs, xml.Attr{
			Name:  xml.Name{Local: "visibility"},
			Value: strconv.FormatFloat(float64(255-lim.Transparency)/255, 'f', -1, 64),
		})
	}
	var op string
	switch lim.Mode {
	//case limage.CompositeNormal:
	//	op = "svg:src-over"
	case limage.CompositeMultiply:
		op = "svg:multiply"
	case limage.CompositeScreen:
		op = "svg:screen"
	case limage.CompositeOverlay:
		op = "svg:overlay"
	case limage.CompositeDarkenOnly:
		op = "svg:darken"
	case limage.CompositeLightenOnly:
		op = "svg:lighten"
	case limage.CompositeDodge:
		op = "svg:color-dodge"
	case limage.CompositeBurn:
		op = "svg:color-burn"
	case limage.CompositeHardLight:
		op = "svg:hard-light"
	case limage.CompositeSoftLight:
		op = "svg:soft-light"
	case limage.CompositeDifference:
		op = "svg:difference"
	case limage.CompositeColor:
		op = "svg:color"
	case limage.CompositeLuminosity:
		op = "svg:luminosity"
	case limage.CompositeHue:
		op = "svg:hue"
	case limage.CompositeSaturation:
		op = "svg:saturation"
	case limage.CompositePlus:
		op = "svg:plus"
	case limage.CompositeDestinationIn:
		op = "svg:dst-in"
	case limage.CompositeDestinationOut:
		op = "svg:dst-out"
	case limage.CompositeSourceAtop:
		op = "svg:src-atop"
	case limage.CompositeDestinationAtop:
		op = "svg:dst-atop"
	}
	if op != "" {
		attrs = append(attrs, xml.Attr{
			Name:  xml.Name{Local: "composite-op"},
			Value: op,
		})
	}
	var err error
	switch im := lim.Image.(type) {
	case limage.Image:
		layerNum, err = writeGroupStack(e, im, attrs, layerNum)
	case *limage.Image:
		layerNum, err = writeGroupStack(e, *im, attrs, layerNum)
	// case limage.Text, *limage.Text: // text is not yet in the spec
	default:
		layerNum++
		attrs = append(attrs, xml.Attr{
			Name:  xml.Name{Local: "src"},
			Value: "data/" + strconv.Itoa(layerNum) + ".png",
		})
		err = e.EncodeToken(xml.StartElement{
			Name: xml.Name{Local: "layer"},
			Attr: attrs,
		})
		if err == nil {
			err = e.EncodeToken(xml.EndElement{
				Name: xml.Name{Local: "layer"},
			})
		}
	}
	return layerNum, err
}

func writeGroupStack(e *xml.Encoder, lim limage.Image, attrs []xml.Attr, layerNum int) (int, error) {
	err := e.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "stack"},
		Attr: attrs,
	})
	if err != nil {
		return 0, err
	}
	for _, l := range lim {
		layerNum, err = writeStack(e, l, layerNum)
		if err != nil {
			return 0, err
		}
	}
	return layerNum, e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "stack"}})
}
