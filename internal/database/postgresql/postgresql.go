package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/notblinkyet/song-library-api/internal/config"
)

var (
	ErrNoAffectedRows = errors.New("no affected row")
	ErrNotFound       = errors.New("not found")
)

type PostgreSQL struct {
	pool *pgxpool.Pool
}

func NewPostgreSQL(config *config.Config) (*PostgreSQL, error) {
	const op = "postgresql.New"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", config.DbUser, config.DbPassword,
		config.DbHost, config.DbPort, config.DbName)
	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &PostgreSQL{
		pool: pool,
	}, nil
}

func (p *PostgreSQL) Close() {
	p.pool.Close()
}
