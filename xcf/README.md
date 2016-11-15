# xcf
--
    import "github.com/MJKWoolnough/limage/xcf"


## Usage

```go
const (
	ErrInvalidFileTypeID       errors.Error = "invalid file type identification"
	ErrUnsupportedVersion      errors.Error = "unsupported file version"
	ErrInvalidHeader           errors.Error = "invalid header"
	ErrInvalidProperties       errors.Error = "invalid property list"
	ErrInvalidOpacity          errors.Error = "opacity not in valid range"
	ErrInvalidGuideLength      errors.Error = "invalid guide length"
	ErrInvalidUnit             errors.Error = "invalid unit"
	ErrInvalidSampleLength     errors.Error = "invalid sample points length"
	ErrInvalidGroup            errors.Error = "invalid or unknown group specified for layer"
	ErrUnknownCompression      errors.Error = "unknown compression method"
	ErrMissingAlpha            errors.Error = "non-bottom layer missing alpha channel"
	ErrInvalidLayerType        errors.Error = "invalid layer type"
	ErrInvalidItemPathLength   errors.Error = "invalid item path length"
	ErrInconsistantData        errors.Error = "inconsistant data read"
	ErrInvalidParasites        errors.Error = "invalid parasites layout"
	ErrNoOpen                  errors.Error = "didn't receive Open token"
	ErrNoName                  errors.Error = "didn't receive Name token"
	ErrInconsistantClosedState errors.Error = "inconsistant closed state"
	ErrUnknownPathsVersion     errors.Error = "unknown paths version"
	ErrInvalidString           errors.Error = "string is invalid"
	ErrStringTooLong           errors.Error = "string exceeds maximum length"
	ErrInvalidSeek             errors.Error = "invalid seek"
	ErrUnknownVectorVersion    errors.Error = "unknown vector version"
	ErrUnknownStrokeType       errors.Error = "unknown stroke type"
	ErrInvalidFloatsNumber     errors.Error = "invalids number of floats"
	ErrInvalidBoolean          errors.Error = "invalid boolean value"
	ErrInvalidRLE              errors.Error = "invalid RLE data"
	ErrTooBig                  errors.Error = "write too big"
)
```
Errors

#### func  Decode

```go
func Decode(r io.ReaderAt) (limage.Image, error)
```
Decode reads an XCF layered image from the given ReaderAt

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
