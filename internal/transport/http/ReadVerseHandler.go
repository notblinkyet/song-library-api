package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	parseurl "github.com/notblinkyet/song-library-api/internal/lib/ParseURL"
	"github.com/notblinkyet/song-library-api/internal/lib/sl"
	"github.com/notblinkyet/song-library-api/internal/services"
)

// @Summary Retrieve verses of a song by ID
// @Description Retrieves the verses of a song by its ID.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param start query int false "Start index for verse retrieval (1-based index). Defaults to 1."
// @Param count query int false "Number of verses to retrieve. Defaults to 1."
// @Success 200 {object} []models.Verse "Verses of the song"
// @Failure 400 {object} string "Invalid request parameters or song does not contain requested verses."
// @Failure 404 {object} string "Song not found."
// @Router /songs/{id} [get]
func (h *Handler) ReadVerse(w http.ResponseWriter, r *http.Request) {
	h.log.Info("received request to read text by ID")

	// Parse the song ID from the URL parameter.
	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		h.log.Error("failed to parse song ID", sl.Error(err))
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}
	h.log.Debug("parsed song ID", slog.Int("id", id))

	// Parse additional query parameters for verse retrieval.
	start := parseurl.ParseInt(r.URL.Query(), "start", 1)
	count := parseurl.ParseInt(r.URL.Query(), "count", 1)

	// Call the service layer to retrieve the requested verses.
	verse, err := h.service.ReadVerse(id, start, count)
	if err != nil {
		if errors.Is(err, services.ErrVerseOutOfBound) {
			h.log.Warn("song does not contain requested verses", slog.Int("id", id))
			http.Error(w, services.ErrVerseOutOfBound.Error(), http.StatusBadRequest)
			return
		}
		h.log.Error("failed to retrieve verses", sl.Error(err))
		http.Error(w, "failed to retrieve verses", http.StatusBadRequest)
		return
	}
	h.log.Info("verses retrieved successfully", slog.Int("id", id))

	// Return the verses in the response.
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent(" ", "\t")
	err = encoder.Encode(verse)
	if err != nil {
		h.log.Error("failed to encode verses", sl.Error(err))
		return
	}
}
