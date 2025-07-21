package public

import (
	"embed"
)

//go:embed *.html
var PublicFS embed.FS
