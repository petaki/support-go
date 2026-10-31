package web

import (
	"io/fs"
	"net/http"
)

// Static function.
func Static(debug bool, directory string, files fs.FS) http.Handler {
	if debug {
		return http.FileServer(http.Dir(directory))
	}

	return http.FileServer(http.FS(files))
}
