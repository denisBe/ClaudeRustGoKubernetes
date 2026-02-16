package main

import (
	"embed"
	"io/fs"
)

//go:embed static
var staticFiles embed.FS

// StaticFS returns the embedded static files as an fs.FS rooted at "static/".
// This strips the "static/" prefix so files are served at /static/index.html etc.
func StaticFS() (fs.FS, error) {
	return fs.Sub(staticFiles, "static")
}
