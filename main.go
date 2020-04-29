package main

import (
	"golang-songs/model"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"encoding/json"
	"fmt"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/davecgh/go-spew/spew"
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

//MySQLへの接続
func gormConnect() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println(".envファイルの読み込み失敗")
	}

	mysqlConfig := os.Getenv("mysqlConfig")
	db, err := gorm.Open("mysql", mysqlConfig)
	if err != nil {
		fmt.Println(err)
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

	dec := json.NewDecoder(r.Body)
	var d Form
	dec.Decode(&d)

	email := d.Email
	password := d.Password

	if email == "" {
		error := model.Error{}
		error.Message = "Emailは必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	if password == "" {
		error := model.Error{}
		error.Message = "パスワードは必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	// dump も出せる
	fmt.Println("---------------------")
	spew.Dump(user)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		fmt.Println(err)
	}

	user.Email = email
	user.Password = string(hash)
	password = string(hash)

	db := gormConnect()
	defer db.Close()
	if err := db.Create(&model.User{Email: email, Password: password}).Error; err != nil {
		error := model.Error{}
		error.Message = "アカウントの作成に失敗しました"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	user.Password = ""
	w.Header().Set("Content-Type", "application/json")

	v, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(v)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	user := model.User{}

	jwt := model.JWT{}

	dec := json.NewDecoder(r.Body)
	var d Form
	dec.Decode(&d)

	email := d.Email
	password := d.Password

	if email == "" {
		error := model.Error{}
		error.Message = "Email は必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	if password == "" {
		error := model.Error{}
		error.Message = "パスワードは必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
	}

	user.Email = email
	user.Password = password

	db := gormConnect()
	defer db.Close()

	userData := model.User{}
	row := db.Where("email = ?", user.Email).Find(&userData)
	if err := db.Where("email = ?", user.Email).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	if _, err := json.Marshal(row); err != nil {
		error := model.Error{}
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	passwordData := userData.Password

	err := bcrypt.CompareHashAndPassword([]byte(passwordData), []byte(password))
	if err != nil {
		error := model.Error{}
		error.Message = "無効なパスワードです。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	//トークン作成
	token, err := createToken(user)
	if err != nil {
		error := model.Error{}
		error.Message = "トークンの作成に失敗しました"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	w.WriteHeader(http.StatusOK)
	jwt.Token = token

	v2, err := json.Marshal(jwt)
	if err != nil {
		fmt.Println(err)
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
		"iss":   "__init__", // JWT の発行者が入る(文字列(__init__)は任意)
	})

	//Dumpを吐く
	spew.Dump(token)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println(err)
	}

	return tokenString, nil
}

//JWT認証のテスト 成功
var TestHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	post := "test"
	json.NewEncoder(w).Encode(post)
})

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/signup", SignUpHandler).Methods("POST")
	r.HandleFunc("/api/login", LoginHandler).Methods("POST")
	//JWT認証のテスト
	r.Handle("/api/test", JwtMiddleware.Handler(TestHandler)).Methods("GET")
	if err := http.ListenAndServe(":8081", r); err != nil {
		fmt.Println(err)
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
