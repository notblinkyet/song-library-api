package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/notblinkyet/song-library-api/internal/lib/sl"
	"github.com/notblinkyet/song-library-api/internal/models"
)

// @Summary Update a song by ID
// @Description Updates an existing song in the library by its ID. Only the specified fields in the request body will be updated. The request body must contain a valid JSON representation of the `models.Song` struct.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Updated song information. Only fields with non-empty values will be updated."
// @Success 200 {object} models.Song "Successfully updated song"
// @Failure 400 {object} string "Invalid request (e.g., invalid JSON or missing required fields)."
// @Failure 404 {object} string "Song not found."
// @Failure 500 {object} string "Internal server error during the update process."
// @Router /songs/{id} [patch]]
func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	h.log.Info("received request to update a song")

	// Parse the song ID from the URL parameter.
	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		h.log.Error("failed to parse song ID", sl.Error(err))
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	// Retrieve the current version of the song from the service.
	song, err := h.service.ReadByID(id)
	if err != nil {
		h.log.Error("failed to retrieve song", slog.Int("id", id), sl.Error(err))
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	// Parse and decode the updated song data from the request body.
	var newInfo models.Song
	err = json.NewDecoder(r.Body).Decode(&newInfo)
	if err != nil {
		h.log.Error("failed to decode request body", sl.Error(err))
		return
	}

	// Update the song fields with the new data.
	if newInfo.Title != "" {
		song.Title = newInfo.Title
	}
	if newInfo.Group != "" {
		song.Group = newInfo.Group
	}
	if newInfo.ReleaseDate.Equal(time.Time{}) {
		song.ReleaseDate = newInfo.ReleaseDate
	}
	if newInfo.Text != "" {
		song.Text = newInfo.Text
	}
	if newInfo.Link != "" {
		song.Link = newInfo.Link
	}
	h.log.Debug("updated song fields", slog.Any("song", song))

	// Save the updated song data.
	err = h.service.UpdateSong(song)
	if err != nil {
		h.log.Error("failed to update song", slog.Int("id", id), sl.Error(err))
		http.Error(w, "Failed to update song", http.StatusInternalServerError)
		return
	}
	h.log.Info("song updated successfully", slog.Int("id", id))
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent(" ", "\t")
	err = encoder.Encode(&song)
	if err != nil {
		h.log.Error("failed to decode request body", sl.Error(err))
		return
	}
}
