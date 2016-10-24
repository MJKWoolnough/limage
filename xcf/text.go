package xcf

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html"
	"image/color"
	"io"
	"strconv"
	"strings"

	"github.com/MJKWoolnough/limage"
	"github.com/MJKWoolnough/limage/lcolor"
)

func parseTextData(t *parasite) (limage.TextData, error) {
	tags, err := t.Parse()
	if err != nil {
		return nil, err
	}
	var (
		textData    string
		defaultText limage.TextDatum
	)
	defaultText.BackColor = color.Alpha{}
	defaultText.ForeColor = color.Gray{}
	for _, tag := range tags {
		switch tag.Name {
		case "text":
			defaultText.Data, _ = tag.Values[0].(string)
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
				f, _ := tag.Values[0].(float64)
				defaultText.Size = uint32(f)
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
					defaultText.ForeColor = lcolor.RGB{uint8(r), uint8(g), uint8(b)}
				}
			}
		case "justify":
		case "box-mode":
		case "box-width":
		case "box-height":
		case "box-unit":
		case "hinting":
		}
	}
	if defaultText.Data != "" {
		return limage.TextData{defaultText}, nil
	}
	xd := xml.NewDecoder(strings.NewReader(textData))
	stack := limage.TextData{defaultText}
	td := make(limage.TextData, 0, 32)
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
						s, err := strconv.ParseUint(a.Value, 10, 32)
						if err != nil {
							return nil, err
						}
						nt.Size = uint32(s) >> 10
					case "letter_spacing":
						ls, err := strconv.ParseUint(a.Value, 10, 32)
						if err != nil {
							return nil, err
						}
						nt.LetterSpacing = uint32(ls) >> 10
					case "rise":
						r, err := strconv.ParseUint(a.Value, 10, 32)
						if err != nil {
							return nil, err
						}
						nt.Rise = uint32(r) >> 10
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

func (e *encoder) WriteText(text limage.TextData, dx, dy uint32) {
	e.WriteUint32(propTextLayerFlags)
	e.WriteUint32(4)
	e.WriteUint32(1)
	e.WriteUint32(21)

	var (
		buf  bytes.Buffer
		base limage.TextDatum
	)

	if len(text) == 1 {
		base = text[0]
		fmt.Fprintf(&buf, "(text %q)\n", base.Data)
	} else {
		base = limage.TextDatum{
			BackColor: lcolor.RGB{0, 0, 0},
			ForeColor: lcolor.RGB{0, 0, 0},
			Font:      "Sans",
			Size:      18,
		}

		buf.WriteString("(markup \"<markup>")

		for _, td := range text {
			if td.Font != "Sans" {
				fmt.Fprintf(&buf, "<span font=%q>", strconv.Quote(td.Font))
			}
			if td.Bold {
				buf.WriteString("<b>")
			}
			if td.Italic {
				buf.WriteString("<i>")
			}
			if td.Underline {
				buf.WriteString("<u>")
			}
			if td.Strikethrough {
				buf.WriteString("<s>")
			}
			if td.LetterSpacing != 0 {
				fmt.Fprintf(&buf, "<span letter_spacing=\"%d\">", td.LetterSpacing<<10)
			}
			if td.Size != 18 {
				fmt.Fprintf(&buf, "<span size=\"%d\">", td.Size<<10)
			}
			if td.Rise != 0 {
				fmt.Fprintf(&buf, "<span rise=\"%d\">", td.Rise<<10)
			}
			buf.WriteString(strconv.Quote(html.EscapeString(td.Data)))
			if td.Rise != 0 {
				buf.WriteString("</span>")
			}
			if td.Size != 18 {
				buf.WriteString("</span>")
			}
			if td.LetterSpacing != 0 {
				buf.WriteString("</span>")
			}
			if td.Strikethrough {
				buf.WriteString("</s>")
			}
			if td.Underline {
				buf.WriteString("</u>")
			}
			if td.Italic {
				buf.WriteString("</i>")
			}
			if td.Bold {
				buf.WriteString("</b>")
			}
			if td.Font != "Sans" {
				buf.WriteString("</span>")
			}
		}

		buf.WriteString("</markup>\")\n")
	}

	r, g, b, _ := base.ForeColor.RGBA()

	fmt.Fprintf(&buf, "(font %q)\n"+
		"(font-size %.9f)\n"+
		"(font-size-units pixels)\n"+
		"(antialias yes)\n"+
		"(base-direction ltr)\n"+
		"(color (color-rgb %.6f %.6f %.6f))\n"+
		"(justify left)\n"+
		"(box-mode dynamic)\n"+
		"(box-width %.6f)\n"+
		"(box-height %.6f)\n"+
		"(box-unit pixels)\n"+
		"(hinting yes)\n"+
		"\x00", base.Font, base.Size, float32(r>>8), float32(g>>8), float32(b>>8), float64(dx), float64(dy))

	// write base

	data := buf.Bytes()

	e.WriteUint32(uint32(4 + len(textParasiteName) + 1 + 4 + 4 + len(data)))

	e.WriteString(textParasiteName)
	e.WriteUint32(0)
	e.WriteUint32(uint32(len(data)))
	e.Write(data)
}
