package interfaces

import (
	"golang-songs/model"
	"golang-songs/usecases"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type FakeBookmarkRepository struct{}

func (fbr *FakeBookmarkRepository) Bookmark(userEmail string, songID int) error {
	return nil
}

func (fbr *FakeBookmarkRepository) RemoveBookmark(userEmail string, songID int) error {
	return nil
}

func TestBookmarkHandler(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/song/1/bookmark", nil)

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
	fakeBookmarkController := &BookmarkController{
		BookmarkInteractor: usecases.BookmarkInteractor{
			BookmarkRepository: &FakeBookmarkRepository{},
		},
	}
	// テスト用にルーティング用意.
	r := mux.NewRouter()
	r.Handle("/api/song/{id}/bookmark", http.HandlerFunc(fakeBookmarkController.BookmarkHandler)).Methods("POST")
	r.ServeHTTP(res, req)

	// レスポンスのステータスコードのテスト.
	if res.Code != http.StatusCreated {
		t.Errorf("invalid code: %d", res.Code)
	}
}

func TestRemoveBookmarkHandler(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/song/1/remove-bookmark", nil)

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
	fakeBookmarkController := &BookmarkController{
		BookmarkInteractor: usecases.BookmarkInteractor{
			BookmarkRepository: &FakeBookmarkRepository{},
		},
	}
	// テスト用にルーティング用意.
	r := mux.NewRouter()
	r.Handle("/api/song/{id}/remove-bookmark", http.HandlerFunc(fakeBookmarkController.RemoveBookmarkHandler)).Methods("POST")
	r.ServeHTTP(res, req)

	// レスポンスのステータスコードのテスト.
	if res.Code != http.StatusCreated {
		t.Errorf("invalid code: %d", res.Code)
	}
}
