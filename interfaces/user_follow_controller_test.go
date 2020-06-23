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

func TestFollowUserHandler(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/user/2/follow", nil)

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
			UserFollowRepository: &FakeUserFollowRepository{},
		},
	}
	//テスト用にルーティング用意
	r := mux.NewRouter()
	r.Handle("/api/user/{id}/follow", http.HandlerFunc(fakeUserFollowController.FollowUserHandler)).Methods("POST")
	r.ServeHTTP(res, req)

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusCreated {
		t.Errorf("invalid code: %d", res.Code)
	}
}

func TestUnfollowUserHandler(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/user/2/unfollow", nil)

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
			UserFollowRepository: &FakeUserFollowRepository{},
		},
	}
	//テスト用にルーティング用意
	r := mux.NewRouter()
	r.Handle("/api/user/{id}/unfollow", http.HandlerFunc(fakeUserFollowController.UnfollowUserHandler)).Methods("POST")
	r.ServeHTTP(res, req)

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusCreated {
		t.Errorf("invalid code: %d", res.Code)
	}
}
