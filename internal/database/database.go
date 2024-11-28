package database

import "github.com/notblinkyet/song-library-api/internal/models"

type Storage interface {
	ReadFilteredSongs(filter *models.Filter) ([]models.Song, error)
	ReadByID(id int) (*models.Song, error)
	DeleteSong(id int) error
	UpdateSong(song *models.Song) error
	CreateSong(song *models.Song) (int, error)
}
