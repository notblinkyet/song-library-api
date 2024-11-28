package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/notblinkyet/song-library-api/internal/models"
)

var (
	ErrBadRequest     = errors.New("bad request")
	ErrInternalServer = errors.New("internal server error")
)

type ApiClient struct {
	ApiURL string
}

func NewApiClient(url string) *ApiClient {
	return &ApiClient{
		ApiURL: url,
	}
}

func (a *ApiClient) GetMoreAboutSong(req *models.CreateSongRequest) (*models.Song, error) {
	const op = "api.GetMoreAboutSong"
	resp, err := http.Get(fmt.Sprintf("%s/info?song=%s&group=%s", a.ApiURL, req.Title, req.Group))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 400 {
			return nil, fmt.Errorf("%s: %w", op, ErrBadRequest)
		}
		return nil, fmt.Errorf("%s: %w", op, ErrInternalServer)
	}

	var res models.Song

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	res.Group = req.Group
	res.Title = req.Title
	return &res, nil
}
