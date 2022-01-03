package interfaces

import (
	"encoding/json"
	"golang-songs/model"
	"net/http"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
)

type SpotifyRepository struct {
	DB *sqlx.DB
}

func (spr *SpotifyRepository) GetTracks(token string, title string) (*model.Response, error) {
	values := url.Values{}

	values.Add("type", "track")
	values.Add("q", title)
	values.Add("market", "JP")
	values.Add("limit", "10")

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/search", nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = values.Encode()

	// ヘッダにアクセストークン入れている.
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var tracks model.Tracks

	if err := json.NewDecoder(resp.Body).Decode(&tracks); err != nil {
		return nil, err
	}

	response := &model.Response{
		Tracks: tracks,
	}

	return response, nil
}
