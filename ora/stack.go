package ora

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"image/png"
	"strconv"

	"github.com/MJKWoolnough/limage"
)

type decoder struct {
	zr *zip.Reader
	x  *xml.Decoder
}

func (d decoder) readStack() (limage.Image, error) {
	i := make(limage.Image, 32)
Loop:
	for {
		t, err := d.x.Token()
		if err != nil {
			return nil, err
		}
		switch t := t.(type) {
		case *xml.StartElement:
			switch t.Name {
			case "stack", "layer":
				l, err := d.readLayer(t)
				if err != nil {
					return nil, err
				}
				i = append(i, l)
			default:
				if err := d.skipTag(); err != nil {
					return nil, err
				}
			}
		case *xml.EndElement:
			break Loop
		}
	}
	if len(i) != cap(i) {
		j := make(limage.Image, len(i))
		copy(j, i)
		i = j
	}
	return i, nil
}

func (d decoder) readLayer(s *xml.StartElement) (limage.Layer, error) {
	var (
		l      limage.Layer
		source string
	)
	for _, a := range s.Attr {
		switch a.Name.Local {
		case "name":
			l.Name = a.Value
		case "x":
			offset, err := strconv.Atoi(a.Value)
			if err != nil {
				return l, err
			}
			l.LayerBounds.Min.X = offset
		case "y":
			offset, err := strconv.Atoi(a.Value)
			if err != nil {
				return l, err
			}
			l.LayerBounds.Min.Y = offset
		case "opacity":
			o, err := strconv.ParseFloat(a.Value, 64)
			if err != nil {
				return l, err
			}
			l.Transparency = uint8(255 * (1 - o))
		case "visibility":
			l.Invisible = a.Value == "hidden"
		case "composite-op":
			switch a.Value {
			case "svg:src-over":
				l.Mode = limage.CompositeNormal
			case "svg:multiply":
				l.Mode = limage.CompositeMultiply
			case "svg:screen":
				l.Mode = limage.CompositeScreen
			case "svg:overlay":
				l.Mode = limage.CompositeOverlay
			case "svg:darken":
				l.Mode = limage.CompositeDarkenOnly
			case "svg:lighten":
				l.Mode = limage.CompositeLightenOnly
			case "svg:color-dodge":
				l.Mode = limage.CompositeDodge
			case "svg:color-burn":
				l.Mode = limage.CompositeBurn
			case "svg:hard-light":
				l.Mode = limage.CompositeHardLight
			case "svg:soft-light":
				l.Mode = limage.CompositeSoftLight
			case "svg:difference":
				l.Mode = limage.CompositeDifference
			case "svg:color":
				l.Mode = limage.CompositeColor
			case "svg:luminosity":
				l.Mode = limage.CompositeLuminosity
			case "svg:hue":
				l.Mode = limage.CompositeHue
			case "svg:saturation":
				l.Mode = limage.CompositeSaturation
			case "svg:plus":
				l.Mode = limage.CompositePlus
			case "svg:dst-in":
				l.Mode = limage.CompositeDestinationIn
			case "svg:dst-out":
				l.Mode = limage.CompositeDestinationOut
			case "svg:src-atop":
				l.Mode = limage.CompositeSourceAtop
			case "svg:dst-atop":
				l.Mode = limage.CompositeDestinationAtop
			}
		case "src":
			source = a.Value
		}
	}
	if s.Name == "stack" {
		var err error
		l.Image, err = d.readStack()
		return l, err
	}
	for _, f := range d.zr.File {
		if f.Name == source {
			fr, err := f.Open()
			if err != nil {
				return l, err
			}
			l.Image, err = png.Decode(fr)
			if err != nil {
				return l, err
			}
			fr.Close()
			break
		}
	}
	if l.Image == nil {
		return l, ErrInvalidSource
	}
	if err := d.skipTag(); err != nil {
		return nil, err
	}
	return l, nil
}

func (d decoder) skipTag() error {
	toSkip := 0
	for {
		t, err := d.x.Token()
		if err != nil {
			return err
		}
		switch t.(type) {
		case *xml.StartElement:
			toSkip++
		case *xml.EndElement:
			if toSkip == 0 {
				return nil
			}
			toSkip--
		}
	}
}

// Errors
var (
	ErrInvalidSource = errors.New("invalid source")
)
