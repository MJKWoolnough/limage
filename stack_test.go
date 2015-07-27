package ora

import (
	"reflect"
	"strings"
	"testing"
)

func TestProcessStack(t *testing.T) {
	tests := []struct {
		xml   string
		image *imageStack
	}{
		{
			`<?xml version='1.0' encoding='UTF-8'?>
<image h="500" w="600">
        <stack>
                <layer composite-op="svg:src-over" name="New Layer" opacity="1.0" src="data/000.png" visibility="visible" x="0" y="0" />
                <layer composite-op="svg:src-atop" name="@$%^*£" opacity="0.5" src="data/001.png" visibility="invisible" x="-12" y="2" />
        </stack>
</image>`,
			&imageStack{
				Width:  600,
				Height: 500,
				XRes:   72,
				YRes:   72,
				Stack: []stackItem{
					&layer{
						props{
							X:         0,
							Y:         0,
							Name:      "New Layer",
							Opacity:   1,
							Composite: CompositeSrcOver,
						},
						"data/000.png",
					},
					&layer{
						props{
							X:         -12,
							Y:         2,
							Name:      "@$%^*£",
							Opacity:   0.5,
							Invisible: true,
							Composite: CompositeSrcAtop,
						},
						"data/001.png",
					},
				},
			},
		},
	}

	for n, test := range tests {
		i, err := processLayerStack(strings.NewReader(test.xml))
		if err != nil {
			t.Errorf("test %d: unexpected error: %s", n+1, err)
		} else if !reflect.DeepEqual(i, test.image) {
			t.Errorf("test %d: output not as expected", n+1)
		}
	}
}
