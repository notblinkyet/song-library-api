package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	parseurl "github.com/notblinkyet/song-library-api/internal/lib/ParseURL"
	"github.com/notblinkyet/song-library-api/internal/lib/sl"
	"github.com/notblinkyet/song-library-api/internal/models"
)

// @Summary Retrieve songs based on filters
// @Description Retrieves a list of songs from the library based on optional filters.
// @Tags songs
// @Accept json
// @Produce json
// @Param song query string false "Song title"
// @Param group query string false "Song group"
// @Param release_date query string false "Song release date YYYY.MM.DD"
// @Param text query string false "Text search in song details"
// @Param link query string false "Link search in song details"
// @Success 200 {array} models.Song "Successfully retrieved songs"
// @Failure 400 {object} string "Invalid request (e.g., invalid filter parameters)"
// @Failure 500 {object} string "Internal server error"
// @Router /songs [get]
func (h *Handler) ReadFilteredSongs(w http.ResponseWriter, r *http.Request) {
	h.log.Info("received request to read filtered songs")

	// Extract filtering parameters from the query string.
	var filter models.Filter
	values := r.URL.Query()
	filter.Title = parseurl.ParseString(values, "song", "")
	filter.Group = parseurl.ParseString(values, "group", "")
	filter.ReleaseDate = parseurl.ParseTime(values, "release_date", time.Time{})
	filter.Text = parseurl.ParseString(values, "text", "")
	filter.Link = parseurl.ParseString(values, "link", "")
	filter.Limit = parseurl.ParseInt(values, "limit", 0)
	filter.Offset = parseurl.ParseInt(values, "offset", 0)

	h.log.Debug("filter parameters extracted", slog.Any("filter", filter))

	// Call the service layer to retrieve the filtered list of songs.
	songs, err := h.service.ReadFilteredSongs(&filter)
	if err != nil {
		h.log.Error("failed to retrieve songs by filter", sl.Error(err))
		http.Error(w, "Failed to retrieve songs", http.StatusBadRequest)
		return
	}
	h.log.Info("songs retrieved successfully", slog.Int("count", len(songs)))

	// Return the filtered songs in the response.
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent(" ", "\t")
	err = encoder.Encode(&songs)
	if err != nil {
		h.log.Error("failed to encode songs", sl.Error(err))
		http.Error(w, "Failed to encode", http.StatusInternalServerError)
		return
	}
}
