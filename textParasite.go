package xcf

import (
	"errors"
	"image/color"
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
	FontColor, ForeColor, BackColor        color.Color
	Size, LetterSpacing, Rise              uint
	Bold, Italic, Underline, Strikethrough bool
	Data                                   string
}

func parseTextParasite(data []byte) (TextData, error) {
	p := parser.New(parser.NewByteTokeniser(data))
	p.TokeniserState(openTK)
	for {
		readTag(&p)
	}
	return TextData{}, nil
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
	switch string(c) {
	case open:
		return openTK(t)
	case close:
		return parser.Token{
			Type: tokenClose,
		}, valueTK
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
