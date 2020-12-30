package interfaces

import (
	"bytes"
	"encoding/json"
	"golang-songs/model"
	"golang-songs/usecases"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/gorilla/mux"
)

type FakeSongRepository struct{}

func (fsr *FakeSongRepository) FindAll() (*[]model.Song, error) {
	song1 := model.Song{
		ID:             1,
		CreatedAt:      time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		UpdatedAt:      time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		Title:          "title1",
		Artist:         "artist1",
		MusicAge:       0,
		Image:          "",
		Video:          "",
		Album:          "",
		Description:    "",
		SpotifyTrackId: "",
		UserID:         1,
	}

	song2 := model.Song{
		ID:             2,
		CreatedAt:      time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		UpdatedAt:      time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		Title:          "title2",
		Artist:         "artist2",
		MusicAge:       0,
		Image:          "",
		Video:          "",
		Album:          "",
		Description:    "",
		SpotifyTrackId: "",
		UserID:         2,
	}

	songs := []model.Song{song1, song2}

	return &songs, nil
}

func (fsr *FakeSongRepository) FindByID(songID int) (*model.Song, error) {
	song := model.Song{
		ID:             1,
		CreatedAt:      time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		UpdatedAt:      time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		Title:          "title1",
		Artist:         "artist1",
		MusicAge:       0,
		Image:          "",
		Video:          "",
		Album:          "",
		Description:    "",
		SpotifyTrackId: "",
		UserID:         1,
	}

	return &song, nil
}

func (fsr *FakeSongRepository) Save(userEmail string, p model.Song) error {
	return nil
}
func (fsr *FakeSongRepository) UpdateByID(userEmail string, songID int, p model.Song) error {
	return nil
}

func (fsr *FakeSongRepository) DeleteByID(songID int) error {
	return nil
}

func TestAllSongsHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/songs", nil)

	// headerをセット.
	req.Header.Set("Content-Type", "application/json")

	// リクエストユーザー作成.
	user := model.User{Email: "a@test.co.jp", Password: "aaaaaa"}

	// トークン作成.
	token, err := createToken(user)
	if err != nil {
		t.Fatal("トークンの作成に失敗しました")
	}

	jointToken := "Bearer" + " " + token

	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成.
	res := httptest.NewRecorder()

	f := &SongController{SongInteractor: usecases.SongInteractor{
		SongRepository: &FakeSongRepository{},
	}}
	f.AllSongsHandler(res, req)

	// レスポンスのステータスコードのテスト.
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}

	// レスポンスボディをDecode.
	var p []model.Song
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		t.Fatal("レスポンスボディのデコードに失敗しました。")
	}

	song1 := model.Song{
		ID:             1,
		CreatedAt:      time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		UpdatedAt:      time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		Title:          "title1",
		Artist:         "artist1",
		MusicAge:       0,
		Image:          "",
		Video:          "",
		Album:          "",
		Description:    "",
		SpotifyTrackId: "",
		UserID:         1,
	}

	song2 := model.Song{
		ID:             2,
		CreatedAt:      time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		UpdatedAt:      time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		Title:          "title2",
		Artist:         "artist2",
		MusicAge:       0,
		Image:          "",
		Video:          "",
		Album:          "",
		Description:    "",
		SpotifyTrackId: "",
		UserID:         2,
	}

	// 期待値(アサート用の構造体).
	expected := []model.Song{song1, song2}

	// レスポンスのボディが期待通りか確認.
	if diff := cmp.Diff(p, expected); diff != "" {
		t.Errorf("handler returned unexpected body: %v",
			diff)
	}
}

// idで指定した曲の情報を返すハンドラのテスト.
func TestGetSongHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/song/1", nil)

	// headerをセット.
	req.Header.Set("Content-Type", "application/json")

	// リクエストユーザー作成.
	user := model.User{Email: "a@test.co.jp", Password: "aaaaaa"}

	// トークン作成.
	token, err := createToken(user)
	if err != nil {
		t.Fatal("トークンの作成に失敗しました")
	}

	jointToken := "Bearer" + " " + token

	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成.
	res := httptest.NewRecorder()

	// テスト用にコントローラ用意.
	fakeSongController := &SongController{
		SongInteractor: usecases.SongInteractor{
			SongRepository: &FakeSongRepository{},
		},
	}

	// テスト用にルーティング用意.
	r := mux.NewRouter()
	r.Handle("/api/song/{id}", http.HandlerFunc(fakeSongController.GetSongHandler)).Methods("GET")
	r.ServeHTTP(res, req)

	// レスポンスのステータスコードのテスト.
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}

	// レスポンスボディをDecode.
	var p model.Song
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		t.Fatal("レスポンスボディのデコードに失敗しました。。")
	}

	// 期待値(アサート用の構造体).
	expected := model.Song{
		ID:             1,
		CreatedAt:      time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		UpdatedAt:      time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		Title:          "title1",
		Artist:         "artist1",
		MusicAge:       0,
		Image:          "",
		Video:          "",
		Album:          "",
		Description:    "",
		SpotifyTrackId: "",
		UserID:         1,
	}

	// レスポンスのボディが期待通りか確認.
	if diff := cmp.Diff(p, expected); diff != "" {
		t.Errorf("handler returned unexpected body: %v",
			diff)
	}
}

