package web

import (
	"embed"
	"io/fs"
	"net/http"
)

var (
	//go:embed dist/index.html
	HomeHTML []byte
	//go:embed dist/index.html
	LoginHTML []byte
)

//go:embed all:dist
var files embed.FS

func NewFileSystem() http.FileSystem {
	subfs, _ := fs.Sub(files, "dist")
	return http.FS(subfs)
}
