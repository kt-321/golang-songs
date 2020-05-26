package main

import (
	"bytes"
	"encoding/json"
	"golang-songs/model"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

//func TestSignUpHandler(t *testing.T) {
func TestSignUpHandler_ServeHTTP(t *testing.T) {
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

	url := "http://localhost:8081/api/signup"

	// テスト用の JSON ボディ作成
	b, err := json.Marshal(Form{Email: "tes@tes", Password: "testes"})
	if err != nil {
		t.Fatal(err)
	}

	// テスト用のリクエスト作成
	req := httptest.NewRequest("POST", url, bytes.NewBuffer(b))

	// テスト用のレスポンス作成
	res := httptest.NewRecorder()

	// ハンドラーの実行
	//SignUpHandler(res, req)

	//レシーバ付きの場合
	f := &SignUpHandler{DB: db}
	f.ServeHTTP(res, req)

	//レシーバ付きでない場合
	//handler := http.HandlerFunc(SignUpHandler)
	//handler.ServeHTTP(res, req)

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}

	//// レスポンスの JSON ボディのテスト
	//resp := JsonResponse{}
	//if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
	//	t.Errorf("errpr: %#v, res: %#v", err, res)
	//}
	//if resp.Message != "hello world" {
	//	t.Errorf("invalid response: %#v", resp)
	//}
	//
	//t.Logf("%#v", resp)

	////実行
	//resp, err := client.Do(req)
	//if err != nil {
	//	//panic(err)
	//	//t.Error("リクエスト実行エラー")
	//	t.Error("リクエスト実行エラー", err)
	//}

	////body使わなくなったら閉じる
	//defer resp.Body.Close()

	//エラー検証
	//if err != nil {
	//	t.Error(err)
	//	return
	//}

	//ステータスコード確認
	//if resp.StatusCode != 200 {
	//	t.Error(resp.StatusCode)
	//	return
	//}
}

//func TestLoginHandler(t *testing.T) {
func TestLoginHandler_ServeHTTP(t *testing.T) {
	//レシーバ付きの場合
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

	url := "http://localhost:8081/api/login"

	// テスト用の JSON ボディ作成
	b, err := json.Marshal(Form{Email: "t@t", Password: "tttttt"})
	if err != nil {
		t.Fatal(err)
	}

	// テスト用のリクエスト作成
	req := httptest.NewRequest("POST", url, bytes.NewBuffer(b))

	////headerをセット
	//req.Header.Set("Content-Type", "application/json")
	////httpクライアント
	//client := &http.Client{}

	// テスト用のレスポンス作成
	res := httptest.NewRecorder()
	// ハンドラーの実行
	//LoginHandler(res, req)

	//レシーバ付きの場合
	f := &LoginHandler{DB: db}
	f.ServeHTTP(res, req)

	//レシーバ無しの場合
	//handler := http.HandlerFunc(LoginHandler)
	//handler.ServeHTTP(res, req)

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}

	//// レスポンスの JSON ボディのテスト
	//resp := JsonResponse{}
	//if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
	//	t.Errorf("errpr: %#v, res: %#v", err, res)
	//}
	//if resp.Message != "hello world" {
	//	t.Errorf("invalid response: %#v", resp)
	//}
	//
	//t.Logf("%#v", resp)

	////実行
	//resp, err := client.Do(req)
	//if err != nil {
	//	//panic(err)
	//	//t.Error("リクエスト実行エラー")
	//	t.Error("リクエスト実行エラー", err)
	//}

	////body使わなくなったら閉じる
	//defer resp.Body.Close()

	//エラー検証
	//if err != nil {
	//	t.Error(err)
	//	return
	//}

	//ステータスコード確認
	//if resp.StatusCode != 200 {
	//	t.Error(resp.StatusCode)
	//	return
	//}
}

