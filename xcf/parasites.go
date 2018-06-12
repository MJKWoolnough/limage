package xcf

import (
	"io"
	"strconv"
	"strings"

	"vimagination.zapto.org/errors"
	"vimagination.zapto.org/parser"
)

const (
	//iccProfileParasiteName = "icc-profile"
	//commentParasiteName    = "gimp-comment"
	textParasiteName = "gimp-text-layer"
)

type parasite struct {
	name  string
	flags uint32
	data  []byte
}

type parasites []parasite

func (p parasites) Get(name string) *parasite {
	for n := range p {
		if p[n].name == name {
			return &p[n]
		}
	}
	return nil
}

func (d *reader) ReadParasites(l uint32) parasites {
	ps := make(parasites, 0, 32)
	for l > 0 {
		var p parasite
		p.name = d.ReadString()
		p.flags = d.ReadUint32()
		pplength := d.ReadUint32()
		read := 4 + uint32(len(p.name)) + 1 // length (uint32) + string([]byte) + \0 (byte)
		read += 4                           // flags
		read += 4                           // pplength
		read += pplength                    // len(data)
		if read > l {
			d.SetError(ErrInvalidParasites)
			return nil
		}
		l -= read
		p.data = make([]byte, pplength)
		d.Read(p.data)

		ps = append(ps, p)
	}
	return ps
}

func (d *reader) ReadParasite() parasite {
	var p parasite

	p.name = d.ReadString()
	p.flags = d.ReadUint32()
	pplength := d.ReadUint32()

	p.data = make([]byte, pplength)
	d.Read(p.data)

	return p
}

func (ps *parasite) Parse() ([]tag, error) {
	p := parser.New(parser.NewByteTokeniser(ps.data))
	p.TokeniserState(openTK)
	tags := make([]tag, 0, 32)
	for {
		tag, err := readTag(&p)
		if err != nil {
			if p.Err != io.EOF {
				return nil, err
			}
			break
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

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

// Tag represents a single tag from a parsed Parasite
type tag struct {
	Name   string
	Values []interface{}
}

func readTag(p *parser.Parser) (tag, error) {
	if p.Accept(parser.TokenDone) {
		return tag{}, io.EOF
	}
	if !p.Accept(tokenOpen) {
		return tag{}, ErrNoOpen
	}
	p.Get()
	if !p.Accept(tokenName) {
		return tag{}, ErrNoName
	}
	nt := p.Get()
	var tg tag
	tg.Name = nt[0].Data
	for {
		tt := p.AcceptRun(tokenValueString, tokenValueNumber)
		for _, v := range p.Get() {
			switch v.Type {
			case tokenValueString:
				tg.Values = append(tg.Values, v.Data)
			case tokenValueNumber:
				num, err := strconv.ParseFloat(v.Data, 64)
				if err != nil {
					return tag{}, err
				}
				tg.Values = append(tg.Values, num)
			}
		}
		switch tt {
		case tokenClose:
			p.Accept(tokenClose)
			p.Get()
			return tg, nil
		case tokenOpen:
			ttg, err := readTag(p)
			p.TokeniserState(valueTK)
			if err != nil {
				return tag{}, err
			}
			tg.Values = append(tg.Values, ttg)
		case parser.TokenDone:
			return tag{}, io.EOF
		default:
			return tag{}, ErrInvalidParasites
		}
	}
}

func openTK(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	t.AcceptRun(whitespace)
	switch t.Peek() {
	case -1, 0:
		return t.Done()
	}
	if !t.Accept(open) {
		t.Err = ErrInvalidParasites
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
	}, openTK
}

func nameTK(t *parser.Tokeniser) (parser.Token, parser.TokenFunc) {
	if !t.Accept(valName) {
		t.Err = ErrInvalidParasites
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
	} else if c < 0 {
		t.Err = ErrInvalidParasites
		return t.Error()
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
const (
	ErrInvalidParasites errors.Error = "invalid parasites layout"
	ErrNoOpen           errors.Error = "didn't receive Open token"
	ErrNoName           errors.Error = "didn't receive Name token"
)
