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
)

//type FakeAuthRepository struct{}
type FakeAuthRepository struct{}

func (far *FakeAuthRepository) SignUp(form model.Form) error {
	return nil
}

func (far *FakeAuthRepository) Login(form model.Form) (*model.User, error) {
	var user model.User

	user.ID = 1
	user.CreatedAt = time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local)
	user.UpdatedAt = time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local)
	//user.DeletedAt = null
	user.Password = "uuuuuu"
	user.Name = ""
	user.Email = "u@u"
	user.Age = 0
	user.Gender = 0
	user.ImageUrl = ""
	user.FavoriteMusicAge = 0
	user.FavoriteArtist = ""
	user.Comment = ""
	//user.Followings = []
	//user.Bookmarkings = []

	return &user, nil
}

func TestSignUpHandler(t *testing.T) {
	// テスト用の JSON ボディ作成
	b, err := json.Marshal(model.Form{Email: "test@test", Password: "testtest"})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}

	//headerをセット
	req.Header.Set("Content-Type", "application/json")

	// テスト用のレスポンス作成
	res := httptest.NewRecorder()
	f := &AuthController{AuthInteractor: usecases.AuthInteractor{
		AuthRepository: &FakeAuthRepository{},
	}}
	f.SignUpHandler(res, req)

	//actual := SignUp(model.Form{Email: "test@test", Password: "testtest"})

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}
}

func TestLoginHandler(t *testing.T) {
	// テスト用の JSON ボディ作成
	b, err := json.Marshal(model.Form{Email: "u@u", Password: "uuuuuu"})
	if err != nil {
		t.Fatal(err)
	}

	// テスト用のリクエスト作成
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(b))

	//headerをセット
	req.Header.Set("Content-Type", "application/json")

	// テスト用のレスポンス作成
	res := httptest.NewRecorder()

	//f := &AuthController{AuthInteractor: usecases.AuthInteractor{
	//	AuthRepository: &FakeAuthRepository{},
	//}}
	f := &AuthController{AuthInteractor: usecases.AuthInteractor{
		AuthRepository: &FakeAuthRepository{},
	}}
	f.LoginHandler(res, req)

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}
	expected := `{"id":1,"createdAt":"2020-06-01T09:00:00+09:00","updatedAt":"2020-06-01T09:00:00+09:00","deletedAt":null,"name":"","email":"u@u","age":0,"gender":0,"imageUrl":"","favoriteMusicAge":0,"favoriteArtist":"","comment":"","followings":[],"bookmarkings":[]}`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			res.Body.String(), expected)
	}
}
