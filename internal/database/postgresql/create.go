package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/notblinkyet/song-library-api/internal/models"
)

func (p PostgreSQL) CreateSong(song *models.Song) (int, error) {
	const op = "postgresql.CreateSong"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var id int

	query := "INSERT INTO songs (title, group_name, release_date, song_text, link) VALUES ($1, $2, $3, $4, $5) RETURNING id;"

	err := p.pool.QueryRow(ctx, query, &song.Title, &song.Group, &song.ReleaseDate, &song.Text, &song.Link).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return id, nil
}
