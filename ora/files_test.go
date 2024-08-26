package ora

import _ "embed"

var (
	//go:embed testfiles/abc.ora
	abcFile string

	//go:embed testfiles/blackRedBlue.ora
	blackRedBlueFile string

	//go:embed testfiles/blackMask.ora
	blackMaskFile string

	//go:embed testfiles/black.ora
	blackFile string

	//go:embed testfiles/blackRed.ora
	blackRedFile string

	//go:embed testfiles/red.ora
	redFile string

	//go:embed testfiles/white.ora
	whiteFile string
)
