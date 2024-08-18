# limage
--
    import "vimagination.zapto.org/limage"

Package limage introduces structures to read and build layered images.

## Usage

#### type Composite

```go
type Composite uint32
```

Composite determines how two layers are composed together.

```go
const (
	CompositeNormal Composite = iota
	CompositeDissolve
	CompositeBehind
	CompositeMultiply
	CompositeScreen
	CompositeOverlay
	CompositeDifference
	CompositeAddition
	CompositeSubtract
	CompositeDarkenOnly
	CompositeLightenOnly
	CompositeHue
	CompositeSaturation
	CompositeColor
	CompositeValue
	CompositeDivide
	CompositeDodge
	CompositeBurn
	CompositeHardLight
	CompositeSoftLight
	CompositeGrainExtract
	CompositeGrainMerge
	CompositeLuminosity
	CompositePlus
	CompositeDestinationIn
	CompositeDestinationOut
	CompositeSourceAtop
	CompositeDestinationAtop
	CompositeColorErase
	CompositeChroma
	CompositeLightness
	CompositeVividLight
	CompositePinLight
	CompositeLinearLight
	CompositeHardMix
	CompositeExclusion
	CompositeLinearBurn
	CompositeLuminance
	CompositeErase
	CompositeMerge
	CompositeSplit
	CompositePassThrough
)
```
Composite constants.

#### func (Composite) Composite

```go
func (c Composite) Composite(b, t color.Color) color.Color
```
Composite performs the composition of two layers.

#### func (Composite) String

```go
func (c Composite) String() string
```
String returns the name of the composition.

#### type GrayAlpha

```go
type GrayAlpha struct {
	Pix    []lcolor.GrayAlpha
	Stride int
	Rect   image.Rectangle
}
```

GrayAlpha is an image of GrayAlpha pixels.

#### func  NewGrayAlpha

```go
func NewGrayAlpha(r image.Rectangle) *GrayAlpha
```
NewGrayAlpha create a new GrayAlpha image with the given bounds.

#### func (*GrayAlpha) At

```go
func (g *GrayAlpha) At(x, y int) color.Color
```
At returns the color for the pixel at the specified coords.

#### func (*GrayAlpha) Bounds

```go
func (g *GrayAlpha) Bounds() image.Rectangle
```
Bounds returns the limits of the image.

#### func (*GrayAlpha) ColorModel

```go
func (g *GrayAlpha) ColorModel() color.Model
```
ColorModel returns a color model to transform arbitrary colours into a GrayAlpha
color.

#### func (*GrayAlpha) GrayAlphaAt

```go
func (g *GrayAlpha) GrayAlphaAt(x, y int) lcolor.GrayAlpha
```
GrayAlphaAt returns a GrayAlpha colr for the specified coords.

#### func (*GrayAlpha) Opaque

```go
func (g *GrayAlpha) Opaque() bool
```
Opaque returns true if all pixels have full alpha.

#### func (*GrayAlpha) PixOffset

```go
func (g *GrayAlpha) PixOffset(x, y int) int
```
PixOffset returns the index of the element of Pix corresponding to the given
coords.

#### func (*GrayAlpha) Set

```go
func (g *GrayAlpha) Set(x, y int, c color.Color)
```
Set converts the given colour to a GrayAlpha colour and sets it at the given
coords.

#### func (*GrayAlpha) SetGrayAlpha

```go
func (g *GrayAlpha) SetGrayAlpha(x, y int, ga lcolor.GrayAlpha)
```
SetGrayAlpha sets the colour at the given coords.

#### func (*GrayAlpha) SubImage

```go
func (g *GrayAlpha) SubImage(r image.Rectangle) image.Image
```
SubImage retuns the Image viewable through the given bounds.

#### type Image

```go
type Image []Layer
```

Image represents a collection of layers.

#### func (Image) At

```go
func (g Image) At(x, y int) color.Color
```
At returns the colour at the specified coords.

#### func (Image) Bounds

```go
func (g Image) Bounds() image.Rectangle
```
Bounds returns the limits for the dimensions of the group.

#### func (Image) ColorModel

```go
func (g Image) ColorModel() color.Model
```
ColorModel represents the color model of the group. It uses the first layer to
determine the color model.

#### func (Image) SubImage

```go
func (g Image) SubImage(r image.Rectangle) image.Image
```
SubImage returns an image representing the portion of the image p visible
through r. The returned value shares pixels with the original image.

#### type Layer

```go
type Layer struct {
	Name         string
	LayerBounds  image.Rectangle // Bounds within the layer group.
	Mode         Composite
	Invisible    bool
	Transparency uint8
	image.Image
}
```

Layer represents a single layer of a multilayered image.

#### func (Layer) At

```go
func (l Layer) At(x, y int) color.Color
```
At returns the colour at the specified coords.

#### func (Layer) Bounds

```go
func (l Layer) Bounds() image.Rectangle
```
Bounds returns the limits for the dimensions of the layer.

#### func (Layer) SubImage

```go
func (l Layer) SubImage(r image.Rectangle) image.Image
```
SubImage returns an image representing the portion of the image p visible
through r. The returned value shares pixels with the original image.

