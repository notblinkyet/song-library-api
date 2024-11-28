package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/notblinkyet/song-library-api/internal/lib/api"
	"github.com/notblinkyet/song-library-api/internal/lib/sl"
	"github.com/notblinkyet/song-library-api/internal/models"
)

// Handler provides HTTP handlers for song-related operations.
type Handler struct {
	service SongLibraryService
	log     *slog.Logger
}

// NewHandler initializes and returns a new Handler instance.
func NewHandler(service SongLibraryService, log *slog.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

// @Summary Create a new song
// @Description Creates a new song in the library. Requires a valid JSON request body containing the song's details.
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.CreateSongRequest true "Song details"
// @Success 201 {object} models.Id "Successfully created song"
// @Failure 400 {object} string "Invalid request (e.g., missing required fields)"
// @Failure 500 {object} string "Internal server error"
// @Router /songs [post]
func (h *Handler) CreateSong(w http.ResponseWriter, r *http.Request) {
	h.log.Info("received request to create a new song")

	// Parse and decode the request body into the CreateSongRequest model.
	var req models.CreateSongRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.log.Error("failed to decode request body", sl.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	h.log.Debug("decoded request body", slog.Any("request", req))

	// Validate required fields in the request.
	if req.Group == "" || req.Title == "" {
		h.log.Warn("missing required fields: group or title")
		http.Error(w, "Group and title are required fields", http.StatusBadRequest)
		return
	}

	// Call the service layer to create the song and retrieve the new song's ID.
	id, err := h.service.Create(&req)
	if err != nil {
		// Handle specific errors returned by the service.
		if errors.Is(err, api.ErrInternalServer) {
			h.log.Error("internal server error during song creation", sl.Error(err))
			http.Error(w, "internal server error during song creation", http.StatusInternalServerError)
		} else {
			h.log.Error("failed to create song", sl.Error(err))
			http.Error(w, "failed to create song", http.StatusBadRequest)
		}
		return
	}
	h.log.Info("song created successfully", slog.Int("songID", id))

	// Return the ID of the newly created song in the response.
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.SetIndent(" ", "\t")
	err = encoder.Encode(models.NewId(id))
	if err != nil {
		h.log.Error("failed to encode song ID", sl.Error(err))
		return
	}
}
