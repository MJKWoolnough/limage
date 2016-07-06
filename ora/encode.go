package ora

import (
	"archive/zip"
	"fmt"
	"image"
	"image/png"
	"io"
)

const mimetypeStr = "image/openraster"

type subimage interface {
	SubImage(image.Rectangle) image.Image
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
	b := m.Bounds().Max
	fmt.Fprintf(fw, "<?xml version='1.0' encoding='UTF-8'?>\n<image w=\"%d\" h=\"%d\"><stack><layer composite-op=\"svg:src-over\" name=\"Layer\" opacity=\"1.0\" src=\"data/layer.png\" visibility=\"visible\" x=\"0\" y=\"0\" /></stack></image>", b.X, b.Y)
	fw, err = zw.Create("data/layer.png")
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

	// TODO: Create an actual thumbnail
	err = png.Encode(fw, m.(subimage).SubImage(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: 256, Y: 256},
	}))
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
	return nil
}