// 新しく曲を追加するハンドラのテスト.
func TestCreateSongHandler(t *testing.T) {
	// テスト用の JSON ボディ作成.
	b, err := json.Marshal(model.Song{Title: "song1", Artist: "artist1", MusicAge: 1980, Image: "", Video: "", Album: "", Description: "テスト曲です。", SpotifyTrackId: "", UserID: 1})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/api/song", bytes.NewBuffer(b))

	// headerをセット.
	req.Header.Set("Content-Type", "application/json")

	// リクエストユーザー作成.
	user := model.User{Email: "a@test.co.jp", Password: "aaaaaa"}

	// トークン作成.
	token, err := createToken(user)
	if err != nil {
		t.Fatal("トークンの作成に失敗しました")
	}

	jointToken := "Bearer" + " " + token

	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成.
	res := httptest.NewRecorder()

	f := &SongController{SongInteractor: usecases.SongInteractor{
		SongRepository: &FakeSongRepository{},
	}}
	f.CreateSongHandler(res, req)

	// レスポンスのステータスコードのテスト(201).
	if res.Code != http.StatusCreated {
		t.Errorf("invalid code: %d", res.Code)
	}
}

// idで指定した曲を更新するハンドラのテスト.
func TestUpdateSongHandler(t *testing.T) {
	// テスト用の JSON ボディ作成
	b, err := json.Marshal(model.Song{Title: "song2", Artist: "artist2", MusicAge: 2000, Image: "", Video: "", Album: "", Description: "テスト曲です。2", SpotifyTrackId: "", UserID: 1})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("PUT", "/api/song/1", bytes.NewBuffer(b))

	// headerをセット.
	req.Header.Set("Content-Type", "application/json")

	// リクエストユーザー作成.
	user := model.User{Email: "a@test.co.jp", Password: "aaaaaa"}

	// トークン作成.
	token, err := createToken(user)
	if err != nil {
		t.Fatal("トークンの作成に失敗しました")
	}

	jointToken := "Bearer" + " " + token

	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成.
	res := httptest.NewRecorder()

	// テスト用にコントローラ用意.
	fakeSongController := &SongController{
		SongInteractor: usecases.SongInteractor{
			SongRepository: &FakeSongRepository{},
		},
	}

	// テスト用にルーティング用意.
	r := mux.NewRouter()
	r.Handle("/api/song/{id}", http.HandlerFunc(fakeSongController.UpdateSongHandler)).Methods("PUT")
	r.ServeHTTP(res, req)

	// レスポンスのステータスコードのテスト(204).
	if res.Code != http.StatusNoContent {
		t.Errorf("invalid code: %d", res.Code)
	}
}

// idで指定した曲を削除するハンドラのテスト.
func TestDeleteSongHandler(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/api/song/1", nil)

	// リクエストユーザー作成.
	user := model.User{Email: "a@test.co.jp", Password: "aaaaaa"}

	// トークン作成.
	token, err := createToken(user)
	if err != nil {
		t.Fatal("トークンの作成に失敗しました")
	}

	jointToken := "Bearer" + " " + token

	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成.
	res := httptest.NewRecorder()

	// テスト用にコントローラ用意.
	fakeSongController := &SongController{
		SongInteractor: usecases.SongInteractor{
			SongRepository: &FakeSongRepository{},
		},
	}

	// テスト用にルーティング用意.
	r := mux.NewRouter()
	r.Handle("/api/song/{id}", http.HandlerFunc(fakeSongController.DeleteSongHandler)).Methods("DELETE")
	r.ServeHTTP(res, req)

	// レスポンスのステータスコードのテスト(204).
	if res.Code != http.StatusNoContent {
		t.Errorf("invalid code: %d", res.Code)
	}
}
