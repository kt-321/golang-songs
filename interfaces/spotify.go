package interfaces

import (
	"encoding/json"
	"golang-songs/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func GetTracks(token string, title string) (*model.Response, error) {

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

	//ヘッダに」アクセストークン入れている
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tracks model.Tracks

	if err := json.Unmarshal(b, &tracks); err != nil {
		return nil, err
	}

	response := &model.Response{
		Tracks: tracks,
	}

	return response, nil
}
