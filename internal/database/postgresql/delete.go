package postgresql

import (
	"context"
	"fmt"
	"time"
)

func (p PostgreSQL) DeleteSong(id int) error {
	const op = "postgresql.DeleteSong"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM songs WHERE id = $1;"

	commandTag, err := p.pool.Exec(ctx, query, &id)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
