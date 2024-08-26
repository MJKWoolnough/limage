# xcf
--
    import "vimagination.zapto.org/limage/xcf"

Package xcf implements an image encoder and decoder for GIMPs XCF format

## Usage

```go
var (
	ErrInvalidFileTypeID   = errors.New("invalid file type identification")
	ErrUnsupportedVersion  = errors.New("unsupported file version")
	ErrInvalidHeader       = errors.New("invalid header")
	ErrInvalidProperties   = errors.New("invalid property list")
	ErrInvalidOpacity      = errors.New("opacity not in valid range")
	ErrInvalidGuideLength  = errors.New("invalid guide length")
	ErrInvalidUnit         = errors.New("invalid unit")
	ErrInvalidSampleLength = errors.New("invalid sample points length")
	ErrInvalidGroup        = errors.New("invalid or unknown group specified for layer")
	ErrUnknownCompression  = errors.New("unknown compression method")
	ErrMissingAlpha        = errors.New("non-bottom layer missing alpha channel")
	ErrNeedReaderAt        = errors.New("need a io.ReaderAt")
)
```
Errors.

```go
var (
	ErrInvalidLayerType      = errors.New("invalid layer type")
	ErrInvalidItemPathLength = errors.New("invalid item path length")
	ErrInconsistantData      = errors.New("inconsistent data read")
)
```
Errors.

```go
var (
	ErrInvalidParasites = errors.New("invalid parasites layout")
	ErrNoOpen           = errors.New("didn't receive Open token")
	ErrNoName           = errors.New("didn't receive Name token")
)
```
Errors.

```go
var (
	ErrInconsistantClosedState = errors.New("inconsistent closed state")
	ErrUnknownPathsVersion     = errors.New("unknown paths version")
)
```
Errors.

```go
var (
	ErrInvalidString = errors.New("string is invalid")
	ErrStringTooLong = errors.New("string exceeds maximum length")
	ErrInvalidSeek   = errors.New("invalid seek")
)
```
Errors.

```go
var (
	ErrUnknownVectorVersion = errors.New("unknown vector version")
	ErrUnknownStrokeType    = errors.New("unknown stroke type")
	ErrInvalidFloatsNumber  = errors.New("invalids number of floats")
)
```
Errors.

```go
var (
	ErrInvalidBoolean = errors.New("invalid boolean value")
)
```
Errors.

```go
var (
	ErrInvalidRLE = errors.New("invalid RLE data")
)
```
Errors.

```go
var (
	ErrTooBig = errors.New("write too big")
)
```
Errors.

#### func  Decode

```go
func Decode(r io.ReaderAt) (limage.Image, error)
```
Decode reads an XCF layered image from the given ReaderAt.

#### func  DecodeCompressed

```go
func DecodeCompressed(r io.ReaderAt) (limage.Image, error)
```
DecodeCompressed reads an XCF layered image, as Decode, but defers decoding and
decompressing, doing so upon an At method.

#### func  DecodeConfig

```go
func DecodeConfig(r io.ReaderAt) (image.Config, error)
```
DecodeConfig retrieves the color model and dimensions of the XCF image.

#### func  Encode

```go
func Encode(w io.WriterAt, im image.Image) error
```
Encode encodes the given image as an XCF file to the given WriterAt.

#### type CompressedGray

```go
type CompressedGray struct {
	Rect image.Rectangle
}
```

CompressedGray is an image.Image for which the data remains in a compressed form
until read.

#### func (*CompressedGray) At

```go
func (c *CompressedGray) At(x, y int) color.Color
```
At returns colour at the specified coords.

#### func (*CompressedGray) Bounds

```go
func (c *CompressedGray) Bounds() image.Rectangle
```
Bounds returns a Rect containing the boundary data for the image.

#### func (CompressedGray) ColorModel

```go
func (CompressedGray) ColorModel() color.Model
```
ColorModel returns the Gray Color Model.

#### func (*CompressedGray) GrayAt

```go
func (c *CompressedGray) GrayAt(x, y int) color.Gray
```
GrayAt returns Gray colour at the specified coords.

#### type CompressedGrayAlpha

```go
type CompressedGrayAlpha struct {
	Rect image.Rectangle
}
```

