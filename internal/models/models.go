package models

type Song struct {
	ID          int    `json:"id"`
	Title       string `json:"song"`
	Group       string `json:"group"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type CreateSongRequest struct {
	Title string `json:"song"`
	Group string `json:"group"`
}

type Filter struct {
	Title       string
	Group       string
	ReleaseDate string
	Text        string
	Link        string
	Limit       int
	Offset      int
}

type Verse struct {
	Verse string `json:"verse"`
}

type Id struct {
	Id int `json:"id"`
}

func NewVerse(verse string) *Verse {
	return &Verse{
		Verse: verse,
	}
}

func NewId(id int) *Id {
	return &Id{
		Id: id,
	}
}
