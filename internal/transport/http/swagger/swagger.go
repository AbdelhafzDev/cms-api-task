package swagger

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed ui
var uiFS embed.FS

func RegisterRoutes(r chi.Router, specFS embed.FS, specPath string) {
	r.Get("/api/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		data, err := specFS.ReadFile(specPath)
		if err != nil {
			http.Error(w, "OpenAPI spec not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/x-yaml")
		_, _ = w.Write(data)
	})

	uiSubFS, err := fs.Sub(uiFS, "ui")
	if err != nil {
		panic("failed to get swagger ui subdirectory: " + err.Error())
	}

	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	r.Mount("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.FS(uiSubFS))))
}
