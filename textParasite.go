package xcf

import (
	"encoding/xml"
	"errors"
	"image/color"
	"io"
	"strconv"
	"strings"

	"github.com/MJKWoolnough/parser"
)

const (
	open       = "("
	close      = ")"
	chars      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	valName    = chars + "-"
	digit      = "1234567890"
	quoted     = "\""
	whitespace = "\n\r "
)

const (
	tokenOpen parser.TokenType = iota
	tokenClose
	tokenName
	tokenValueString
	tokenValueNumber
)

type TextData []Text

func (t TextData) String() string {
	var s string
	for _, d := range t {
		s += d.Data
	}
	return s
}

type Text struct {
	ForeColor, BackColor                   color.Color
	Size, LetterSpacing, Rise              float64
	Bold, Italic, Underline, Strikethrough bool
	Font, Data                             string
	FontUnit                               uint8
}

func parseTextParasite(data []byte) (TextData, error) {
	p := parser.New(parser.NewByteTokeniser(data))
	p.TokeniserState(openTK)
	var (
		markup      string
		defaultText Text
	)
	for {
		t, err := readTag(&p)
		if err != nil {
			if p.Err != io.EOF {
				return nil, err
			}
			break
		}
		switch t.name {
		case "markup":
			if len(t.values) != 1 {
				// Error
			}
			str, ok := t.values[0].(string)
			if !ok {
				// Error
			}
			markup = str
		case "font":
			if len(t.values) == 1 {
				defaultText.Font, _ = t.values[0].(string)
			}
		case "font-size":
			if len(t.values) == 1 {
				defaultText.Size, _ = t.values[0].(float64)
			}
		case "font-size-unit":
		case "antialias":
		case "language":
		case "base-direction":
		case "color":
			if len(t.values) == 1 {
				defaultText.ForeColor, _ = t.values[0].(color.Color)
			}
		case "justify":
			//		case "box-mode":
			//		case "box-width":
			//		case "box-height":
			//		case "box-unit":
		case "hinting":
		}
	}
	xd := xml.NewDecoder(strings.NewReader(markup))
	stack := []Text{defaultText}
	td := make(TextData, 0, 32)
	for {
		t, err := xd.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		switch t := t.(type) {
		case xml.StartElement:
			nt := stack[len(stack)-1]
			switch t.Name.Space {
			case "markup":
			case "span":
				for _, a := range t.Attr {
					switch a.Name.Local {
					case "font":
						nt.Font = a.Value
					case "foreground":
						if len(a.Value) == 7 && a.Value[0] == '#' {
							n, err := strconv.ParseUint(a.Value[1:], 16, 32)
							if err != nil {
								return nil, err
							}
							nt.ForeColor = color.RGBA{uint8(n >> 16), uint8((n >> 8) & 255), uint8(n & 255), 255}
						} else if len(a.Value) == 4 && a.Value[0] == '#' {
							n, err := strconv.ParseUint(a.Value[1:], 16, 32)
							if err != nil {
								return nil, err
							}
							r := (n >> 4) & 240
							r |= r >> 4
							g := n & 240
							g |= g >> 4
							b := n & 15
							b |= b << 4
							nt.ForeColor = color.RGBA{uint8(r), uint8(g), uint8(b), 255}
						}
					case "size":
						nt.Size, err = strconv.ParseFloat(a.Value, 64)
						if err != nil {
							return nil, err
						}
					case "letter_spacing":
						nt.LetterSpacing, err = strconv.ParseFloat(a.Value, 64)
						if err != nil {
							return nil, err
						}
					case "rise":
						nt.Rise, err = strconv.ParseFloat(a.Value, 64)
						if err != nil {
							return nil, err
						}
					}
				}
			case "b":
				nt.Bold = true
			case "i":
				nt.Italic = true
			case "s":
				nt.Strikethrough = true
			case "u":
				nt.Underline = true
			}
			stack = append(stack, nt)
		case xml.CharData:
			nt := stack[len(stack)-1]
			nt.Data = string(t)
			td = append(td, nt)
		case xml.EndElement:
			stack = stack[:len(stack)-1]
		}
	}
	return td, nil
}

type tag struct {
	name   string
	values []interface{}
}

func readTag(p *parser.Parser) (tag, error) {
	if !p.Accept(tokenOpen) {
		return tag{}, ErrInvalidLayout
	}
	p.Get()
	if !p.Accept(tokenName) {
		return tag{}, ErrInvalidLayout
	}
	nt := p.Get()
	var tg tag
	tg.name = nt[0].Data
	for {
		tt := p.AcceptRun(tokenValueString, tokenValueNumber)
		for _, v := range p.Get() {
			switch v.Type {
			case tokenValueString:
				tg.values = append(tg.values, v.Data)
			case tokenValueNumber:
				num, err := strconv.ParseFloat(v.Data, 64)
				if err != nil {
					return tag{}, err
				}
				tg.values = append(tg.values, num)
			}
		}
		switch tt {
		case tokenClose:
			p.Accept(tokenClose)
			p.Get()
			return tg, nil
		case tokenOpen:
			ttg, err := readTag(p)
			if err != nil {
				return tag{}, err
			}
			switch ttg.name {
			case "color-rgb":
				if len(ttg.values) == 3 {
					r, rok := ttg.values[0].(float64)
					g, gok := ttg.values[1].(float64)
					b, bok := ttg.values[2].(float64)
					if !rok || !gok || !bok {
						//error
					}
					tg.values = append(tg.values, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
				} else {
					//error??
				}
			}
		case parser.TokenDone:
			return tag{}, io.EOF
		default:
			return tag{}, ErrInvalidLayout
		}
	}
}

func openTK(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.AcceptRun(whitespace)
	if !t.Accept(open) {
		t.Err = ErrInvalidLayout
		return t.Error()
	}
	t.Get()
	return parser.Token{
		Type: tokenOpen,
	}, nameTK
}

func closeTK(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.Accept(close)
	t.Get()
	return parser.Token{
		Type: tokenClose,
	}, valueTK
}

func nameTK(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if !t.Accept(valName) {
		t.Err = ErrInvalidLayout
		return t.Error()
	}
	t.AcceptRun(valName)
	return parser.Token{
		Type: tokenName,
		Data: t.Get(),
	}, valueTK
}

func valueTK(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.AcceptRun(whitespace)
	t.Get()
	c := t.Peek()
	if c == 0 {
		return t.Done()
	}
	switch string(c) {
	case open:
		return openTK(t)
	case close:
		return closeTK(t)
	case quoted:
		return parser.Token{
			Type: tokenValueString,
			Data: quotedString(t),
		}, valueTK
	}
	if strings.ContainsRune(digit, c) {
		t.AcceptRun(digit)
		t.Accept(".")
		t.AcceptRun(digit)
		return parser.Token{
			Type: tokenValueNumber,
			Data: t.Get(),
		}, valueTK
	}
	t.AcceptRun(valName)
	return parser.Token{
		Type: tokenValueString,
		Data: t.Get(),
	}, valueTK
}

func quotedString(t *parser.Tokeniser) string {
	t.Accept(quoted)
	t.Get()
	var s string
	for {
		t.ExceptRun(quoted + "\\")
		s += t.Get()
		if t.Accept("\\") {
			c := string(t.Peek())
			switch c {
			case "\"", "\\":
				s += c
			default:
				s += "\\" + c
			}
			t.Accept(c)
			t.Get()
			continue
		}
		break
	}
	t.Accept(quoted)
	t.Get()
	return s
}

// Errors
var (
	ErrInvalidLayout = errors.New("invalid layout")
)
