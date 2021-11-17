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

type FakeUserRepository struct{}

func (fur *FakeUserRepository) FindAll() (*[]model.User, error) {
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

func (fur *FakeUserRepository) GetUser(userEmail string) (*model.User, error) {
	user := model.User{
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

	return &user, nil
}

func (fur *FakeUserRepository) FindByID(userID int) (*model.User, error) {
	user := model.User{
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

	return &user, nil
}

func (fur *FakeUserRepository) Update(userID int, p model.User) error {
	return nil
}

// idで指定したユーザーの情報を返すハンドラのテスト.
func TestGetUserHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/user/1", nil)

	// headerをセット.
	req.Header.Set("Content-Type", "application/json")

	// リクエストユーザー作成.
	user := model.User{Email: "a@test.co.jp", Password: "aaaaaa"}

	// トークン作成.
	token, err := CreateToken(user)
	if err != nil {
		t.Fatal("トークンの作成に失敗しました")
	}

	jointToken := "Bearer" + " " + token

	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成.
	res := httptest.NewRecorder()

	fakeUserController := &UserController{
		UserInteractor: usecases.UserInteractor{
			UserRepository: &FakeUserRepository{},
		},
	}

	r := mux.NewRouter()
	r.Handle("/api/user/{id}", http.HandlerFunc(fakeUserController.GetUserHandler)).Methods("GET")
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
	}

	// レスポンスのボディが期待通りか確認
	if diff := cmp.Diff(p, expected); diff != "" {
		t.Errorf("handler returned unexpected body: %v",
			diff)
	}
}

// リクエストユーザーの情報を返すハンドラのテスト.
func TestUserHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/user", nil)

	// headerをセット.
	req.Header.Set("Content-Type", "application/json")

	// リクエストユーザー作成.
	user := model.User{Email: "a@test.co.jp", Password: "aaaaaa"}

	// トークン作成.
	token, err := CreateToken(user)
	if err != nil {
		t.Fatal("トークンの作成に失敗しました")
	}

	jointToken := "Bearer" + " " + token

	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成.
	res := httptest.NewRecorder()

	f := &UserController{UserInteractor: usecases.UserInteractor{
		UserRepository: &FakeUserRepository{},
	}}
	f.UserHandler(res, req)

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
	}

	// レスポンスのボディが期待通りか確認.
	if diff := cmp.Diff(p, expected); diff != "" {
		t.Errorf("handler returned unexpected body: %v",
			diff)
	}
}

// 全てのユーザーの情報を返すハンドラのテスト.
func TestAllUsersHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/users", nil)

	// headerをセット.
	req.Header.Set("Content-Type", "application/json")

	// リクエストユーザー作成.
	user := model.User{Email: "a@test.co.jp", Password: "aaaaaa"}

	// トークン作成.
	token, err := CreateToken(user)
	if err != nil {
		t.Fatal("トークンの作成に失敗しました")
	}

	jointToken := "Bearer" + " " + token

	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成.
	res := httptest.NewRecorder()

	f := &UserController{UserInteractor: usecases.UserInteractor{
		UserRepository: &FakeUserRepository{},
	}}
	f.AllUsersHandler(res, req)

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
	}

	// 期待値(アサート用の構造体).
	expected := []model.User{user1, user2}

	// レスポンスのボディが期待通りか確認.
	if diff := cmp.Diff(p, expected); diff != "" {
		t.Errorf("handler returned unexpected body: %v",
			diff)
	}
}

// idで指定したユーザーの情報を更新するハンドラのテスト.
func TestUpdateUserHandler(t *testing.T) {
	// テスト用の JSON ボディ作成.
	b, err := json.Marshal(model.User{Email: "hoge@test.co.jp", Name: "hogehoge", Age: 0, Gender: 0, FavoriteMusicAge: 0, FavoriteArtist: "椎名林檎", Comment: "テストユーザーです。"})
	if err != nil {
		t.Fatal(err)
	}

	// テスト用のリクエスト作成.
	req := httptest.NewRequest("PUT", "/api/user/2/update", bytes.NewBuffer(b))

	// headerをセット.
	req.Header.Set("Content-Type", "application/json")

	// リクエストユーザー作成.
	user := model.User{Email: "a@test.co.jp", Password: "aaaaaa"}

	// トークン作成.
	token, err := CreateToken(user)
	if err != nil {
		t.Fatal("トークンの作成に失敗しました")
	}

	jointToken := "Bearer" + " " + token

	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成.
	res := httptest.NewRecorder()

	fakeUserController := &UserController{
		UserInteractor: usecases.UserInteractor{
			UserRepository: &FakeUserRepository{},
		},
	}

	r := mux.NewRouter()
	r.Handle("/api/user/{id}/update", http.HandlerFunc(fakeUserController.UpdateUserHandler)).Methods("PUT")
	r.ServeHTTP(res, req)

	// レスポンスのステータスコードのテスト(204).
	if res.Code != http.StatusNoContent {
		t.Errorf("invalid code: %d", res.Code)
	}
}
