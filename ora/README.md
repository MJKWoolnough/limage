# ora
--
    import "vimagination.zapto.org/limage/ora"


## Usage

```go
var (
	ErrMissingStack    = errors.New("missing stack file")
	ErrInvalidMimeType = errors.New("invalid mime type")
	ErrInvalidStack    = errors.New("invalid stack")
)
```
Errors

```go
var (
	ErrInvalidSource = errors.New("invalid source")
)
```
Errors

#### func  Decode

```go
func Decode(zr *zip.Reader) (limage.Image, error)
```
Decode reads an ORA layered image from the given Reader

It accepts a *zip.Reader and it is the callers responsibility to handle it

#### func  DecodeConfig

```go
func DecodeConfig(zr *zip.Reader) (image.Config, error)
```
DecodeConfig retrieves the color model and dimensions of the ORA image.

It accepts a *zip.Reader and it is the callers responsibility to handle it

#### func  Encode

```go
func Encode(w io.Writer, m image.Image) error
```
Encode encodes the given image as an ORA file to the given Writer
