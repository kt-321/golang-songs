package main

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"golang-songs/model"
	"os"

	"encoding/json"
	"fmt"
	"log"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/davecgh/go-spew/spew"
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
)

func gormConnect() *gorm.DB {
	db, err := gorm.Open("mysql", "root:@/golang_songs?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	return db
}

// レスポンスにエラーを突っ込んで、返却するメソッド
func errorInResponse(w http.ResponseWriter, status int, error model.Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
	return
}

type Form struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	error := model.Error{}

	dec := json.NewDecoder(r.Body)
	var d Form
	dec.Decode(&d)

	email := d.Email
	password := d.Password

	if email == "" {
		error.Message = "Emailは必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	if password == "" {
		error.Message = "パスワードは必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	// dump も出せる
	fmt.Println("---------------------")
	spew.Dump(user)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Fatal(err)
	}

	user.Email = email
	user.Password = string(hash)
	password = string(hash)

	db := gormConnect()
	defer db.Close()
	db.Create(&model.User{Email: email, Password: password})

	user.Password = ""
	w.Header().Set("Content-Type", "application/json")

	v, err := json.Marshal(user)
	if err != nil {
		println(string(v))
	}
	w.Write(v)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	user2 := model.User{}

	error := model.Error{}
	jwt := model.JWT{}

	dec := json.NewDecoder(r.Body)
	var d Form
	dec.Decode(&d)

	email := d.Email
	password := d.Password


	if email == "" {
		error.Message = "Email は,必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	if password == "" {
		error.Message = "パスワードは、必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
	}

	user.Email = email
	user.Password = password


	db := gormConnect()
	defer db.Close()

	row := db.Where("email = ?", user.Email).Find(&user2)

	_, err := json.Marshal(row)
	if err != nil {
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	password2 := user2.Password

	err = bcrypt.CompareHashAndPassword([]byte(password2), []byte(password))

	if err != nil {
		error.Message = "無効なパスワードです。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	//トークン作成
	token, err := createToken(user)

	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	jwt.Token = token

	v2, err := json.Marshal(jwt)
	if err != nil {
		println(string(v2))
	}
	w.Write(v2)
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
		"iss": "__init__", // JWT の発行者が入る(文字列(__init__)は任意)
	})


	//Dumpを吐く
	spew.Dump(token)

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}

//JWT認証のテスト 成功
var TestHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	post := "test"
	json.NewEncoder(w).Encode(post)
})


func main() {
	db := gormConnect()
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/api/signup", SignUpHandler).Methods("POST")
	r.HandleFunc("/api/login", LoginHandler).Methods("POST")
	//JWT認証のテスト
	r.Handle("/api/test", JwtMiddleware.Handler(TestHandler)).Methods("GET")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatal(err)
	}
}

// JwtMiddleware check token
var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		secret := os.Getenv("SIGNINGKEY")
		return []byte(secret), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})