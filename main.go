package main

import (
	"golang-songs/controller"
	"golang-songs/model"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"encoding/json"
	"fmt"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/davecgh/go-spew/spew"
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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

// Auth は署名前の認証トークン情報を表す。
type Auth struct {
	Email string
	Iss   int64
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

	db, _ := gormConnect()
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

	db, _ := gormConnect()
	defer db.Close()

	var userData model.User
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
	//log.Println(userData)

	w.WriteHeader(http.StatusOK)
	jwt.Token = token

	v2, err := json.Marshal(jwt)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(v2)
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	header_hoge := r.Header.Get("Authorization")
	bearerToken := strings.Split(header_hoge, " ")
	authToken := bearerToken[1]

	parsedToken, err := Parse(authToken)
	userEmail := parsedToken.Email

	db, _ := gormConnect()
	defer db.Close()

	var user model.User
	if err := db.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当するユーザーが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	followings := []model.User{}
	db.Preload("Followings").Find(&user)
	db.Model(&user).Related(&followings, "Followings")

	bookmarkings := []model.Song{}
	db.Preload("Bookmarkings").Find(&user)
	db.Model(&user).Related(&bookmarkings, "Bookmarkings")

	log.Println("user:", user)

	//if _, err := json.Marshal(row); err != nil {
	//	error := model.Error{}
	//	error.Message = "該当するアカウントが見つかりません。"
	//	errorInResponse(w, http.StatusUnauthorized, error)
	//	return
	//}
	//

	v, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}
	log.Println(v)
	w.Write(v)
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

// Parse は jwt トークンから元になった認証情報を取り出す。
func Parse(signedString string) (*Auth, error) {
	//追加
	secret := os.Getenv("SIGNINGKEY")
	//var err model.Error
	var err error

	token, err := jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		//if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		//return "", err.Errorf("unexpected signing method: %v", token.Header["alg"])
		//}
		return []byte(secret), nil
	})

	//if err != nil {
	//	if ve, ok := err.(*jwt.ValidationError); ok {
	//		if ve.Errors&jwt.ValidationErrorExpired != 0 {
	//			return nil, err.Wrapf(err, "%s is expired", signedString)
	//		} else {
	//			return nil, err.Wrapf(err, "%s is invalid", signedString)
	//		}
	//	} else {
	//		return nil, err.Wrapf(err, "%s is invalid", signedString)
	//	}
	//}

	if token == nil {
		//err.Message = "not found token in signedString"
		return nil, err
	}

	//claims, ok := token.Claims.(jwt.MapClaims)
	claims, _ := token.Claims.(jwt.MapClaims)
	//if !ok {
	//	return nil, err.Errorf("not found claims in %s", signedString)
	//}

	//userID, ok := claims[userIDKey].(string)
	//if !ok {
	//	return nil, err.Errorf("not found %s in %s", userIDKey, signedString)
	//}
	//email, ok := claims[email].(string)
	email, _ := claims["email"].(string)
	//if !ok {
	//	error := model.Error{}
	//	error.Message = "JSONへの変換失敗"
	//errorInResponse(w, http.StatusUnauthorized, error)
	//return nil, errorInResponse(w, http.StatusUnauthorized, error)
	//return nil, error.Message
	//return nil, err.Errorf("not found %s in %s", email, signedString)
	//}
	//iss, ok := claims[iss].(float64)
	iss, _ := claims["iss"].(float64)
	//if !ok {
	//	return nil, err.Errorf("not found %s in %s", iss, signedString)
	//}

	//return &Auth{
	//	Email: email,
	//	Iss:   int64(iss),
	//}, nil

	return &Auth{
		Email: email,
		Iss:   int64(iss),
	}, nil
}

