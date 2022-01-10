package userQuery

import (
	"encoding/json"
	"golang-songs/interfaces"
	"golang-songs/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/gorilla/mux"
)

type FakeDataAccessor struct{}

func (fur *FakeDataAccessor) GetAllUsers() (*[]model.User, error) {
	user1 := model.User{
		ID:               1,
		CreatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		UpdatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		Password:         "aaaaaa",
		Name:             "",
		Email:            "a@test.co.jp",
		Age:              0,
		Gender:           0,
		ImageUrl:         "",
		FavoriteMusicAge: 0,
		FavoriteArtist:   "",
		Comment:          "",
	}

	user2 := model.User{
		ID:               2,
		CreatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		UpdatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		Password:         "iiiiii",
		Name:             "",
		Email:            "i@test.co.jp",
		Age:              0,
		Gender:           0,
		ImageUrl:         "",
		FavoriteMusicAge: 0,
		FavoriteArtist:   "",
		Comment:          "",
	}

	users := []model.User{user1, user2}

	return &users, nil
}

func (fur *FakeDataAccessor) GetUserInfoByID(userID int) (*getUserInfoByIDRes, error) {
	user := getUserInfoByIDRes{
		ID:               1,
		CreatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		UpdatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		Name:             "",
		Email:            "a@test.co.jp",
		Age:              0,
		Gender:           0,
		ImageUrl:         "",
		FavoriteMusicAge: 0,
		FavoriteArtist:   "",
		Comment:          "",
	}

	return &user, nil
}

func (fur *FakeDataAccessor) GetUserInfoByEmail(userEmail string) (*getUserInfoByEmailRes, error) {
	user := getUserInfoByEmailRes{
		ID:               1,
		CreatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		UpdatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		Name:             "",
		Email:            "a@test.co.jp",
		Age:              0,
		Gender:           0,
		ImageUrl:         "",
		FavoriteMusicAge: 0,
		FavoriteArtist:   "",
		Comment:          "",
	}

	return &user, nil
}

func (fur *FakeDataAccessor) GetBookmarkings(userID int) (*getBookmarkingsRes, error) {
	res := &getBookmarkingsRes{
		Bookmarkings: []*model.Song{
			{
				ID:             1,
				CreatedAt:      time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
				UpdatedAt:      time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
				DeletedAt:      nil,
				Title:          "song1",
				Artist:         "artist1",
				MusicAge:       2010,
				Image:          "image1",
				Video:          "video1",
				Album:          "album1",
				Description:    "description1",
				SpotifyTrackId: "spotifySong1",
				UserID:         2,
			},
		},
	}

	return res, nil
}

func (fur *FakeDataAccessor) GetFollowees(userID int) (*getFolloweesRes, error) {
	res := &getFolloweesRes{
		Followees: []*model.User{
			{
				ID:               2,
				CreatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
				UpdatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
				Name:             "followee",
				Email:            "a@test.co.jp",
				Age:              20,
				Gender:           1,
				ImageUrl:         "image1",
				FavoriteMusicAge: 1990,
				FavoriteArtist:   "artist1",
				Comment:          "comment",
			},
		},
	}

	return res, nil
}

//全てのユーザーの情報を返すハンドラのテスト.
func TestAllUsersHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/users", nil)

	// headerをセット.
	req.Header.Set("Content-Type", "application/json")

	// リクエストユーザー作成.
	user := model.User{Email: "a@test.co.jp", Password: "aaaaaa"}

	// トークン作成.
	token, err := interfaces.CreateToken(user)
	if err != nil {
		t.Fatal("トークンの作成に失敗しました")
	}

	jointToken := "Bearer" + " " + token

	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成.
	res := httptest.NewRecorder()

	fakeUserServer := &userQueryServer{
		usecase: usecase{
			da: &FakeDataAccessor{},
		},
	}
	fakeUserServer.GetAllUsers(res, req)

	// レスポンスのステータスコードのテスト.
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}

	// レスポンスボディをDecode.
	var p []model.User
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		t.Fatal("レスポンスボディのデコードに失敗しました。")
	}

	user1 := model.User{
		ID:               1,
		CreatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		UpdatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		Name:             "",
		Email:            "a@test.co.jp",
		Age:              0,
		Gender:           0,
		ImageUrl:         "",
		FavoriteMusicAge: 0,
		FavoriteArtist:   "",
		Comment:          "",
		Bookmarkings:     nil,
		Followings:       nil,
	}

	user2 := model.User{
		ID:               2,
		CreatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		UpdatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		Name:             "",
		Email:            "i@test.co.jp",
		Age:              0,
		Gender:           0,
		ImageUrl:         "",
		FavoriteMusicAge: 0,
		FavoriteArtist:   "",
		Comment:          "",
		Bookmarkings:     nil,
		Followings:       nil,
	}

	// 期待値(アサート用の構造体).
	expected := []model.User{user1, user2}

	// レスポンスのボディが期待通りか確認.
	if diff := cmp.Diff(p, expected); diff != "" {
		t.Errorf("handler returned unexpected body: %v",
			diff)
	}
}

