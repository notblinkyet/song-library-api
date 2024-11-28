package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/notblinkyet/song-library-api/internal/models"
)

func (p PostgreSQL) UpdateSong(song *models.Song) error {
	const op = "postgresql.UpdateSong"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "UPDATE songs SET title=$1, group_name=$2, release_date=$3, song_text=$4, link=$5 WHERE id=$6;"

	commandTag, err := p.pool.Exec(ctx, query, &song.Title,
		&song.Group, &song.ReleaseDate, &song.Text, &song.Link, &song.ID)
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
