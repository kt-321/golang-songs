package main

import (
	"golang-songs/controller"
	"golang-songs/model"
	"log"
	"os"
	"strings"

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

func gormConnect() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println(".envファイルの読み込み失敗")
	}
	mysqlConfig := os.Getenv("mysqlConfig")
	log.Println("mysqlConfig:", mysqlConfig)
	db, err := gorm.Open("mysql", mysqlConfig)

	//db, err := gorm.Open("mysql", "root:@/golang_songs?charset=utf8&parseTime=True&loc=Local")
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
	//log.Println("header_hoge:", header_hoge)
	bearerToken := strings.Split(header_hoge, " ")
	authToken := bearerToken[1]
	//fmt.Println("bearerToken: ", bearerToken)
	fmt.Println("authToken: ", authToken)

	parsedToken, err := Parse(authToken)
	userEmail := parsedToken.Email
	log.Println("userEmail:", userEmail)

	db := gormConnect()
	defer db.Close()

	//userData := model.User{}

	vars := mux.Vars(r)
	id := vars["id"]
	log.Println(id)

	var user model.User
	//var user model.UserInResponse

	//row := db.Where("email = ?", user.Email).Find(&userData)
	row := db.Where("email = ?", userEmail).Find(&user)
	if err := db.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当するユーザーが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	log.Println(row)

	//if _, err := json.Marshal(row); err != nil {
	//	error := model.Error{}
	//	error.Message = "該当するアカウントが見つかりません。"
	//	errorInResponse(w, http.StatusUnauthorized, error)
	//	return
	//}
	//

	//user.Password = ""

	log.Println("user:", user.Age)
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

//JWT認証のテスト 成功
var TestHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	post := "test"
	json.NewEncoder(w).Encode(post)
})

//ユーザー情報取得
var GetUserHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Println("id:", id)

	header_hoge := r.Header.Get("Authorization")
	//log.Println("header_hoge:", header_hoge)
	bearerToken := strings.Split(header_hoge, " ")
	authToken := bearerToken[1]
	//fmt.Println("bearerToken: ", bearerToken)
	fmt.Println("authToken: ", authToken)

	parsedToken, err := Parse(authToken)
	userEmail := parsedToken.Email
	log.Println("userEmail:", userEmail)
	//parsedToken2 := *parsedToken
	//userEmail := parsedToken2
	//log.Println("parsedToken", parsedToken)
	//log.Println("parsedToken2", parsedToken2)

	db := gormConnect()
	defer db.Close()

	user := model.User{}
	//userData := model.User{}

	//row := db.Where("id = ?", id).Find(&userData)
	row := db.Where("id = ?", id).Find(&user)
	if err := db.Where("id = ?", id).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当するユーザーが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}
	if _, err := json.Marshal(row); err != nil {
		error := model.Error{}
		error.Message = "JSONへの変換失敗"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}
	log.Println("user:", user.Age)
	v, err := json.Marshal(user)
	if err != nil {
		println(string(v))
	}
	w.Write(v)
})

func AllUsersHandler(w http.ResponseWriter, r *http.Request) {
	db := gormConnect()
	defer db.Close()
	allUsers := []model.User{}

	db.Find(&allUsers)
	v, _ := json.Marshal(allUsers)
	w.Write(v)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	//クエリパラメータの取得
	vars := mux.Vars(r)
	id := vars["id"]

	//フォームの値を取得
	name := r.FormValue("Name")
	log.Println(name)

	//UPDATE成功
	db := gormConnect()
	defer db.Close()
	var user model.User

	db.Model(&user).Where("id = ?", id).Update("name", name)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/signup", SignUpHandler).Methods("POST")
	r.HandleFunc("/api/login", LoginHandler).Methods("POST")
	r.HandleFunc("/api/user", UserHandler).Methods("GET")
	r.HandleFunc("/api/user/{id}", GetUserHandler).Methods("GET")
	r.HandleFunc("/api/users", AllUsersHandler).Methods("GET")

	//JWT認証のテスト
	r.Handle("/api/test", JwtMiddleware.Handler(TestHandler)).Methods("GET")

	//r.GET("/api/tracks", controller.GetTracks)
	//r.Handle("/api/oauth", JwtMiddleware.Handler(controller.OAuth)).Methods("GET")
	//r.Handle("/api/oauth", JwtMiddleware.Handler(controller.OAuth)).Methods("POST")
	r.HandleFunc("/api/oauth", controller.OAuth).Methods("POST")
	r.HandleFunc("/api/getToken", controller.GetToken).Methods("POST")
	r.HandleFunc("/api/tracks", controller.GetTracks).Methods("POST")
	r.HandleFunc("/api/getRedirectUrl", controller.GetRedirectURL).Methods("GET")
	//r.Handle("/api/test", JwtMiddleware.Handler(controller.GetTracks)).Methods("GET")

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