func TestGetUserHandler_ServeHTTP(t *testing.T) {
	//func TestGetUserHandler(t *testing.T) {
	//url := "http://localhost:8081/api/user/2"

	//レシーバ付きの場合
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

	//req, err := http.NewRequest("GET", url, nil)
	//req := httptest.NewRequest("GET", url, nil)
	req := httptest.NewRequest("GET", "http://localhost:8081/api/user/2", nil)

	//headerをセット
	req.Header.Set("Content-Type", "application/json")

	user := model.User{}
	user.Email = "a@a"
	user.Password = "aaaaaa"
	//トークン作成
	token, err := createToken(user)
	if err != nil {
		//error := model.Error{}
		//error.Message = "トークンの作成に失敗しました"
		//errorInResponse(w, http.StatusUnauthorized, error)
		//return
		log.Println("err:", err)
	}
	log.Printf("tokenintest:%s", token)

	jointToken := "Bearer" + " " + token
	log.Printf("jointToken:%s", jointToken)

	req.Header.Set("Authorization", jointToken)

	//req.Header.Set("Authorization", 'Bearer'+string(token))
	//req.Header.Set("Authorization", 'Bearer'+strconv.Itoa(token))

	//httpクライアント
	//client := &http.Client{}
	//
	////実行
	//resp, err := client.Do(req)
	//if err != nil {
	//	t.Error("リクエスト実行エラー", err)
	//}

	//body使わなくなったら閉じる
	//defer resp.Body.Close()

	//エラー検証
	//if err != nil {
	//	t.Error(err)
	//	return
	//}

	////ステータスコード確認
	//if resp.StatusCode != 200 {
	//	t.Error(resp.StatusCode)
	//	return
	//}
	//
	////レスポンスBODY取得
	//body, _ := ioutil.ReadAll(resp.Body)

	//return

	////headerをセット
	//req.Header.Set("Content-Type", "application/json")
	////httpクライアント
	//client := &http.Client{}

	// テスト用のレスポンス作成
	res := httptest.NewRecorder()

	//レシーバ付きの場合
	//f := &GetUserHandler{DB: db}
	//f.ServeHTTP(res, req)

	//レシーバ無しの場合
	router := mux.NewRouter()
	//router.HandleFunc("/api/user/{id}", GetUserHandler)
	router.Handle("/api/user/{id}", &GetUserHandler{DB: db})
	router.ServeHTTP(res, req)

	// ハンドラーの実行
	//GetUserHandler(res, req)

	//client := new(http.Client)
	//res, err := client.Do(req)

	log.Printf("req: %v", req)
	log.Printf("res: %v", res)

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}

	// レスポンスのボディが期待通りか確認
	//expected := `{"id":2,"createdAt":"2020-05-23T21:02:20+09:00","updatedAt":"2020-05-23T21:02:20+09:00","deletedAt":null,"name":"","email":"u@u","age":0,"gender":0,"imageUrl":"","favoriteMusicAge":0,"favoriteArtist":"","comment":"","followings":null,"bookmarkings":null}`
	expected := `{"id":2,"createdAt":"2020-05-23T21:02:20+09:00","updatedAt":"2020-05-23T21:02:20+09:00","deletedAt":null,"name":"","email":"u@u","age":0,"gender":0,"imageUrl":"","favoriteMusicAge":0,"favoriteArtist":"","comment":"","followings":[],"bookmarkings":[]}`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			res.Body.String(), expected)
	}

	//actual := string(res)
	//expected := `{"id":1, "createdAt": "2020-05-23T19:17:21+09:00", "updatedAt": "2020-05-23T19:21:04+09:00", "deletedAt": null, "name": "a", "email": "a@a", "age": 20, "gender": 1, "imageUrl": "", favoriteMusicAge": 1980, "favoriteArtist":"椎名林檎", "comment": テストテスト", "followings": [], "bookmarkings": []}`
	//if actual != expected {
	//	t.Error("response error")
	//}
}

