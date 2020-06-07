package main

import (
	"golang-songs/infrastructure"
	"golang-songs/model"
	"log"
	"os"

	"github.com/pkg/errors"

	"github.com/jinzhu/gorm"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".envファイルの読み込み失敗")
	}

	mysqlConfig := os.Getenv("mysqlConfig")

	db, err := gorm.Open("mysql", mysqlConfig)
	if err != nil {
		log.Println(err)
	}

	infrastructure.Dispatch(db)

	db.DB().SetMaxIdleConns(10)
	defer db.Close()
}

// Parse は jwt トークンから元になった認証情報を取り出す。
func Parse(signedString string) (*model.Auth, error) {
	secret := os.Getenv("SIGNINGKEY")

	token, err := jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.Errorf("unexpected signing method: %v", token.Header)
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
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
