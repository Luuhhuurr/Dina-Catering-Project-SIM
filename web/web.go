// Package web membungkus template HTML dan aset statis agar ikut tertanam di
// dalam binary aplikasi.
package web

import (
	"embed"
	"io/fs"
)

//go:embed templates
var TemplatesFS embed.FS

//go:embed static
var StaticFS embed.FS

// Static mengembalikan filesystem aset statis (css, js) siap dilayani.
func Static() fs.FS {
	sub, err := fs.Sub(StaticFS, "static")
	if err != nil {
		panic(err)
	}
	return sub
}
