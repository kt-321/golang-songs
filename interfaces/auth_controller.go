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

// ユーザー登録.
func (ac *AuthController) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var d model.Form
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		errorInResponse(w, http.StatusInternalServerError, DecodeError)

		return
	}

	if d.Email == "" {
		errorInResponse(w, http.StatusBadRequest, RequiredEmailError)

		return
	}

	if d.Password == "" {
		errorInResponse(w, http.StatusBadRequest, RequiredPasswordError)

		return
	}

	email, password := d.Email, d.Password

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		errorInResponse(w, http.StatusBadRequest, InvalidPasswordError)

		return
	}

	user := model.User{Email: email, Password: string(hash)}

	d.Password = string(hash)

	err = ac.AuthInteractor.SignUp(d)

	if err != nil {
		errorInResponse(w, http.StatusUnauthorized, CreateAccountError)

		return
	}

	user.Password = ""

	w.Header().Set("Content-Type", "application/json")

	v, err := json.Marshal(user)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, JsonEncodeError)

		return
	}

	if _, err := w.Write(v); err != nil {
		errorInResponse(w, http.StatusInternalServerError, GetUserDetailError)

		return
	}
}

// ログイン.
func (ac *AuthController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var d model.Form

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		errorInResponse(w, http.StatusInternalServerError, DecodeError)

		return
	}

	if d.Email == "" {
		errorInResponse(w, http.StatusBadRequest, RequiredEmailError)

		return
	}

	if d.Password == "" {
		errorInResponse(w, http.StatusBadRequest, RequiredPasswordError)

		return
	}

	email, password := d.Email, d.Password

	userData, err := ac.AuthInteractor.Login(d)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, GetUserDetailError)

		return
	}

	passwordData := userData.Password

	err = bcrypt.CompareHashAndPassword([]byte(passwordData), []byte(password))

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, InvalidPasswordError)

		return
	}

	user := model.User{Email: email, Password: password}

	// トークン作成
	token, err := createToken(user)

	if err != nil {
		errorInResponse(w, http.StatusUnauthorized, CreateTokenError)

		return
	}

	var jwt model.JWT

	w.WriteHeader(http.StatusOK)

	jwt.Token = token

	v, err := json.Marshal(jwt)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, JsonEncodeError)

		return
	}

	if _, err := w.Write(v); err != nil {
		errorInResponse(w, http.StatusInternalServerError, GetJwtTokenError)

		return
	}
}

// JWT.
func createToken(user model.User) (string, error) {
	var err error

	secret := os.Getenv("SIGNINGKEY")

	// Token を作成.
	// jwt -> JSON Web Token - JSON をセキュアにやり取りするための仕様
	// jwtの構造 -> {Base64 encoded Header}.{Base64 encoded Payload}.{Signature}
	// HS254 -> 証明生成用(https://ja.wikipedia.org/wiki/JSON_Web_Token)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss":   "__init__", // JWT の発行者が入る(文字列(__init__)は任意)
	})

	// Dumpを吐く.
	spew.Dump(token)

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return tokenString, fmt.Errorf("failed to change token to string: %v", err)
	}

	return tokenString, nil
}