//ユーザー情報取得
//var GetUserHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//id := vars["id"]
	//log.Println("id:", id)

	vars := mux.Vars(r)
	log.Println("vars:", vars)
	id, ok := vars["id"]
	log.Println("id:", id)

	if !ok {
		var error model.Error
		error.Message = "ユーザーのidを取得できません。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	header_hoge := r.Header.Get("Authorization")
	//log.Println("header_hoge:", header_hoge)
	bearerToken := strings.Split(header_hoge, " ")
	authToken := bearerToken[1]

	parsedToken, err := Parse(authToken)

	userEmail := parsedToken.Email
	log.Println("userEmail:", userEmail)

	db, _ := gormConnect()
	defer db.Close()

	//user := model.User{}
	//if err := db.Where("id = ?", id).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
	//	error := model.Error{}
	//	error.Message = "該当するユーザーが見つかりません。"
	//	errorInResponse(w, http.StatusUnauthorized, error)
	//	return
	//}
	//v, err := json.Marshal(user)
	//if err != nil {
	//	println(string(v))
	//}
	//w.Write(v)

	var user model.User

	if err := db.Where("id = ?", id).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
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

//type GetUserHandler struct {
//	DB *gorm.DB
//}
//
////指定されたユーザーの情報を返す
//func (f *GetUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	id, ok := vars["id"]
//	if !ok {
//		var error model.Error
//		error.Message = "ユーザーのidを取得できません。"
//		errorInResponse(w, http.StatusBadRequest, error)
//		return
//	}
//
//	var user model.User
//
//	if err := f.DB.Where("id = ?", id).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
//		error := model.Error{}
//		error.Message = "該当するアカウントが見つかりません。"
//		errorInResponse(w, http.StatusUnauthorized, error)
//		return
//	}
//	v, err := json.Marshal(user)
//	if err != nil {
//		var error model.Error
//		error.Message = "JSONへの変換に失敗しました"
//		errorInResponse(w, http.StatusInternalServerError, error)
//		return
//	}
//
//	if _, err := w.Write(v); err != nil {
//		var error model.Error
//		error.Message = "ユーザー情報の取得に失敗しました。"
//		errorInResponse(w, http.StatusInternalServerError, error)
//		return
//	}
//}

func AllUsersHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := gormConnect()
	defer db.Close()
	allUsers := []model.User{}

	if err := db.Find(&allUsers).Error; gorm.IsRecordNotFoundError(err) {
		var error model.Error
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	followings := []model.User{}
	db.Preload("Followings").Find(&allUsers)
	db.Model(&allUsers).Related(&followings, "Followings")

	log.Println("allUsers:", allUsers)

	v, _ := json.Marshal(allUsers)

	log.Println("v:", v)
	log.Println("allUsers:", allUsers)

	w.Write(v)
}

//func (f *AllUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//
//	allUsers := []model.User{}
//
//	if err := f.DB.Find(&allUsers).Error; gorm.IsRecordNotFoundError(err) {
//		var error model.Error
//		error.Message = "該当するアカウントが見つかりません。"
//		errorInResponse(w, http.StatusInternalServerError, error)
//		return
//	}
//
//	v, err := json.Marshal(allUsers)
//	if err != nil {
//		var error model.Error
//		error.Message = "ユーザー一覧の取得に失敗しました"
//		errorInResponse(w, http.StatusInternalServerError, error)
//		return
//	}
//	if _, err := w.Write(v); err != nil {
//		var error model.Error
//		error.Message = "ユーザー一覧の取得に失敗しました。"
//		errorInResponse(w, http.StatusInternalServerError, error)
//		return
//	}
//}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
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

	db.Model(&user).Where("id = ?", id).Update(model.User{Email: email, Name: name, Age: age, Gender: gender, FavoriteMusicAge: favoriteMusicAge, FavoriteArtist: favoriteArtist, Comment: comment})
}

var CreateSongHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var d model.Song
	dec.Decode(&d)

	header_hoge := r.Header.Get("Authorization")
	bearerToken := strings.Split(header_hoge, " ")
	authToken := bearerToken[1]

	parsedToken, _ := Parse(authToken)
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

	title := d.Title
	artist := d.Artist
	musicAge := d.MusicAge
	image := d.Image
	video := d.Video
	album := d.Album
	description := d.Description
	spotifyTrackId := d.SpotifyTrackId

	if err := db.Create(&model.Song{
		Title:          title,
		Artist:         artist,
		MusicAge:       musicAge,
		Image:          image,
		Video:          video,
		Album:          album,
		Description:    description,
		SpotifyTrackId: spotifyTrackId,
		UserID:         user.ID}).Error; err != nil {
		var error model.Error
		error.Message = "曲の追加に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
})

var GetSongHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	db, _ := gormConnect()
	defer db.Close()

	song := model.Song{}

	row := db.Where("id = ?", id).Find(&song)
	if err := db.Where("id = ?", id).Find(&song).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当する曲が見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
	if _, err := json.Marshal(row); err != nil {
		error := model.Error{}
		error.Message = "JSONへの変換失敗"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	v, err := json.Marshal(song)
	if err != nil {
		var error model.Error
		error.Message = "曲の取得に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	if _, err := w.Write(v); err != nil {
		var error model.Error
		error.Message = "曲の取得に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
})

var AllSongsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	db, _ := gormConnect()
	defer db.Close()

	allSongs := []model.Song{}

	if err := db.Find(&allSongs).Error; gorm.IsRecordNotFoundError(err) {
		var error model.Error
		error.Message = "曲が見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	v, err := json.Marshal(allSongs)
	if err != nil {
		var error model.Error
		error.Message = "曲一覧の取得に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	if _, err := w.Write(v); err != nil {
		var error model.Error
		error.Message = "曲一覧の取得に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
})

var UpdateSongHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	dec := json.NewDecoder(r.Body)
	var d model.Song
	dec.Decode(&d)

	title := d.Title
	artist := d.Artist
	musicAge := d.MusicAge
	image := d.Image
	video := d.Video
	album := d.Album
	description := d.Description
	spotifyTrackId := d.SpotifyTrackId

	db, _ := gormConnect()
	defer db.Close()

	var song model.Song

	if err := db.Model(&song).Where("id = ?", id).Update(model.Song{
		Title:          title,
		Artist:         artist,
		MusicAge:       musicAge,
		Image:          image,
		Video:          video,
		Album:          album,
		Description:    description,
		SpotifyTrackId: spotifyTrackId}).Error; err != nil {
		var error model.Error
		error.Message = "曲の更新に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
})

var DeleteSongHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	dec := json.NewDecoder(r.Body)
	var d model.Song
	dec.Decode(&d)

	db, _ := gormConnect()
	defer db.Close()

	var song model.Song

	if err := db.Where("id = ?", id).Delete(&song).Error; err != nil {
		var error model.Error
		error.Message = "曲の削除に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
})

type FollowUserHandler struct {
	DB *gorm.DB
}

