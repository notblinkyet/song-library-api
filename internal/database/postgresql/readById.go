package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/notblinkyet/song-library-api/internal/models"
)

func (p PostgreSQL) ReadByID(id int) (*models.Song, error) {
	const op = "postgresql.ReadText"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var song models.Song

	query := `SELECT s.id, s.title, g.name, s.release_date, s.song_text, s.link 
		FROM s songs JOIN g groups ON g.id = s.group_id
		WHERE id=$1
	`
	err := p.pool.QueryRow(ctx, query, &id).Scan(&song.ID, &song.Title, &song.Group,
		&song.ReleaseDate, &song.Text, &song.Link)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &song, nil
}