CompressedGrayAlpha is an image.Image for which the data remains in a compressed
form until read.

#### func (*CompressedGrayAlpha) At

```go
func (c *CompressedGrayAlpha) At(x, y int) color.Color
```
At returns colour at the specified coords.

#### func (*CompressedGrayAlpha) Bounds

```go
func (c *CompressedGrayAlpha) Bounds() image.Rectangle
```
Bounds returns a Rect containing the boundary data for the image.

#### func (CompressedGrayAlpha) ColorModel

```go
func (CompressedGrayAlpha) ColorModel() color.Model
```
ColorModel returns the Gray Alpha Color Model.

#### func (*CompressedGrayAlpha) GrayAlphaAt

```go
func (c *CompressedGrayAlpha) GrayAlphaAt(x, y int) lcolor.GrayAlpha
```
GrayAlphaAt returns Gray+Alpha colour at the specified coords.

#### type CompressedNRGBA

```go
type CompressedNRGBA struct {
	Rect image.Rectangle
}
```

CompressedNRGB is an image.Image for which the data remains in a compressed form
until read.

#### func (*CompressedNRGBA) At

```go
func (c *CompressedNRGBA) At(x, y int) color.Color
```
At returns colour at the specified coords.

#### func (*CompressedNRGBA) Bounds

```go
func (c *CompressedNRGBA) Bounds() image.Rectangle
```
Bounds returns a Rect containing the boundary data for the image.

#### func (CompressedNRGBA) ColorModel

```go
func (CompressedNRGBA) ColorModel() color.Model
```
ColorModel returns the NRGBA Color Model.

#### func (*CompressedNRGBA) NRGBAAt

```go
func (c *CompressedNRGBA) NRGBAAt(x, y int) color.NRGBA
```
NRGBAAt returns NRGBA colour at the specified coords.

#### type CompressedPaletted

```go
type CompressedPaletted struct {
	Rect    image.Rectangle
	Palette color.Palette
}
```

CompressedPaletted is an image.Image for which the data remains in a compressed
form until read.

#### func (*CompressedPaletted) At

```go
func (c *CompressedPaletted) At(x, y int) color.Color
```
At returns colour at the specified coords.

#### func (*CompressedPaletted) Bounds

```go
func (c *CompressedPaletted) Bounds() image.Rectangle
```
Bounds returns a Rect containing the boundary data for the image.

#### func (*CompressedPaletted) ColorModel

```go
func (c *CompressedPaletted) ColorModel() color.Model
```
ColorModel returns the Palette of the image.

#### type CompressedPalettedAlpha

```go
type CompressedPalettedAlpha struct {
	Rect    image.Rectangle
	Palette lcolor.AlphaPalette
}
```

CompressedPalettedAlpha is an image.Image for which the data remains in a
compressed form until read.

#### func (*CompressedPalettedAlpha) At

```go
func (c *CompressedPalettedAlpha) At(x, y int) color.Color
```
At returns colour at the specified coords.

#### func (*CompressedPalettedAlpha) Bounds

```go
func (c *CompressedPalettedAlpha) Bounds() image.Rectangle
```
Bounds returns a Rect containing the boundary data for the image

#### func (*CompressedPalettedAlpha) ColorModel

```go
func (c *CompressedPalettedAlpha) ColorModel() color.Model
```
ColorModel returns the Palette of the image

#### type CompressedRGB

```go
type CompressedRGB struct {
	Rect image.Rectangle
}
```

CompressedRGB is an image.Image for which the data remains in a compressed form
until read.

#### func (*CompressedRGB) At

```go
func (c *CompressedRGB) At(x, y int) color.Color
```
At returns colour at the specified coords.

#### func (*CompressedRGB) Bounds

```go
func (c *CompressedRGB) Bounds() image.Rectangle
```
Bounds returns a Rect containing the boundary data for the image.

#### func (CompressedRGB) ColorModel

```go
func (CompressedRGB) ColorModel() color.Model
```
ColorModel returns the RGB Color Model.

#### func (*CompressedRGB) RGBAt

```go
func (c *CompressedRGB) RGBAt(x, y int) lcolor.RGB
```
RGBAt returns RGB colour at the specified coords.
