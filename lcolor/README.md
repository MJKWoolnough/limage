# lcolor
--
    import "vimagination.zapto.org/limage/lcolor"


## Usage

```go
var (
	GrayAlphaModel = color.ModelFunc(grayAlphaColourModel)
	RGBModel       = color.ModelFunc(rgbColourModel)
)
```
Color Models.

#### type AlphaPalette

```go
type AlphaPalette color.Palette
```

AlphaPalette is an extension to the normal color.Palette type.

#### func (AlphaPalette) Convert

```go
func (ap AlphaPalette) Convert(c color.Color) color.Color
```
Convert converts the given color to the nearest color in the palette, but
preserves the alpha channel.

#### func (AlphaPalette) Index

```go
func (ap AlphaPalette) Index(c color.Color) int
```
Index returns the palette index of the nearest color.

#### type GrayAlpha

```go
type GrayAlpha struct {
	Y, A uint8
}
```

GrayAlpha represents a Gray color with an Alpha channel.

#### func (GrayAlpha) RGBA

```go
func (c GrayAlpha) RGBA() (r, g, b, a uint32)
```
RGBA implements the color.Color interface.

#### func (GrayAlpha) ToNRGBA

```go
func (c GrayAlpha) ToNRGBA() color.NRGBA64
```
ToNRGBA converts the HSV color into the RGB colorspace.

#### type HSLA

```go
type HSLA struct {
	H, S, L, A uint16
}
```

HSLA represents the Hue, Saturation, Lightness and Alpha of a pixel.

#### func  RGBToHSL

```go
func RGBToHSL(cl color.Color) HSLA
```
RGBToHSL converts an RGC color.Color to HSLA format.

#### func (HSLA) RGBA

```go
func (h HSLA) RGBA() (uint32, uint32, uint32, uint32)
```
RGBA implements the color.Color interface.

#### func (HSLA) ToNRGBA

```go
func (h HSLA) ToNRGBA() color.NRGBA64
```
ToNRGBA converts the HSL color into the RGB colorspace.

#### type HSVA

```go
type HSVA struct {
	H, S, V, A uint16
}
```

HSVA represents the Hue, Saturation, Value and Alpha of a pixel.

#### func  RGBToHSV

```go
func RGBToHSV(cl color.Color) HSVA
```
RGBToHSV converts a color to the HSV color space.

#### func (HSVA) RGBA

```go
func (h HSVA) RGBA() (uint32, uint32, uint32, uint32)
```
RGBA implements the color.Color interface.

#### func (HSVA) ToNRGBA

```go
func (h HSVA) ToNRGBA() color.NRGBA64
```
ToNRGBA converts the HSV color into the RGB colorspace.

#### type IndexedAlpha

```go
type IndexedAlpha struct {
	I, A uint8
}
```

IndexedAlpha is the combination of a palette index and an Alpha channel.

#### type RGB

```go
type RGB struct {
	R, G, B uint8
}
```

RGB is a standard colour type whose Alpha channel is always full.

#### func (RGB) RGBA

```go
func (rgb RGB) RGBA() (r, g, b, a uint32)
```
RGBA implements the color.Color interface.

#### func (RGB) ToNRGBA

```go
func (rgb RGB) ToNRGBA() color.NRGBA64
```
ToNRGBA returns itself as a non-alpha-premultiplied value As the alpha is always
full, this only returns the normal values.
