package main

import (
	"fmt"
	"time"

	"github.com/notblinkyet/song-library-api/internal/config"
	"github.com/notblinkyet/song-library-api/internal/database/postgresql"
	"github.com/notblinkyet/song-library-api/internal/models"
)

func main() {
	cfg := config.MustLoadConfig()
	db, err := postgresql.NewPostgreSQL(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	t, err := time.Parse("02.01.2006", "16.07.2006")
	if err != nil {
		panic(err)
	}
	song := models.Song{
		Title:       "Supermassive Black Hole",
		Group:       "Muse",
		ReleaseDate: t,
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	id, err := db.CreateSong(&song)
	if err != nil {
		panic(err)
	}
	fmt.Println(id)

}