// リクエストユーザーの情報を返すハンドラのテスト.
func TestGetAuthUserHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/user", nil)

	// headerをセット.
	req.Header.Set("Content-Type", "application/json")

	// リクエストユーザー作成.
	user := model.User{Email: "a@test.co.jp", Password: "aaaaaa"}

	// トークン作成.
	token, err := interfaces.CreateToken(user)
	if err != nil {
		t.Fatal("トークンの作成に失敗しました")
	}

	jointToken := "Bearer" + " " + token

	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成.
	res := httptest.NewRecorder()

	fakeUserServer := &userQueryServer{
		usecase: usecase{
			da: &FakeDataAccessor{},
		},
	}
	fakeUserServer.GetAuthUser(res, req)

	// レスポンスのステータスコードのテスト.
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}

	// レスポンスボディをDecode.
	var p model.User
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		t.Fatal("レスポンスボディのデコードに失敗しました。")
	}

	// 期待値(アサート用の構造体).
	expected := model.User{
		ID:               1,
		CreatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		UpdatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		Name:             "",
		Email:            "a@test.co.jp",
		Age:              0,
		Gender:           0,
		ImageUrl:         "",
		FavoriteMusicAge: 0,
		FavoriteArtist:   "",
		Comment:          "",
		Bookmarkings: []*model.Song{
			{
				ID:             1,
				CreatedAt:      time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
				UpdatedAt:      time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
				DeletedAt:      nil,
				Title:          "song1",
				Artist:         "artist1",
				MusicAge:       2010,
				Image:          "image1",
				Video:          "video1",
				Album:          "album1",
				Description:    "description1",
				SpotifyTrackId: "spotifySong1",
				UserID:         2,
			},
		},
		Followings: []*model.User{
			{
				ID:               2,
				CreatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
				UpdatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
				Name:             "followee",
				Email:            "a@test.co.jp",
				Age:              20,
				Gender:           1,
				ImageUrl:         "image1",
				FavoriteMusicAge: 1990,
				FavoriteArtist:   "artist1",
				Comment:          "comment",
			},
		},
	}

	// レスポンスのボディが期待通りか確認.
	if diff := cmp.Diff(p, expected); diff != "" {
		t.Errorf("handler returned unexpected body: %v",
			diff)
	}
}

// idで指定したユーザーの情報を返すハンドラのテスト.
func TestGetUserHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/user/1", nil)

	// headerをセット.
	req.Header.Set("Content-Type", "application/json")

	// リクエストユーザー作成.
	user := model.User{Email: "a@test.co.jp", Password: "aaaaaa"}

	// トークン作成.
	token, err := interfaces.CreateToken(user)
	if err != nil {
		t.Fatal("トークンの作成に失敗しました")
	}

	jointToken := "Bearer" + " " + token

	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成.
	res := httptest.NewRecorder()

	fakeUserServer := &userQueryServer{
		usecase: usecase{
			da: &FakeDataAccessor{},
		},
	}

	r := mux.NewRouter()
	r.Handle("/api/user/{id}", http.HandlerFunc(fakeUserServer.FindUser)).Methods("GET")
	r.ServeHTTP(res, req)

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}

	// レスポンスボディをDecode.
	var p model.User
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		t.Fatal("レスポンスボディのデコードに失敗しました。")
	}

	// 期待値(アサート用の構造体).
	expected := model.User{
		ID:               1,
		CreatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		UpdatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		Name:             "",
		Email:            "a@test.co.jp",
		Age:              0,
		Gender:           0,
		ImageUrl:         "",
		FavoriteMusicAge: 0,
		FavoriteArtist:   "",
		Comment:          "",
		Bookmarkings: []*model.Song{
			{
				ID:             1,
				CreatedAt:      time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
				UpdatedAt:      time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
				DeletedAt:      nil,
				Title:          "song1",
				Artist:         "artist1",
				MusicAge:       2010,
				Image:          "image1",
				Video:          "video1",
				Album:          "album1",
				Description:    "description1",
				SpotifyTrackId: "spotifySong1",
				UserID:         2,
			},
		},
		Followings: []*model.User{
			{
				ID:               2,
				CreatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
				UpdatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
				Name:             "followee",
				Email:            "a@test.co.jp",
				Age:              20,
				Gender:           1,
				ImageUrl:         "image1",
				FavoriteMusicAge: 1990,
				FavoriteArtist:   "artist1",
				Comment:          "comment",
			},
		},
	}

	// レスポンスのボディが期待通りか確認
	if diff := cmp.Diff(p, expected); diff != "" {
		t.Errorf("handler returned unexpected body: %v",
			diff)
	}
}
