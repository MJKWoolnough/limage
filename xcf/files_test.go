package xcf

import _ "embed"

var (
	//go:embed testfiles/abc.xcf
	abcFile string

	//go:embed testfiles/blackRedBlue.xcf
	blackRedBlueFile string

	//go:embed testfiles/blackMask.xcf
	blackMaskFile string

	//go:embed testfiles/black.xcf
	blackFile string

	//go:embed testfiles/blackRed.xcf
	blackRedFile string

	//go:embed testfiles/red.xcf
	redFile string

	//go:embed testfiles/white.xcf
	whiteFile string
)
