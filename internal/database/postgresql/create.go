package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/notblinkyet/song-library-api/internal/models"
)

func (p PostgreSQL) CreateSong(song *models.Song) (int, error) {
	const op = "postgresql.CreateSong"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var id, groupID int

	query := `
		SELECT id from groups WHERE name = $1
	`

	err := p.pool.QueryRow(ctx, query, &song.Group).Scan(&groupID)

	if err == pgx.ErrNoRows {
		query = `
			INSERT INTO groups(name)
			VALUES($1) RETURNING id;
		`
		err = p.pool.QueryRow(ctx, query, &song.Group).Scan(&groupID)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	} else if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	query = "INSERT INTO songs (title, group_id, release_date, song_text, link) VALUES ($1, $2, $3, $4, $5) RETURNING id;"

	err = p.pool.QueryRow(ctx, query, &song.Title, &groupID, &song.ReleaseDate, &song.Text, &song.Link).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return id, nil
}
