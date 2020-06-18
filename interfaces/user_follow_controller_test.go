package interfaces

import (
	"golang-songs/model"
	"golang-songs/usecases"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type FakeUserFollowRepository struct{}

func (fufr *FakeUserFollowRepository) Follow(requestUserEmail string, targetUserID int) error {
	return nil
}

func (fufr *FakeUserFollowRepository) Unfollow(requestUserEmail string, targetUserID int) error {
	return nil
}

func FollowUserHandler(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/user/{id}/follow", nil)

	//リクエストユーザー作成
	user := model.User{Email: "a@test.co.jp", Password: "aaaaaa"}
	//トークン作成
	token, err := createToken(user)
	if err != nil {
		log.Println("err:", err)
	}
	jointToken := "Bearer" + " " + token
	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成
	res := httptest.NewRecorder()

	//テスト用にコントローラ用意
	fakeUserFollowController := &UserFollowController{
		UserFollowInteractor: usecases.UserFollowInteractor{
			BookmarkRepository: &BookmarkRepository{},
		},
	}
	//テスト用にルーティング用意
	r := mux.NewRouter()
	r.Handle("/api/user/{id}/follow", http.HandlerFunc(fakeBookmarkController.BookmarkHandler)).Methods("POST")
	r.ServeHTTP(res, req)

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}
}

func UnfollowUserHandler(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/user/{id}/unfollow", nil)

	//リクエストユーザー作成
	user := model.User{Email: "a@test.co.jp", Password: "aaaaaa"}
	//トークン作成
	token, err := createToken(user)
	if err != nil {
		log.Println("err:", err)
	}
	jointToken := "Bearer" + " " + token
	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成
	res := httptest.NewRecorder()

	//テスト用にコントローラ用意
	fakeUserFollowController := &BookmarkController{
		BookmarkInteractor: usecases.BookmarkInteractor{
			BookmarkRepository: &BookmarkRepository{},
		},
	}
	//テスト用にルーティング用意
	r := mux.NewRouter()
	r.Handle("/api/user/{id}/unfollow", http.HandlerFunc(fakeBookmarkController.BookmarkHandler)).Methods("POST")
	r.ServeHTTP(res, req)

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}
}
