# xcf
--
    import "vimagination.zapto.org/limage/xcf"


## Usage

```go
const (
	ErrInvalidFileTypeID   errors.Error = "invalid file type identification"
	ErrUnsupportedVersion  errors.Error = "unsupported file version"
	ErrInvalidHeader       errors.Error = "invalid header"
	ErrInvalidProperties   errors.Error = "invalid property list"
	ErrInvalidOpacity      errors.Error = "opacity not in valid range"
	ErrInvalidGuideLength  errors.Error = "invalid guide length"
	ErrInvalidUnit         errors.Error = "invalid unit"
	ErrInvalidSampleLength errors.Error = "invalid sample points length"
	ErrInvalidGroup        errors.Error = "invalid or unknown group specified for layer"
	ErrUnknownCompression  errors.Error = "unknown compression method"
	ErrMissingAlpha        errors.Error = "non-bottom layer missing alpha channel"
)
```
Errors

```go
const (
	ErrInvalidLayerType      errors.Error = "invalid layer type"
	ErrInvalidItemPathLength errors.Error = "invalid item path length"
	ErrInconsistantData      errors.Error = "inconsistant data read"
)
```
Errors

```go
const (
	ErrInvalidParasites errors.Error = "invalid parasites layout"
	ErrNoOpen           errors.Error = "didn't receive Open token"
	ErrNoName           errors.Error = "didn't receive Name token"
)
```
Errors

```go
const (
	ErrInconsistantClosedState errors.Error = "inconsistant closed state"
	ErrUnknownPathsVersion     errors.Error = "unknown paths version"
)
```
Errors

```go
const (
	ErrInvalidString errors.Error = "string is invalid"
	ErrStringTooLong errors.Error = "string exceeds maximum length"
	ErrInvalidSeek   errors.Error = "invalid seek"
)
```
Errors

```go
const (
	ErrUnknownVectorVersion errors.Error = "unknown vector version"
	ErrUnknownStrokeType    errors.Error = "unknown stroke type"
	ErrInvalidFloatsNumber  errors.Error = "invalids number of floats"
)
```
Errors

```go
const (
	ErrInvalidBoolean errors.Error = "invalid boolean value"
)
```
Errors

```go
const (
	ErrInvalidRLE errors.Error = "invalid RLE data"
)
```
Errors

```go
const (
	ErrTooBig errors.Error = "write too big"
)
```
Errors

#### func  Decode

```go
func Decode(r io.ReaderAt) (limage.Image, error)
```
Decode reads an XCF layered image from the given ReaderAt

#### func  DecodeCompressed

```go
func DecodeCompressed(r io.ReaderAt) (limage.Image, error)
```
DecodeCompressed reads an XCF layered image, as Decode, but defers decoding and
decompressing, doing so upon an At method

#### func  DecodeConfig

```go
func DecodeConfig(r io.ReaderAt) (image.Config, error)
```
DecodeConfig retrieves the color model and dimensions of the XCF image

#### func  Encode

```go
func Encode(w io.WriterAt, im image.Image) error
```
Encode encodes the given image as an XCF file to the given WriterAt

#### type CompressedGray

```go
type CompressedGray struct {
	Rect image.Rectangle
}
```


#### func (*CompressedGray) At

```go
func (c *CompressedGray) At(x, y int) color.Color
```

#### func (*CompressedGray) Bounds

```go
func (c *CompressedGray) Bounds() image.Rectangle
```

#### func (CompressedGray) ColorModel

```go
func (CompressedGray) ColorModel() color.Model
```

#### func (*CompressedGray) GrayAt

```go
func (c *CompressedGray) GrayAt(x, y int) color.Gray
```

#### type CompressedGrayAlpha

```go
type CompressedGrayAlpha struct {
	Rect image.Rectangle
}
```


#### func (*CompressedGrayAlpha) At

```go
func (c *CompressedGrayAlpha) At(x, y int) color.Color
```

#### func (*CompressedGrayAlpha) Bounds

```go
func (c *CompressedGrayAlpha) Bounds() image.Rectangle
```

#### func (CompressedGrayAlpha) ColorModel

```go
func (CompressedGrayAlpha) ColorModel() color.Model
```

#### func (*CompressedGrayAlpha) GrayAlphaAt

```go
func (c *CompressedGrayAlpha) GrayAlphaAt(x, y int) lcolor.GrayAlpha
```

#### type CompressedNRGBA

```go
type CompressedNRGBA struct {
	Rect image.Rectangle
}
```


#### func (*CompressedNRGBA) At

```go
func (c *CompressedNRGBA) At(x, y int) color.Color
```

#### func (*CompressedNRGBA) Bounds

```go
func (c *CompressedNRGBA) Bounds() image.Rectangle
```

#### func (CompressedNRGBA) ColorModel

```go
func (CompressedNRGBA) ColorModel() color.Model
```

#### func (*CompressedNRGBA) NRGBAAt

```go
func (c *CompressedNRGBA) NRGBAAt(x, y int) color.NRGBA
```

#### type CompressedPaletted

```go
type CompressedPaletted struct {
	Rect    image.Rectangle
	Palette color.Palette
}
```


#### func (*CompressedPaletted) At

```go
func (c *CompressedPaletted) At(x, y int) color.Color
```

#### func (*CompressedPaletted) Bounds

```go
func (c *CompressedPaletted) Bounds() image.Rectangle
```

#### func (*CompressedPaletted) ColorModel

```go
func (c *CompressedPaletted) ColorModel() color.Model
```

#### type CompressedPalettedAlpha

```go
type CompressedPalettedAlpha struct {
	Rect    image.Rectangle
	Palette lcolor.AlphaPalette
}
```


#### func (*CompressedPalettedAlpha) At

```go
func (c *CompressedPalettedAlpha) At(x, y int) color.Color
```

#### func (*CompressedPalettedAlpha) Bounds

```go
func (c *CompressedPalettedAlpha) Bounds() image.Rectangle
```

#### func (*CompressedPalettedAlpha) ColorModel

```go
func (c *CompressedPalettedAlpha) ColorModel() color.Model
```

#### type CompressedRGB

```go
type CompressedRGB struct {
	Rect image.Rectangle
}
```


#### func (*CompressedRGB) At

```go
func (c *CompressedRGB) At(x, y int) color.Color
```

#### func (*CompressedRGB) Bounds

```go
func (c *CompressedRGB) Bounds() image.Rectangle
```

#### func (CompressedRGB) ColorModel

```go
func (CompressedRGB) ColorModel() color.Model
```

#### func (*CompressedRGB) RGBAt

```go
func (c *CompressedRGB) RGBAt(x, y int) lcolor.RGB
```
