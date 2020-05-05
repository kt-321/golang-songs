package main

import (
	"encoding/json"
	"fmt"
	"golang-songs/model"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/davecgh/go-spew/spew"
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

//MySQLへの接続
func gormConnect() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println(".envファイルの読み込み失敗")
	}

	mysqlConfig := os.Getenv("mysqlConfig")
	db, err := gorm.Open("mysql", mysqlConfig)
	if err != nil {
		return nil, err
	}
	return db, nil
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
	var user model.User

	dec := json.NewDecoder(r.Body)
	var d Form
	dec.Decode(&d)

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
	password = string(hash)

	db, _ := gormConnect()
	defer db.Close()
	if err := db.Create(&model.User{Email: email, Password: password}).Error; err != nil {
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
	w.Write(v)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User

	var jwt model.JWT

	dec := json.NewDecoder(r.Body)
	var d Form
	dec.Decode(&d)

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

	db, _ := gormConnect()
	defer db.Close()

	var userData model.User
	row := db.Where("email = ?", user.Email).Find(&userData)
	if err := db.Where("email = ?", user.Email).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		var error model.Error
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	if _, err := json.Marshal(row); err != nil {
		var error model.Error
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	passwordData := userData.Password

	err := bcrypt.CompareHashAndPassword([]byte(passwordData), []byte(password))
	if err != nil {
		var error model.Error
		error.Message = "無効なパスワードです。"
		errorInResponse(w, http.StatusUnauthorized, error)
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

	w.WriteHeader(http.StatusOK)
	jwt.Token = token

	v2, err := json.Marshal(jwt)
	if err != nil {
		var error model.Error
		error.Message = "JSONへの変換に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
	w.Write(v2)
}

//リクエストユーザーの情報を返す
var UserHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	headerAuthorization := r.Header.Get("Authorization")
	bearerToken := strings.Split(headerAuthorization, " ")
	authToken := bearerToken[1]

	parsedToken, err := Parse(authToken)
	userEmail := parsedToken.Email

	db, _ := gormConnect()
	defer db.Close()

	var user model.User

	if err := db.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	v, err := json.Marshal(user)
	if err != nil {
		var error model.Error
		error.Message = "JSONへの変換に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	w.Write(v)
})

//全てのユーザーを返す
var AllUsersHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	db, _ := gormConnect()
	defer db.Close()

	allUsers := []model.User{}

	db.Find(&allUsers)

	if err := db.Find(&allUsers).Error; gorm.IsRecordNotFoundError(err) {
		var error model.Error
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	v, err := json.Marshal(allUsers)
	if err != nil {
		var error model.Error
		error.Message = "ユーザー一覧の取得に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
	w.Write(v)
})

var UpdateUserHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	dec := json.NewDecoder(r.Body)
	var d model.User
	dec.Decode(&d)

	email := d.Email
	name := d.Name
	age := d.Age
	gender := d.Gender
	favoriteMusicAge := d.FavoriteMusicAge
	favoriteArtist := d.FavoriteArtist
	comment := d.Comment

	db, _ := gormConnect()
	defer db.Close()

	var user model.User

	if err := db.Model(&user).Where("id = ?", id).Update(model.User{Email: email, Name: name, Age: age, Gender: gender, FavoriteMusicAge: favoriteMusicAge, FavoriteArtist: favoriteArtist, Comment: comment}).Error; err != nil {
		var error model.Error
		error.Message = "ユーザー情報の更新に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
})

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

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/signup", SignUpHandler).Methods("POST")
	r.HandleFunc("/api/login", LoginHandler).Methods("POST")
	r.Handle("/api/user", JwtMiddleware.Handler(UserHandler)).Methods("GET")
	r.Handle("/api/users", JwtMiddleware.Handler(AllUsersHandler)).Methods("GET")
	r.Handle("/api/user/{id}/update", JwtMiddleware.Handler(UpdateUserHandler)).Methods("PUT")

	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Println(err)
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

// Parse は jwt トークンから元になった認証情報を取り出す。
func Parse(signedString string) (*model.Auth, error) {
	//追加
	secret := os.Getenv("SIGNINGKEY")

	token, err := jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.Errorf("unexpected signing method: %v", token.Header)
		}
		return []byte(secret), nil
	})

	if token == nil {
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
