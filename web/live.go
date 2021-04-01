// +build live

package web

import (
	"log"
	"net/http"
)

func getFileSystem() http.FileSystem {
	log.Print("using live mode")
	return http.Dir("./web/static")
}
