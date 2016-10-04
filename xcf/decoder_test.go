package xcf

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/MJKWoolnough/limage"
)

func TestConfigDecoder(t *testing.T) {
	return
	f, err := os.Open("test.xcf")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	c, _, err := image.DecodeConfig(f)
	f.Close()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	fmt.Println(c)
}

func TestDecoder(t *testing.T) {
	return
	f, err := os.Open("test.xcf")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	i, err := Decode(f)
	f.Close()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	fmt.Println(float64(255-i.Transparency)/2.55, "%")
	printGroup(&i.Group, "")
	f, err = os.Create("all.png")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	png.Encode(f, i)
	f.Close()
}

func printGroup(g *limage.Image, indent string) {
	fmt.Println(indent, g.Width, "x", g.Height, " - ")
	indent += "	"
	for _, l := range g.Layers {
		fmt.Print(indent, l.Name, " - ", float64(255-l.Transparency)/2.55, "% - ", l.Mode, " - ")
		/*
			f, err := os.Create(l.Name + ".png")
			if err != nil {
				return
			}
			png.Encode(f, l.Image)
			f.Close()
		*/
		switch i := l.Image.(type) {
		case *limage.Group:
			fmt.Println("Group")
			printGroup(i, indent)
			fmt.Print(indent, "Offset")
		case *limage.Text:
			fmt.Print("Text - ", i.String())
		case *limage.MaskedImage:
			fmt.Print("Masked Image")
		default:
			fmt.Print("Image")
		}
		fmt.Println(" - +", l.OffsetX, "+", l.OffsetY)
	}
}
