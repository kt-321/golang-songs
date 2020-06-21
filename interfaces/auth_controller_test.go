package interfaces

import (
	"bytes"
	"encoding/json"
	"golang-songs/model"
	"golang-songs/usecases"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"

	"golang.org/x/crypto/bcrypt"
)

type FakeAuthRepository struct{}

func (far *FakeAuthRepository) SignUp(form model.Form) error {
	return nil
}

func (far *FakeAuthRepository) Login(form model.Form) (*model.User, error) {
	email := "test@test"
	password := "testtest"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	user := model.User{Email: email, Password: string(hash)}

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

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}
}

func TestLoginHandler(t *testing.T) {
	//テスト用のemail,passwordを準備
	email := "test@test.co.jp"
	password := "testtest"

	// テスト用の JSON ボディ作成
	b, err := json.Marshal(model.Form{Email: email, Password: password})
	if err != nil {
		t.Fatal(err)
	}

	// テスト用のリクエスト作成
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(b))
	//headerをセット
	req.Header.Set("Content-Type", "application/json")

	// テスト用のレスポンス作成
	res := httptest.NewRecorder()

	f := &AuthController{AuthInteractor: usecases.AuthInteractor{
		AuthRepository: &FakeAuthRepository{},
	}}
	f.LoginHandler(res, req)

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}

	user := model.User{Email: email, Password: password}

	//トークン作成
	token, err := createToken(user)
	if err != nil {
		t.Fatal("トークンの作成に失敗しました")
	}

	//レスポンスボディをDecode
	var p model.JWT
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		t.Fatal("JSONへの変換に失敗しました")
	}

	//期待値(アサート用の構造体)
	var expected model.JWT
	expected.Token = token

	// レスポンスのボディが期待通りか確認
	if diff := cmp.Diff(p, expected); diff != "" {
		t.Errorf("handler returned unexpected body: %v",
			diff)
	}
}