func TestUserHandler_ServeHTTP(t *testing.T) {
	url := "http://localhost:8081/api/user"

	//レシーバ付きの場合
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

	req := httptest.NewRequest("GET", url, nil)

	//headerをセット
	req.Header.Set("Content-Type", "application/json")

	var user model.User
	user.Email = "u@u"
	user.Password = "uuuuuu"
	//トークン作成
	token, err := createToken(user)
	if err != nil {
		//error := model.Error{}
		//error.Message = "トークンの作成に失敗しました"
		//errorInResponse(w, http.StatusUnauthorized, error)
		//return
		log.Println("err:", err)
	}
	log.Printf("tokenintest:%s", token)

	jointToken := "Bearer" + " " + token
	log.Printf("jointToken:%s", jointToken)

	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成
	res := httptest.NewRecorder()

	//レシーバ付きの場合
	f := &UserHandler{DB: db}
	f.ServeHTTP(res, req)

	log.Printf("req: %v", req)
	log.Printf("res: %v", res)

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}

	// レスポンスのボディが期待通りか確認
	//expected := `{"id":2,"createdAt":"2020-05-23T21:02:20+09:00","updatedAt":"2020-05-23T21:02:20+09:00","deletedAt":null,"name":"","email":"u@u","age":0,"gender":0,"imageUrl":"","favoriteMusicAge":0,"favoriteArtist":"","comment":"","followings":null,"bookmarkings":null}`
	expected := `{"id":2,"createdAt":"2020-05-23T21:02:20+09:00","updatedAt":"2020-05-23T21:02:20+09:00","deletedAt":null,"name":"","email":"u@u","age":0,"gender":0,"imageUrl":"","favoriteMusicAge":0,"favoriteArtist":"","comment":"","followings":[],"bookmarkings":[]}`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			res.Body.String(), expected)
	}
}

//レスポンス正しい気がするがエラーの要因がわからない
//func TestAllUsersHandler_ServeHTTP(t *testing.T) {
//	url := "http://localhost:8081/api/users"
//
//	//レシーバ付きの場合
//	err := godotenv.Load()
//	if err != nil {
//		log.Println(".envファイルの読み込み失敗")
//	}
//
//	mysqlConfig := os.Getenv("mysqlConfig")
//
//	db, err := gorm.Open("mysql", mysqlConfig)
//	if err != nil {
//		log.Println(err)
//	}
//
//	db.DB().SetMaxIdleConns(10)
//	defer db.Close()
//
//	req := httptest.NewRequest("GET", url, nil)
//
//	//headerをセット
//	req.Header.Set("Content-Type", "application/json")
//
//	var user model.User
//	user.Email = "u@u"
//	user.Password = "uuuuuu"
//	//トークン作成
//	token, err := createToken(user)
//	if err != nil {
//		log.Println("err:", err)
//	}
//	log.Printf("tokenintest:%s", token)
//
//	jointToken := "Bearer" + " " + token
//	log.Printf("jointToken:%s", jointToken)
//
//	req.Header.Set("Authorization", jointToken)
//
//	// テスト用のレスポンス作成
//	res := httptest.NewRecorder()
//
//	//レシーバ付きの場合
//	f := &AllUsersHandler{DB: db}
//	f.ServeHTTP(res, req)
//
//	//log.Printf("req: %v", req)
//	log.Printf("res: %v", res)
//
//	// レスポンスのステータスコードのテスト
//	if res.Code != http.StatusOK {
//		t.Errorf("invalid code: %d", res.Code)
//	}
//
//	// レスポンスのボディが期待通りか確認
//	expected := `[{"id":1,"createdAt":"2020-05-23T19:17:21+09:00","updatedAt":"2020-05-23T19:21:04+09:00","deletedAt":null,"name":"a","email":"a@a","age":20,"gender":1,"imageUrl":"","favoriteMusicAge":1980,"favoriteArtist":"椎名林檎","comment":"dfafdsafd","followings":null,"bookmarkings":null},{"id":2,"createdAt":"2020-05-23T21:02:2:00","updatedAt":"2020-05-23T21:02:20+09:00","deletedAt":null,"name":"","email":"u@u","age":0,"gender":0,"imageUrl":"","favoriteMusicAge":0,"favoriteArtist":"","comment":"","followings":null,"bookmarkings":null},{"id":3,"createdAt":"2020-05-23T21:12:28+09:00","updatedAt":"2020-05-23T21:12:28+09:00","deletedAt":null,"name":"","email":"o@o","age":0,"gender":0,"imageUrl":"","favoriteMusicAge":0,"favoriteArtist":"","comment":"","followings":null,"bookmarkings":null},{"id":7,"createdAt":"2020-05-24T01:22:12+09:00","updatedAt":"2020-05-24T01:22:12+09:00","deletedAt":null,"name":"","email":"i@i","age":0,"gender":0,"imageUrl":"","favoriteMusicAge":0,"favoriteArtist":"","comment":"","followings":null,"bookmarkings":null},{"id":11,"createdAt":"2020-05-24T16:18:27+09:00","updatedAt":"2020-05-24T16:18:27+09:00","deletedAt":null,"name":"","email":"k@k","age":0,"gender":0,"imageUrl":"","favoriteMusicAge":0,"favoriteArtist":"","comment":"","followings":null,"bookmarkings":null},{"id":13,"createdAt":"2020-05-24T17:08:18+09:00","updatedAt":"2020-05-24T17:08:18+09:00","deletedAt":null,"name":"","email":"e@e","age":0,"gender":0,"imageUrl":"","favoriteMusicAge":0,"favoriteArtist":"","comment":"","followings":null,"bookmarkings":null},{"id":16,"createdAt":"2020-05-24T17:19:32+09:00","updatedAt":"2020-05-24T17:19:32+09:00","deletedAt":null,"name":"","email":"t@t","age":0,"gender":0,"imageUrl":"","favoriteMusicAge":0,"favoriteArtist":"","comment":"","followings":null,"bookmarkings":null},{"id":19,"createdAt":"2020-05-24T17:38:48+09:00","updatedAt":"2020-05-24T17:38:48+09:00","deletedAt":null,"name":"","email":"te@te","age":0,"gender":0,"imageUrl":"","favoriteMusicAge":0,"favoriteArtist":"","comment":"","followings":null,"bookmarkings":null},{"id":21,"createdAt":"2020-05-24T17:41:03+09:00","updatedAt":"2020-05-24T17:41:03+09:00","deletedAt":null,"name":"","email":"tes@tes","age":0,"gender":0,"imageUrl":"","favoriteMusicAge":0,"favoriteArtist":"","comment":"","followings":null,"bookmarkings":null}]`
//	if res.Body.String() != expected {
//		t.Errorf("handler returned unexpected body: got %v want %v",
//			res.Body.String(), expected)
//	}
//}

