package public

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed login.html
var PublicFS embed.FS

func init() {
	entries, err := fs.ReadDir(content, ".")
	if err != nil {
		panic(err)
	}
	for _, e := range entries {
		fmt.Println("🔍 embedded file:", e.Name())
	}
}
