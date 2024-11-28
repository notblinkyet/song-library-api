package http

import (
	"github.com/notblinkyet/song-library-api/internal/models"
)

type SongLibraryService interface {
	Create(req *models.CreateSongRequest) (int, error)
	ReadFilteredSongs(filter *models.Filter) ([]models.Song, error)
	ReadVerse(id, start, count int) ([]*models.Verse, error)
	UpdateSong(song *models.Song) error
	DeleteSong(id int) error
	ReadByID(id int) (*models.Song, error)
}
