package postgresql

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/notblinkyet/song-library-api/internal/models"
)

func (p PostgreSQL) ReadFilteredSongs(filter *models.Filter) ([]models.Song, error) {
	const op = "postgresql.ReadFilteredSongs"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Валидация ввода
	if filter.Limit < 0 {
		return nil, fmt.Errorf("limit must be non-negative")
	}
	if filter.Offset < 0 {
		return nil, fmt.Errorf("offset must be non-negative")
	}

	var query strings.Builder
	query.WriteString("SELECT id, title, group_name, release_date, song_text, link FROM songs")

	args := make([]any, 0)
	whereClauses := make([]string, 0, 5)
	varCount := 1

	if filter.Title != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("title = $%d", varCount))
		args = append(args, &filter.Title)
		varCount++
	}

	if filter.Group != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("group_name = $%d", varCount))
		args = append(args, &filter.Group)
		varCount++
	}

	if filter.ReleaseDate != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("release_date = $%d", varCount))
		args = append(args, &filter.ReleaseDate)
		varCount++
	}

	if filter.Text != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("song_text LIKE $%d", varCount))
		args = append(args, &filter.Text)
		varCount++
	}

	if filter.Link != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("link = $%d", varCount))
		args = append(args, &filter.Link)
		varCount++
	}

	if len(whereClauses) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(whereClauses, " AND "))
	}

	if filter.Limit > 0 {
		query.WriteString(" LIMIT $")
		query.WriteString(strconv.Itoa(len(args) + 1))
		args = append(args, &filter.Limit)
	}

	if filter.Offset > 0 {
		query.WriteString(" OFFSET $")
		query.WriteString(strconv.Itoa(len(args) + 1))
		args = append(args, &filter.Offset)
	}

	query.WriteString(";")

	rows, err := p.pool.Query(ctx, query.String(), args...)

	fmt.Println(query.String(), args)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var songs []models.Song

	for rows.Next() {
		var song models.Song
		err = rows.Scan(&song.ID, &song.Title, &song.Group, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		songs = append(songs, song)
	}

	return songs, nil
}