#### func (Layer) SubImageLayer

```go
func (l Layer) SubImageLayer(r image.Rectangle) Layer
```
SubImageLayer returns an image representing the portion of the image p visible
through r. The returned value shares pixels with the original image.

#### type MaskedImage

```go
type MaskedImage struct {
	image.Image
	Mask *image.Gray
}
```

MaskedImage represents an image that has a to-be-applied mask.

#### func (MaskedImage) At

```go
func (m MaskedImage) At(x, y int) color.Color
```
At returns the colour at the specified coords after masking.

#### type PalettedAlpha

```go
type PalettedAlpha struct {
	Pix     []lcolor.IndexedAlpha
	Stride  int
	Rect    image.Rectangle
	Palette lcolor.AlphaPalette
}
```

PalettedAlpha represents a paletted image with an alpha channel.

#### func  NewPalettedAlpha

```go
func NewPalettedAlpha(r image.Rectangle, p lcolor.AlphaPalette) *PalettedAlpha
```
NewPalettedAlpha creates a new image that uses a palette with an alpha channel.

#### func (*PalettedAlpha) At

```go
func (p *PalettedAlpha) At(x, y int) color.Color
```
At returns the color of the pixel at the given coords.

#### func (*PalettedAlpha) Bounds

```go
func (p *PalettedAlpha) Bounds() image.Rectangle
```
Bounds returns the limits of the image.

#### func (*PalettedAlpha) ColorModel

```go
func (p *PalettedAlpha) ColorModel() color.Model
```
ColorModel a color model to tranform arbitrary colors to one in the palette.

#### func (*PalettedAlpha) IndexAlphaAt

```go
func (p *PalettedAlpha) IndexAlphaAt(x, y int) lcolor.IndexedAlpha
```
IndexAlphaAt returns the palette index and Alpha component of the given coords.

#### func (*PalettedAlpha) Opaque

```go
func (p *PalettedAlpha) Opaque() bool
```
Opaque returns true if the image is completely opaque.

#### func (*PalettedAlpha) PixOffset

```go
func (p *PalettedAlpha) PixOffset(x, y int) int
```
PixOffset returns the index of the Pix array corresponding to the given coords.

#### func (*PalettedAlpha) Set

```go
func (p *PalettedAlpha) Set(x, y int, c color.Color)
```
Set converts the given colour to the closest in the palette and sets it at the
given coords.

#### func (*PalettedAlpha) SetIndexAlpha

```go
func (p *PalettedAlpha) SetIndexAlpha(x, y int, ia lcolor.IndexedAlpha)
```
SetIndexAlpha directly set the index and alpha channels to the given coords.

#### func (*PalettedAlpha) SubImage

```go
func (p *PalettedAlpha) SubImage(r image.Rectangle) image.Image
```
SubImage retuns the Image viewable through the given bounds.

#### type RGB

```go
type RGB struct {
	Pix    []lcolor.RGB
	Stride int
	Rect   image.Rectangle
}
```

RGB is an image of RGB colours.

#### func  NewRGB

```go
func NewRGB(r image.Rectangle) *RGB
```
NewRGB create a new RGB image with the given bounds.

#### func (*RGB) At

```go
func (r *RGB) At(x, y int) color.Color
```
At returns the colour at the given coords.

#### func (*RGB) Bounds

```go
func (r *RGB) Bounds() image.Rectangle
```
Bounds returns the limits of the image.

#### func (*RGB) ColorModel

```go
func (r *RGB) ColorModel() color.Model
```
ColorModel returns a colour model that converts arbitrary colours to the RGB
space.

#### func (*RGB) Opaque

```go
func (r *RGB) Opaque() bool
```
Opaque just returns true as the alpha channel is fixed.

#### func (*RGB) PixOffset

```go
func (r *RGB) PixOffset(x, y int) int
```
PixOffset returns the index of the Pix array correspinding to the given coords.

#### func (*RGB) RGBAt

```go
func (r *RGB) RGBAt(x, y int) lcolor.RGB
```
RGBAt returns the exact RGB colour at the given coords.

#### func (*RGB) Set

```go
func (r *RGB) Set(x, y int, c color.Color)
```
Set converts the given colour to the RGB space and sets it at the given coords.

#### func (*RGB) SetRGB

```go
func (r *RGB) SetRGB(x, y int, rgb lcolor.RGB)
```
SetRGB directly set an RGB colour to the given coords.

#### func (*RGB) SubImage

```go
func (r *RGB) SubImage(rt image.Rectangle) image.Image
```
SubImage retuns the Image viewable through the given bounds.

#### type Text

```go
type Text struct {
	image.Image
	TextData
}
```

Text represents a text layer.

#### type TextData

```go
type TextData []TextDatum
```

TextData represents the stylised text.

#### func (TextData) String

```go
func (t TextData) String() string
```
String returns a flattened string.

#### type TextDatum

```go
type TextDatum struct {
	ForeColor, BackColor                   color.Color
	Size, LetterSpacing, Rise              uint32
	Bold, Italic, Underline, Strikethrough bool
	Font, Data                             string
}
```

TextDatum is a collection of styling for a single piece of text.
