package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/notblinkyet/song-library-api/internal/database/postgresql"
	"github.com/notblinkyet/song-library-api/internal/lib/sl"
)

var (
	ErrFailToDeleteSong = errors.New("fail to delete")
)

// @Summary Delete a song by ID
// @Description Deletes a song from the library by its ID.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} nil
// @Failure 400 {object} string "Invalid song ID"
// @Failure 500 {object} string "Internal server error during deletion"
// @Router /songs/{id} [delete]
func (h *Handler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	h.log.Info("received request to delete a song")

	// Parse the song ID from the URL parameter.
	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		h.log.Error("failed to parse song ID", sl.Error(err))
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	// Call the service layer to delete the song by ID.
	err = h.service.DeleteSong(id)
	if err != nil {
		h.log.Error("failed to delete song", slog.Int("id", id), sl.Error(err))
		if errors.Is(err, postgresql.ErrNotFound) {
			http.Error(w, postgresql.ErrNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "failed to delete song", http.StatusInternalServerError)
		return
	}
	h.log.Info("song deleted successfully", slog.Int("id", id))
	w.WriteHeader(http.StatusOK)
}
