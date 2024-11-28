package services

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/notblinkyet/song-library-api/internal/database"
	"github.com/notblinkyet/song-library-api/internal/lib/api"
	"github.com/notblinkyet/song-library-api/internal/lib/sl"
	"github.com/notblinkyet/song-library-api/internal/models"
)

// Predefined error
var (
	ErrVerseOutOfBound = errors.New("this song doesn't have so many verses")
)

// ApiClient defines the interface for external API interactions.
type ApiClient interface {
	GetMoreAboutSong(req *models.CreateSongRequest) (*models.Song, error)
}

// SongLibraryService handles all business logic for song-related operations.
type SongLibraryService struct {
	SingStorage database.Storage // Database storage for songs.
	ApiClient   ApiClient        // External API client.
	log         *slog.Logger     // Logger for structured logging.
}

// NewSongLibraryService initializes and returns a new SongLibraryService instance.
func NewSongLibraryService(SingStorage database.Storage, a ApiClient, log *slog.Logger) *SongLibraryService {
	return &SongLibraryService{
		SingStorage: SingStorage,
		ApiClient:   a,
		log:         log,
	}
}

// Create handles the creation of a new song by saving it to the database.
func (s *SongLibraryService) Create(req *models.CreateSongRequest) (int, error) {
	s.log.Info("saving song in the database")

	// Retrieve additional song details from the external API.
	song, err := s.ApiClient.GetMoreAboutSong(req)
	if err != nil {
		if errors.Is(err, api.ErrBadRequest) {
			// Log the error if the API request is invalid.
			s.log.Error("bad request", sl.Error(err))
		} else {
			// Log a warning for other types of API errors.
			s.log.Warn("failed to get info from API", sl.Error(err))
		}
		return 0, err
	}

	// Log the retrieved song details.
	s.log.Debug("retrieved information about song", slog.Any("song", song))

	// Save the song to the database and return the new song's ID.
	id, err := s.SingStorage.CreateSong(song)
	if err != nil {
		// Log any database insertion errors.
		s.log.Error("failed to insert song into database", sl.Error(err))
		return 0, err
	}

	// Log the successful insertion of the song.
	s.log.Debug("song successfully saved to database", slog.Int("id", id))

	return id, nil
}

// ReadFilteredSongs retrieves a list of songs that match the specified filter criteria.
func (s *SongLibraryService) ReadFilteredSongs(filter *models.Filter) ([]models.Song, error) {
	s.log.Info("reading songs using filter")
	return s.SingStorage.ReadFilteredSongs(filter)
}

// ReadText retrieves a subset of song verses based on the start index and count.
func (s *SongLibraryService) ReadVerse(id, start, count int) ([]*models.Verse, error) {
	// Adjust negative count values to zero.
	if count < 0 {
		count = 0
	}
	// Convert the start index to zero-based indexing.
	start--

	s.log.Info("reading text of the song by id")

	// Retrieve the song from the database by its ID.
	song, err := s.SingStorage.ReadByID(id)
	if err != nil {
		return nil, err
	}

	// Split the song's text into verses.
	verses := strings.Split(song.Text, "\n\n")
	s.log.Info("retrieved verses", slog.Any("verses", verses))

	// Check if the requested range of verses exceeds the available verses.
	if start+count > len(verses) {
		return nil, ErrVerseOutOfBound
	}

	res := make([]*models.Verse, 0, count)

	for _, verse := range verses[start : start+count] {
		res = append(res, models.NewVerse(verse))
	}

	// Return the requested range of verses as a string.
	return res, nil
}

// UpdateSong updates the details of an existing song in the database.
func (s *SongLibraryService) UpdateSong(song *models.Song) error {
	s.log.Info("updating song information")
	return s.SingStorage.UpdateSong(song)
}

// DeleteSong deletes a song from the database by its ID.
func (s *SongLibraryService) DeleteSong(id int) error {
	s.log.Info("deleting song information")
	return s.SingStorage.DeleteSong(id)
}

// ReadByID retrieves all details about a song by its ID.
func (s *SongLibraryService) ReadByID(id int) (*models.Song, error) {
	s.log.Info("retrieving song information by id")
	return s.SingStorage.ReadByID(id)
}