func (f *FollowUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	dec := json.NewDecoder(r.Body)
	var d model.Song
	dec.Decode(&d)

	var targetUser model.User

	if err := f.DB.Where("id = ?", id).Find(&targetUser).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当するユーザーが見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	header_hoge := r.Header.Get("Authorization")
	bearerToken := strings.Split(header_hoge, " ")
	authToken := bearerToken[1]

	parsedToken, _ := Parse(authToken)
	userEmail := parsedToken.Email

	var requestUser model.User

	if err := f.DB.Where("email = ?", userEmail).Find(&requestUser).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	if err := f.DB.Create(&model.UserFollow{
		UserID:   requestUser.ID,
		FollowID: targetUser.ID}).Error; err != nil {
		var error model.Error
		error.Message = "曲の追加に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	f.DB.Preload("Followings").Find(&requestUser)
	f.DB.Model(&requestUser).Association("Followings").Append(&targetUser)

	//nの値は増えてる
	n := f.DB.Model(&requestUser).Association("Followings").Count() //動作
	log.Println("n:", n)

	log.Println("requestUser:", requestUser)

	log.Println("==============")
	followings := []model.User{}
	f.DB.Model(&requestUser).Related(&followings, "Followings")

	log.Println("requestUser:", requestUser) //Followingsに値入っているぽい
	log.Println("followings:", followings)
}

type UnfollowUserHandler struct {
	DB *gorm.DB
}

func (f *UnfollowUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	dec := json.NewDecoder(r.Body)
	var d model.Song
	dec.Decode(&d)

	var targetUser model.User

	if err := f.DB.Where("id = ?", id).Find(&targetUser).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当する曲が見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	header_hoge := r.Header.Get("Authorization")
	bearerToken := strings.Split(header_hoge, " ")
	authToken := bearerToken[1]

	parsedToken, _ := Parse(authToken)
	userEmail := parsedToken.Email

	var requestUser model.User

	if err := f.DB.Where("email = ?", userEmail).Find(&requestUser).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	if err := f.DB.Create(&model.UserFollow{
		UserID:   requestUser.ID,
		FollowID: targetUser.ID}).Error; err != nil {
		var error model.Error
		error.Message = "曲の追加に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	f.DB.Preload("Followings").Find(&requestUser)
	f.DB.Model(&requestUser).Association("Followings").Delete(&targetUser)

	n := f.DB.Model(&requestUser).Association("Followings").Count() //動作
	log.Println("n:", n)

	log.Println("==============")
	followings := []model.User{}
	f.DB.Model(&requestUser).Related(&followings, "Followings")

	log.Println("requestUser:", requestUser) //Followingsに値入っているぽい
	log.Println("followings:", followings)
}

type BookmarkHandler struct {
	DB *gorm.DB
}

func (f *BookmarkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	dec := json.NewDecoder(r.Body)
	var d model.Song
	dec.Decode(&d)

	var song model.Song

	if err := f.DB.Where("id = ?", id).Find(&song).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当する曲が見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	header_hoge := r.Header.Get("Authorization")
	bearerToken := strings.Split(header_hoge, " ")
	authToken := bearerToken[1]

	parsedToken, _ := Parse(authToken)
	userEmail := parsedToken.Email

	var user model.User

	if err := f.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	if err := f.DB.Create(&model.Bookmark{
		UserID: user.ID,
		SongID: song.ID}).Error; err != nil {
		var error model.Error
		error.Message = "曲のお気に入り登録に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	f.DB.Preload("Bookmarkings").Find(&user)
	f.DB.Model(&user).Association("Bookmarkings").Append(&song)

	//nの値は増えてる
	n := f.DB.Model(&user).Association("Bookmarkings").Count() //動作
	log.Println("n:", n)

	log.Println("user:", user)

	log.Println("==============")
	bookmarkings := []model.Song{}
	f.DB.Model(&user).Related(&bookmarkings, "Bookmarikings")

	log.Println("user:", user)
	log.Println("bookmarkings:", bookmarkings)
}

type RemoveBookmarkHandler struct {
	DB *gorm.DB
}

func (f *RemoveBookmarkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	dec := json.NewDecoder(r.Body)
	var d model.Song
	dec.Decode(&d)

	var song model.Song

	if err := f.DB.Where("id = ?", id).Find(&song).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当する曲が見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	header_hoge := r.Header.Get("Authorization")
	bearerToken := strings.Split(header_hoge, " ")
	authToken := bearerToken[1]

	parsedToken, _ := Parse(authToken)
	userEmail := parsedToken.Email

	var user model.User

	if err := f.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	f.DB.Preload("Bookmarkings").Find(&user)
	f.DB.Model(&user).Association("Bookmarkings").Delete(&song)

	//nの値は増えてる
	n := f.DB.Model(&user).Association("Bookmarkings").Count() //動作
	log.Println("n:", n)

	log.Println("user:", user)

	log.Println("==============")
	bookmarkings := []model.Song{}
	f.DB.Model(&user).Related(&bookmarkings, "Bookmarikings")

	log.Println("user:", user)
	log.Println("bookmarkings:", bookmarkings)
}

type UploadSongImageHandler struct {
	DB *gorm.DB
}

func (f *UploadSongImageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	dec := json.NewDecoder(r.Body)
	var d model.Song
	dec.Decode(&d)

	var song model.Song

	if err := f.DB.Where("id = ?", id).Find(&song).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当する曲が見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	header_hoge := r.Header.Get("Authorization")
	bearerToken := strings.Split(header_hoge, " ")
	authToken := bearerToken[1]

	parsedToken, _ := Parse(authToken)
	userEmail := parsedToken.Email

	var user model.User

	if err := f.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	//sess, err := session.NewSession(&aws.Config{
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1")},
	)
	// ファイルを開く
	targetFilePath := "./music_life.jpg"
	file, err := os.Open(targetFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//bucketName := "xxx-bucket"
	bucketName := "your-songs-laravel"
	objectKey := "AKIAZNEYIQFBYJLFB4EK"

	// Uploaderを作成し、ローカルファイルをアップロード
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("done")
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

	//S3に接続するためにのsessionを作成
	//sess := session.Must(session.NewSessionWithOptions(session.Options{
	//	Profile: "di",
	//	//SharedConfigState: session.SharedConfigEnable,
	//}))

	////sess, err := session.NewSession(&aws.Config{
	//sess, _ := session.NewSession(&aws.Config{
	//	Region: aws.String("ap-northeast-1")},
	//)
	//// ファイルを開く
	//targetFilePath := "./music_life.jpg"
	//file, err := os.Open(targetFilePath)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer file.Close()
	//
	////bucketName := "xxx-bucket"
	//bucketName := "your-songs-laravel"
	//objectKey := "AKIAZNEYIQFBYJLFB4EK"
	//
	//// Uploaderを作成し、ローカルファイルをアップロード
	//uploader := s3manager.NewUploader(sess)
	//_, err = uploader.Upload(&s3manager.UploadInput{
	//	Bucket: aws.String(bucketName),
	//	Key:    aws.String(objectKey),
	//	Body:   file,
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println("done")

	r := mux.NewRouter()

	r.HandleFunc("/api/signup", SignUpHandler).Methods("POST")
	r.HandleFunc("/api/login", LoginHandler).Methods("POST")
	r.HandleFunc("/api/user", UserHandler).Methods("GET")
	r.HandleFunc("/api/user/{id}", GetUserHandler).Methods("GET")
	//r.Handle("/api/user/{id}", JwtMiddleware.Handler(&GetUserHandler{DB: db})).Methods("GET")
	r.HandleFunc("/api/users", AllUsersHandler).Methods("GET")
	r.HandleFunc("/api/user/{id}/update", UpdateUserHandler).Methods("PUT")

	r.HandleFunc("/api/oauth", controller.OAuth).Methods("POST")
	r.HandleFunc("/api/get-token", controller.GetToken).Methods("POST")
	r.HandleFunc("/api/tracks", controller.GetTracks).Methods("POST")
	r.HandleFunc("/api/get-redirect-url", controller.GetRedirectURL).Methods("GET")

	r.Handle("/api/song", JwtMiddleware.Handler(CreateSongHandler)).Methods("POST")
	r.Handle("/api/song/{id}", JwtMiddleware.Handler(GetSongHandler)).Methods("GET")
	r.Handle("/api/songs", JwtMiddleware.Handler(AllSongsHandler)).Methods("GET")
	r.Handle("/api/song/{id}", JwtMiddleware.Handler(UpdateSongHandler)).Methods("PUT")
	r.Handle("/api/song/{id}", JwtMiddleware.Handler(DeleteSongHandler)).Methods("DELETE")

	r.Handle("/api/user/{id}/follow", JwtMiddleware.Handler(&FollowUserHandler{DB: db})).Methods("POST")
	r.Handle("/api/user/{id}/unfollow", JwtMiddleware.Handler(&UnfollowUserHandler{DB: db})).Methods("POST")

	r.Handle("/api/song/{id}/bookmark", JwtMiddleware.Handler(&BookmarkHandler{DB: db})).Methods("POST")
	r.Handle("/api/song/{id}/remove-bookmark", JwtMiddleware.Handler(&RemoveBookmarkHandler{DB: db})).Methods("POST")

	r.Handle("/api/song/{id}/upload", JwtMiddleware.Handler(&UploadSongImageHandler{DB: db})).Methods("POST")

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
