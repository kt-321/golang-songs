package interfaces

import (
	"encoding/json"
	"fmt"
	"golang-songs/model"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"

	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	"golang-songs/usecases"
)

type AuthController struct {
	AuthInteractor usecases.AuthInteractor
}

func NewAuthController(DB *gorm.DB) *AuthController {
	return &AuthController{
		AuthInteractor: usecases.AuthInteractor{
			AuthRepository: &AuthRepository{
				DB: DB,
			},
		},
	}
}

//ユーザー登録
func (ac *AuthController) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	dec := json.NewDecoder(r.Body)
	var d model.Form
	if err := dec.Decode(&d); err != nil {
		var error model.Error
		error.Message = "リクエストボディのデコードに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	email := d.Email
	password := d.Password

	if email == "" {
		var error model.Error
		error.Message = "Emailは必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	if password == "" {
		var error model.Error
		error.Message = "パスワードは必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	// dump も出せる
	fmt.Println("---------------------")
	spew.Dump(user)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		var error model.Error
		error.Message = "パスワードの値が不正です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	user.Email = email
	user.Password = string(hash)

	d.Password = string(hash)

	err = ac.AuthInteractor.SignUp(d)
	if err != nil {
		var error model.Error
		error.Message = "アカウントの作成に失敗しました"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	user.Password = ""
	w.Header().Set("Content-Type", "application/json")

	v, err := json.Marshal(user)
	if err != nil {
		var error model.Error
		error.Message = "JSONへの変換に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	if _, err := w.Write(v); err != nil {
		var error model.Error
		error.Message = "ユーザー情報の取得に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}

//ログイン
func (ac *AuthController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User

	dec := json.NewDecoder(r.Body)
	var d model.Form
	if err := dec.Decode(&d); err != nil {
		var error model.Error
		error.Message = "リクエストボディのデコードに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	email := d.Email
	password := d.Password

	if email == "" {
		var error model.Error
		error.Message = "Email は必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	if password == "" {
		var error model.Error
		error.Message = "パスワードは必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
	}

	user.Email = email
	user.Password = password

	userData, err := ac.AuthInteractor.Login(d)
	if err != nil {
		var error model.Error
		error.Message = "無効なパスワードです。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	passwordData := userData.Password

	err = bcrypt.CompareHashAndPassword([]byte(passwordData), []byte(password))
	if err != nil {
		var error model.Error
		error.Message = "無効なパスワードです。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	//トークン作成
	token, err := createToken(user)
	if err != nil {
		var error model.Error
		error.Message = "トークンの作成に失敗しました"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	var jwt model.JWT

	w.WriteHeader(http.StatusOK)
	jwt.Token = token

	v, err := json.Marshal(jwt)
	if err != nil {
		var error model.Error
		error.Message = "JSONへの変換に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	if _, err := w.Write(v); err != nil {
		var error model.Error
		error.Message = "JWTトークンの取得に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}

//JWT
func createToken(user model.User) (string, error) {
	var err error

	secret := os.Getenv("SIGNINGKEY")

	// Token を作成
	// jwt -> JSON Web Token - JSON をセキュアにやり取りするための仕様
	// jwtの構造 -> {Base64 encoded Header}.{Base64 encoded Payload}.{Signature}
	// HS254 -> 証明生成用(https://ja.wikipedia.org/wiki/JSON_Web_Token)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss":   "__init__", // JWT の発行者が入る(文字列(__init__)は任意)
	})

	//Dumpを吐く
	spew.Dump(token)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return tokenString, err
	}

	return tokenString, nil
}
