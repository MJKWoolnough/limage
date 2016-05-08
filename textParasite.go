package xcf

import (
	"errors"
	"image/color"
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

type tokeniser struct {
	parser.Parser
}

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
	t := tokeniser{
		Parser: parser.NewByteParser(data),
	}
	t.State = t.open

	return TextData{}, nil
}

func (t *tokeniser) readTag() (string, []interface{}) {
}

func (t *tokeniser) open() (parser.Token, parser.StateFn) {
	t.AcceptRun(whitespace)
	if !t.Accept(open) {
		// error
	}
	t.Get()
	return parser.Token{
		Type: tokenOpen,
	}, t.name
}

func (t *tokeniser) name() (parser.Token, parser.StateFn) {
	t.AcceptRun(valName)
	return parser.Token{
		Type: tokenName,
		Data: t.Get(),
	}, t.value
}

func (t *tokeniser) value() (parser.Token, parser.StateFn) {
	t.AcceptRun(whitespace)
	t.Get()
	switch c := t.Peek(); string(c) {
	case open:
		return t.open()
	case close:
		return parser.Token{
			Type: tokenClose,
		}, t.value
	case quoted:
		return parser.Token{
			Type: tokenValueString,
			Data: t.quotedString(),
		}, t.value
	}
	if strings.ContainsRune(digit, c) {
		t.AcceptRun(digit)
		t.Accept(".")
		t.AcceptRun(digit)
		return parser.Token{
			Type: tokenValueNumber,
			Data: t.Get(),
		}, t.value
	}
	t.AcceptRun(valName)
	return parser.Token{
		Type: tokenValueString,
		Data: t.Get(),
	}, t.value
}

func (t *tokeniser) quotedString() string {
	t.Accept(quoted)
	t.Get()
	var s string
	for {
		t.ExceptRun(quoted + "\\")
		s += t.Get()
		if t.Accept("\\") {
			switch c := t.Peek(); c {
			case "\"", "\\":
				s += string(c)
			default:
				s += "\\" + string(c)
			}
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
