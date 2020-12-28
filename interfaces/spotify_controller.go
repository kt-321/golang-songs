package interfaces

import (
	"context"
	"encoding/json"
	"golang-songs/model"
	"golang-songs/usecases"
	"net/http"

	"golang.org/x/oauth2"

	"os"

	"github.com/jinzhu/gorm"
)

type SpotifyController struct {
	SpotifyInteractor usecases.SpotifyInteractor
}

func NewSpotifyController(DB *gorm.DB) *SpotifyController {
	return &SpotifyController{
		SpotifyInteractor: usecases.SpotifyInteractor{
			SpotifyRepository: &SpotifyRepository{
				DB: DB,
			},
		},
	}
}

// SpotifyAPIのリダイレクトURLを返す.
func (spc *SpotifyController) GetRedirectURLHandler(w http.ResponseWriter, r *http.Request) {
	config := oauth2.Config{
		ClientID:     os.Getenv("client_id"),
		ClientSecret: os.Getenv("client_secret"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
		RedirectURL: os.Getenv("redirect_url"),
		Scopes:      []string{},
	}

	url := config.AuthCodeURL("state")

	w.Header().Set("Content-Type", "application/json")

	// Encodeを用いたJson変換.
	encoder := json.NewEncoder(w)
	// 自動エスケープを無効に.
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(url); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 6)

		return
	}
}

// SpotifyAPIのトークンを取得して返す.
func (spc *SpotifyController) GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	var d model.Code

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 17)

		return
	}

	config := oauth2.Config{
		ClientID:     os.Getenv("client_id"),
		ClientSecret: os.Getenv("client_secret"),
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://accounts.spotify.com/api/token",
		},
		RedirectURL: os.Getenv("redirect_url"),
		Scopes:      []string{},
	}

	token, err := config.Exchange(context.TODO(), d.Code)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 34)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	v, err := json.Marshal(token.AccessToken)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 6)

		return
	}

	if _, err := w.Write(v); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 35)

		return
	}
}

// SpotifyAPIにより曲を検索して取得する.
func (spc *SpotifyController) GetTracksHandler(w http.ResponseWriter, r *http.Request) {
	var d model.SearchTitle

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 17)

		return
	}

	if d.Token == "" {
		errorInResponse(w, http.StatusInternalServerError, 36)

		return
	}

	// トラック（曲）検索.
	tracks, err := spc.SpotifyInteractor.GetTracks(d.Token, d.Title)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 37)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	v, err := json.Marshal(tracks)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 6)

		return
	}

	if _, err := w.Write(v); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 37)

		return
	}
}
