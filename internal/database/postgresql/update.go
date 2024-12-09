package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/notblinkyet/song-library-api/internal/models"
)

func (p PostgreSQL) UpdateSong(song *models.Song) error {
	const op = "postgresql.UpdateSong"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var group_id int

	query := `SELECT id FROM groups WHERE name=$1`

	err := p.pool.QueryRow(ctx, query, &song.Group).Scan(&group_id)

	if err == pgx.ErrNoRows {
		query = `
			INSERT INTO groups (name)
			VALUES($1) RETURNING id;
		`
		err = p.pool.QueryRow(ctx, query, &song.Group).Scan(&group_id)
		if err != nil {
			return fmt.Errorf("group not found: %s", song.Group)
		}
	} else if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query = "UPDATE songs SET title=$1, group_id=$2, release_date=$3, song_text=$4, link=$5 WHERE id=$6;"

	commandTag, err := p.pool.Exec(ctx, query, &song.Title,
		&group_id, &song.ReleaseDate, &song.Text, &song.Link, &song.ID)
	if err != nil {
		if err == ErrNoAffectedRows {
			return ErrNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrNoAffectedRows)
	}
	return nil
}
