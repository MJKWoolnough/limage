package ora

import (
	"encoding/xml"
	"errors"
	"io"
	"strconv"
)

// Constants for the various composition operations
const (
	CompositeSrcOver composite = iota
	CompositeMultiply
	CompositeScreen
	CompositeOverlay
	CompositeDarken
	CompositeLighten
	CompositeColorDodge
	CompositeColorBurn
	CompositeHardLight
	CompositeSoftLight
	CompositeDifference
	CompositeColor
	CompositeLuminosity
	CompositeHue
	CompositeSaturation
	CompositePlus
	CompositeDstIn
	CompositeDstOut
	CompositeSrcAtop
	CompositeDstAtop
)

type composite int

type imageStack struct {
	Width, Height, XRes, YRes uint
	Stack                     []stackItem
}

func (i *imageStack) Get(name string) (stackItem, int, int) {
	s := stack{
		Items: i.Stack,
	}
	return s.Get(name)
}

type stackItem interface {
	Process(*xml.Decoder) error
	ProcessAttrs([]xml.Attr) error
	Name() string
	Offsets() (int, int)
}

type props struct {
	X, Y      int
	name      string
	Opacity   float32
	Invisible bool
	Composite composite
}

func (p props) Name() string {
	return p.name
}

func (p *props) ProcessAttrs(attrs []xml.Attr) error {
	p.Opacity = 1
	for _, a := range attrs {
		switch a.Name.Local {
		case "x":
			v, err := strconv.Atoi(a.Value)
			if err != nil {
				return err
			}
			p.X = v
		case "y":
			v, err := strconv.Atoi(a.Value)
			if err != nil {
				return err
			}
			p.Y = v
		case "name":
			p.name = a.Value
		case "opacity":
			v, err := strconv.ParseFloat(a.Value, 32)
			if err != nil {
				return err
			}
			if v < 0 || v > 1 {
				return ErrInvalidOpacity
			}
			p.Opacity = float32(v)
		case "visibility":
			switch a.Value {
			case "visible":
			case "invisible":
				p.Invisible = true
			default:
				return ErrInvalidVisibility
			}
		case "composite-op":
			switch a.Value {
			case "svg:src-over":
				p.Composite = CompositeSrcOver
			case "svg:multiply":
				p.Composite = CompositeMultiply
			case "svg:screen":
				p.Composite = CompositeScreen
			case "svg:overlay":
				p.Composite = CompositeOverlay
			case "svg:darken":
				p.Composite = CompositeDarken
			case "svg:lighten":
				p.Composite = CompositeLighten
			case "svg:color-dodge":
				p.Composite = CompositeColorDodge
			case "svg:color-burn":
				p.Composite = CompositeColorBurn
			case "svg:hard-light":
				p.Composite = CompositeHardLight
			case "svg:soft-light":
				p.Composite = CompositeSoftLight
			case "svg:difference":
				p.Composite = CompositeDifference
			case "svg:color":
				p.Composite = CompositeColor
			case "svg:luminosity":
				p.Composite = CompositeLuminosity
			case "svg:hue":
				p.Composite = CompositeHue
			case "svg:saturation":
				p.Composite = CompositeSaturation
			case "svg:plus":
				p.Composite = CompositePlus
			case "svg:dst-in":
				p.Composite = CompositeDstIn
			case "svg:dst-out":
				p.Composite = CompositeDstOut
			case "svg:src-atop":
				p.Composite = CompositeSrcAtop
			case "svg:dst-atop":
				p.Composite = CompositeDstAtop
			default:
				return ErrInvalidCompositeOp
			}
		}
	}
	return nil
}

func (p *props) Offsets() (int, int) {
	return p.X, p.Y
}

type stack struct {
	props
	Items []stackItem
}

func (s *stack) Process(x *xml.Decoder) error {
	for {
		t, err := x.Token()
		if err != nil {
			return err
		}
		switch t := t.(type) {
		case xml.StartElement:
			var i stackItem
			switch t.Name.Local {
			case "stack":
				i = new(stack)
			case "layer":
				i = new(layer)
			case "text":
				i = new(text)
			default:
				err = x.Skip()
				continue
				if err != nil {
					return err
				}
			}
			err = i.ProcessAttrs(t.Attr)
			if err != nil {
				return err
			}
			err = i.Process(x)
			if err != nil {
				return err
			}
			s.Items = append(s.Items, i)
		case xml.EndElement:
			return nil
		}
	}
}

