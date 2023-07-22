package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func ClientRouter() http.Handler {
	// Create a route along /files that will serve contents from
	workDir, _ := os.Getwd()

	r := chi.NewRouter()

	// Specify the relative path to the files folder
	spaPath := filepath.Join(workDir, "client", "dist")

	// Serve static files from the root directory
	rootDir := "/"
	rootFiles := http.FileServer(http.Dir(filepath.Join(spaPath, rootDir)))
	r.Handle(rootDir+"*", http.StripPrefix(rootDir, rootFiles))

	// Serve static files from the specified directory
	assetsDir := "/assets/"
	staticFiles := http.FileServer(http.Dir(filepath.Join(spaPath, assetsDir)))
	r.Handle(assetsDir+"*", http.StripPrefix(assetsDir, staticFiles))

	indexPath := filepath.Join(spaPath, "index.html")

	// Serve React Build files
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, indexPath)
	})

	//? Redirect to client if route not found
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, indexPath)
	})

	return r
}