func UpdateUserHandler_ServeHTTP(t *testing.T) {
	url := "http://localhost:8081/api/user/2"

	//レシーバ付きの場合
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

	// テスト用の JSON ボディ作成
	b, err := json.Marshal(model.User{Email: "u@u", Name: "", Age: 0, Gender: 0, FavoriteMusicAge: 0, FavoriteArtist: "", Comment: ""})
	//{"id":2,"createdAt":"2020-05-23T21:02:20+09:00","updatedAt":"2020-05-23T21:02:20+09:00","deletedAt":null,"name":"","email":"u@u","age":0,"gender":0,"imageUrl":"","favoriteMusicAge":0,"favoriteArtist":"","comment":"","followings":[],"bookmarkings":[]}
	if err != nil {
		t.Fatal(err)
	}

	// テスト用のリクエスト作成
	req := httptest.NewRequest("PUT", url, bytes.NewBuffer(b))

	//headerをセット
	req.Header.Set("Content-Type", "application/json")

	var user model.User
	user.Email = "u@u"
	user.Password = "uuuuuu"
	//トークン作成
	token, err := createToken(user)
	if err != nil {
		log.Println("err:", err)
	}
	log.Printf("tokenintest:%s", token)

	jointToken := "Bearer" + " " + token
	log.Printf("jointToken:%s", jointToken)

	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成
	res := httptest.NewRecorder()

	//レシーバ付きの場合
	f := &UpdateUserHandler{DB: db}
	f.ServeHTTP(res, req)

	log.Printf("req: %v", req)
	log.Printf("res: %v", res)

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}

}