func (s *stack) Get(name string) (stackItem, int, int) {
	for _, i := range s.Items {
		if i.Name() == name {
			x, y := i.Offsets()
			return i, x, y
		}
		if st, ok := i.(*stack); ok {
			si, x, y := st.Get(name)
			if si != nil {
				return si, x + s.X, y + s.Y
			}
		}
	}
	return nil, 0, 0
}

type layer struct {
	props
	Src string
}

func (l *layer) Process(x *xml.Decoder) error {
	for {
		tk, err := x.Token()
		if err != nil {
			return err
		}
		switch tk.(type) {
		case xml.StartElement:
			err = x.Skip()
			if err != nil {
				return nil
			}
		case xml.EndElement:
			return nil
		}
	}
}

func (l *layer) ProcessAttrs(attrs []xml.Attr) error {
	if err := l.props.ProcessAttrs(attrs); err != nil {
		return err
	}
	for _, a := range attrs {
		if a.Name.Local == "src" {
			l.Src = a.Value
			return nil
		}
	}
	return ErrMissingSrc
}

type text struct {
	props
	Text string
}

func (t *text) Process(x *xml.Decoder) error {
	for {
		tk, err := x.Token()
		if err != nil {
			return err
		}
		switch tk := tk.(type) {
		case xml.StartElement:
			err = x.Skip()
			if err != nil {
				return nil
			}
		case xml.EndElement:
			return nil
		case xml.CharData:
			t.Text += string(tk)
		}
	}
}

func processLayerStack(r io.Reader) (*imageStack, error) {
	x := xml.NewDecoder(r)
	i := &imageStack{
		XRes: 72,
		YRes: 72,
	}
	for {
		t, err := x.Token()
		if err != nil {
			return nil, err
		}
		if t, ok := t.(xml.StartElement); ok {
			if t.Name.Local != "image" {
				return nil, ErrInvalidLayerStack
			}
			for _, a := range t.Attr {
				switch a.Name.Local {
				case "w":
					v, err := strconv.ParseUint(a.Value, 10, 0)
					if err != nil {
						return nil, err
					}
					i.Width = uint(v)
				case "h":
					v, err := strconv.ParseUint(a.Value, 10, 0)
					if err != nil {
						return nil, err
					}
					i.Height = uint(v)
				case "xres":
					v, err := strconv.ParseUint(a.Value, 10, 0)
					if err != nil {
						return nil, err
					}
					if v < 1 {
						return nil, ErrInvalidResolution
					}
					i.XRes = uint(v)
				case "yres":
					v, err := strconv.ParseUint(a.Value, 10, 0)
					if err != nil {
						return nil, err
					}
					if v < 1 {
						return nil, ErrInvalidResolution
					}
					i.YRes = uint(v)
				}
			}
			break
		}
	}
	for {
		t, err := x.Token()
		if err != nil {
			return nil, err
		}
		if se, ok := t.(xml.StartElement); ok {
			if se.Name.Local == "stack" {
				var s stack
				err = s.Process(x)
				if err != nil {
					return nil, err
				}
				i.Stack = s.Items
				return i, nil
			}
			x.Skip()
		}
	}
	return nil, ErrMissingStack
}

// Errors
var (
	ErrMissingSrc         = errors.New("layer missing required src attribute")
	ErrInvalidVisibility  = errors.New("invalid visibility attribute value")
	ErrInvalidOpacity     = errors.New("invalid opacity attribute value")
	ErrInvalidCompositeOp = errors.New("invalid or unknown composite-op attribute value")
	ErrInvalidResolution  = errors.New("invalid resolution attribute value")
	ErrInvalidLayerStack  = errors.New("invalid layer stack xml")
	ErrMissingStack       = errors.New("missing stack element in stack.xml")
)
