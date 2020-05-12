package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"golang-songs/model"
	"net/http"
	"os"

	"golang-songs/service"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

// レスポンスにエラーを突っ込んで、返却するメソッド
func errorInResponse(w http.ResponseWriter, status int, error model.Error) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(error); err != nil {
		var error model.Error
		error.Message = "リクエストボディのデコードに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}

func GetRedirectURL(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		var error model.Error
		error.Message = ".envの読み込みに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	config := oauth2.Config{
		ClientID:     os.Getenv("client_id"),
		ClientSecret: os.Getenv("client_secret"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},

		RedirectURL: "http://localhost:3000/spotify/songs",
		Scopes:      []string{},
	}

	url := config.AuthCodeURL("state")

	w.Header().Set("Content-Type", "application/json")

	// Encodeを用いたJson変換
	encoder := json.NewEncoder(w)
	//自動エスケープを無効に
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(url); err != nil {
		var error model.Error
		error.Message = "JSONへの変換に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var d model.Code
	if err := dec.Decode(&d); err != nil {
		var error model.Error
		error.Message = "リクエストボディのデコードに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	err := godotenv.Load()
	if err != nil {
		var error model.Error
		error.Message = ".envの読み込みに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	config := oauth2.Config{
		ClientID:     os.Getenv("client_id"),
		ClientSecret: os.Getenv("client_secret"),
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://accounts.spotify.com/api/token",
		},
		RedirectURL: "http://localhost:3000/spotify/songs",
		Scopes:      []string{},
	}

	token, err := config.Exchange(context.TODO(), d.Code)
	if err != nil {
		var error model.Error
		error.Message = "トークンの取得に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	v, err := json.Marshal(token.AccessToken)
	if err != nil {
		var error model.Error
		error.Message = "JSONへの変換に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	if _, err := w.Write(v); err != nil {
		var error model.Error
		error.Message = "URLの取得に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}

func JSONSafeMarshal(v interface{}, safeEncoding bool) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	if safeEncoding {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}

func GetTracks(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var d model.SearchTitle
	if err := dec.Decode(&d); err != nil {
		var error model.Error
		error.Message = "リクエストボディのデコードに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	if d.Token == "" {
		var error model.Error
		error.Message = "アクセストークンの取得に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	//トラック（曲）検索
	tracks, err := service.GetTracks(d.Token, d.Title)
	if err != nil {
		var error model.Error
		error.Message = "トラックの取得に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	v, err := json.Marshal(tracks)
	if err != nil {
		var error model.Error
		error.Message = "JSONへの変換に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	if _, err := w.Write(v); err != nil {
		var error model.Error
		error.Message = "トラックの取得に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}
