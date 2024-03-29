package interfaces

import (
	"fmt"
	"golang-songs/model"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// Parse は jwt トークンから元になった認証情報を取り出す.
func Parse(signedString string) (*model.Auth, error) {
	secret := os.Getenv("SIGNINGKEY")

	token, err := jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.Errorf("unexpected signing method: %v", token.Header)
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse string type jwttoken: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.Errorf("not found claims in %s", signedString)
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.Errorf("not found %s in %s", email, signedString)
	}

	return &model.Auth{Email: email}, nil
}

// リクエストヘッダーからユーザーのEmailを取得するメソッド
func GetEmail(r *http.Request) (string, *model.ErrorSet) {
	headerAuthorization := r.Header.Get("Authorization")

	if len(headerAuthorization) == 0 {
		return "", &model.ErrorSet{StatusCode: http.StatusInternalServerError, Message: GetAuthenticationTokenError}
	}

	bearerToken := strings.Split(headerAuthorization, " ")

	if len(bearerToken) < 2 {
		return "", &model.ErrorSet{StatusCode: http.StatusUnauthorized, Message: GetBearerTokenError}
	}

	parsedToken, err := Parse(bearerToken[1])

	if err != nil {
		return "", &model.ErrorSet{StatusCode: http.StatusInternalServerError, Message: ParseAuthenticationCodeError}
	}

	return parsedToken.Email, nil
}

// パスパラメータからターゲットとなるidの値を取得するメソッド
func GetId(r *http.Request) (int, *model.ErrorSet) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		e := errors.New("key'id'が存在しない")
		return 0, &model.ErrorSet{StatusCode: http.StatusBadRequest, Message: GetIdError, Err: &e}
	}

	parsedId, err := strconv.Atoi(id)

	if err != nil {
		e := errors.WithStack(err)
		return 0, &model.ErrorSet{StatusCode: http.StatusInternalServerError, Message: ConvertIdToIntError, Err: &e}
	}

	return parsedId, nil
}

// リクエストユーザーのメールアドレスとターゲットのidを取得するメソッド
func GetEmailAndId(r *http.Request) (string, int, *model.ErrorSet) {
	userEmail, errorSet := GetEmail(r)

	if errorSet != nil {
		return "", 0, errorSet
	}

	id, errorSet := GetId(r)

	if errorSet != nil {
		return "", 0, errorSet
	}

	return userEmail, id, nil
}
