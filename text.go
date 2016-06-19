package xcf

import (
	"encoding/xml"
	"image/color"
	"io"
	"strconv"
	"strings"
)

func parseTextData(t parasite) (TextData, error) {
	tags, err := t.Parse()
	if err != nil {
		return nil, err
	}
	var (
		textData    string
		defaultText TextDatum
	)
	defaultText.BackColor = color.Alpha{}
	defaultText.ForeColor = color.Gray{}
	for _, tag := range tags {
		switch tag.Name {
		case "markup":
			if len(tag.Values) == 1 {
				textData, _ = tag.Values[0].(string)
			}
		case "font":
			if len(tag.Values) == 1 {
				defaultText.Font, _ = tag.Values[0].(string)
			}
		case "font-size":
			if len(tag.Values) == 1 {
				defaultText.Size, _ = tag.Values[0].(float64)
			}
		case "font-size-unit":
		case "antialias":
		case "language":
		case "base-direction":
		case "color":
			if len(tag.Values) == 1 {
				t, _ := tag.Values[0].(Tag)
				if t.Name == "color-rgb" && len(t.Values) != 3 {
					r, _ := t.Values[0].(float64)
					g, _ := t.Values[1].(float64)
					b, _ := t.Values[2].(float64)
					defaultText.ForeColor = rgb{uint8(r), uint8(g), uint8(b)}
				}
			}
		case "justify":
			//		case "box-mode":
			//		case "box-width":
			//		case "box-height":
			//		case "box-unit":
		case "hinting":
		}
	}
	xd := xml.NewDecoder(strings.NewReader(textData))
	stack := TextData{defaultText}
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
			switch t.Name.Local {
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
