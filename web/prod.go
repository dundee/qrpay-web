// +build !live

package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static
var embededFiles embed.FS

func getFileSystem() http.FileSystem {
	log.Print("using embed mode")
	fsys, err := fs.Sub(embededFiles, "static")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
