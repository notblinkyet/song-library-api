package http

import (
	"net/http"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (h *Handler) FillEndpoints(r *chi.Mux) {

	r.Post("/songs", h.CreateSong)
	r.Get("/songs", h.ReadFilteredSongs)
	r.Get("/songs/{id}", h.ReadVerse)
	r.Patch("/songs/{id}", h.UpdateSong)
	r.Delete("/songs/{id}", h.DeleteSong)
	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})
	r.Get("/swagger/*", httpSwagger.WrapHandler)
}
