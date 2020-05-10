package main

import (
	"encoding/json"
	"fmt"
	"golang-songs/controller"
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

// レスポンスにエラーを突っ込んで、返却するメソッド
func errorInResponse(w http.ResponseWriter, status int, error model.Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
	return
}

type SignUpHandler struct {
	DB *gorm.DB
}

func (f *SignUpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	password = string(hash)

	if err := f.DB.Create(&model.User{Email: email, Password: password}).Error; err != nil {
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

type LoginHandler struct {
	DB *gorm.DB
}

func (f *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var user model.User

	var jwt model.JWT

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

	var userData model.User
	row := f.DB.Where("email = ?", user.Email).Find(&userData)
	if err := f.DB.Where("email = ?", user.Email).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
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

	if _, err := w.Write(v2); err != nil {
		var error model.Error
		error.Message = "JWTトークンの取得に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}

type UserHandler struct {
	DB *gorm.DB
}

//リクエストユーザーの情報を返す
func (f *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headerAuthorization := r.Header.Get("Authorization")
	if len(headerAuthorization) == 0 {
		var error model.Error
		error.Message = "認証トークンの取得に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	bearerToken := strings.Split(headerAuthorization, " ")
	if len(bearerToken) < 2 {
		var error model.Error
		error.Message = "bearerトークンの取得に失敗しました。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	authToken := bearerToken[1]

	parsedToken, err := Parse(authToken)
	if err != nil {
		var error model.Error
		error.Message = "認証コードのパースに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	userEmail := parsedToken.Email

	var user model.User

	if err := f.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
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

	if _, err := w.Write(v); err != nil {
		var error model.Error
		error.Message = "ユーザー情報の取得に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}

type AllUsersHandler struct {
	DB *gorm.DB
}

//全てのユーザーを返す
func (f *AllUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	allUsers := []model.User{}

	if err := f.DB.Find(&allUsers).Error; gorm.IsRecordNotFoundError(err) {
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
	if _, err := w.Write(v); err != nil {
		var error model.Error
		error.Message = "ユーザー一覧の取得に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}

type UpdateUserHandler struct {
	DB *gorm.DB
}

func (f *UpdateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		var error model.Error
		error.Message = "ユーザーのidを取得できません。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	dec := json.NewDecoder(r.Body)
	var d model.User
	if err := dec.Decode(&d); err != nil {
		var error model.Error
		error.Message = "リクエストボディのデコードに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	var user model.User

	if err := f.DB.Model(&user).Where("id = ?", id).Update(model.User{Email: d.Email, Name: d.Name, Age: d.Age, Gender: d.Gender, FavoriteMusicAge: d.FavoriteMusicAge, FavoriteArtist: d.FavoriteArtist, Comment: d.Comment}).Error; err != nil {
		var error model.Error
		error.Message = "ユーザー情報の更新に失敗しました。"
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

type CreateSongHandler struct {
	DB *gorm.DB
}

func (f *CreateSongHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var d model.Song
	if err := dec.Decode(&d); err != nil {
		var error model.Error
		error.Message = "リクエストボディのデコードに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	header_hoge := r.Header.Get("Authorization")
	bearerToken := strings.Split(header_hoge, " ")
	authToken := bearerToken[1]

	parsedToken, err := Parse(authToken)
	if err != nil {
		var error model.Error
		error.Message = "認証コードのパースに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	userEmail := parsedToken.Email

	var user model.User

	if err := f.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	if err := f.DB.Create(&model.Song{
		Title:          d.Title,
		Artist:         d.Artist,
		MusicAge:       d.MusicAge,
		Image:          d.Image,
		Video:          d.Video,
		Album:          d.Album,
		Description:    d.Description,
		SpotifyTrackId: d.SpotifyTrackId,
		UserID:         user.ID}).Error; err != nil {
		var error model.Error
		error.Message = "曲の追加に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}

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

	db.DB().SetMaxIdleConns(10)
	defer db.Close()

	r := mux.NewRouter()

	r.Handle("/api/signup", &SignUpHandler{DB: db}).Methods("POST")
	r.Handle("/api/login", &LoginHandler{DB: db}).Methods("POST")
	r.Handle("/api/user", JwtMiddleware.Handler(&UserHandler{DB: db})).Methods("GET")
	r.Handle("/api/users", JwtMiddleware.Handler(&AllUsersHandler{DB: db})).Methods("GET")
	r.Handle("/api/user/{id}/update", JwtMiddleware.Handler(&UpdateUserHandler{DB: db})).Methods("PUT")

	r.Handle("/api/song", JwtMiddleware.Handler(&CreateSongHandler{DB: db})).Methods("POST")

	r.HandleFunc("/api/getRedirectUrl", controller.GetRedirectURL).Methods("GET")
	r.HandleFunc("/api/getToken", controller.GetToken).Methods("POST")
	r.HandleFunc("/api/tracks", controller.GetTracks).Methods("POST")

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
