package interfaces

import (
	"bytes"
	"encoding/json"
	"golang-songs/model"
	"golang-songs/usecases"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

type FakeUserRepository struct{}

func (fur *FakeUserRepository) FindAll() (*[]model.User, error) {
	user1 := model.User{
		ID:               1,
		CreatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		UpdatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		Password:         "aaaaaa",
		Name:             "",
		Email:            "a@test.co.jp",
		Age:              0,
		Gender:           0,
		ImageUrl:         "",
		FavoriteMusicAge: 0,
		FavoriteArtist:   "",
		Comment:          "",
	}

	user2 := model.User{
		ID:               2,
		CreatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		UpdatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		Password:         "iiiiii",
		Name:             "",
		Email:            "i@test.co.jp",
		Age:              0,
		Gender:           0,
		ImageUrl:         "",
		FavoriteMusicAge: 0,
		FavoriteArtist:   "",
		Comment:          "",
	}

	users := []model.User{user1, user2}

	return &users, nil
}

func (fur *FakeUserRepository) GetUser(userEmail string) (*model.User, error) {
	user := model.User{
		ID:               1,
		CreatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		UpdatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		Password:         "aaaaaa",
		Name:             "",
		Email:            "a@test.co.jp",
		Age:              0,
		Gender:           0,
		ImageUrl:         "",
		FavoriteMusicAge: 0,
		FavoriteArtist:   "",
		Comment:          "",
	}

	return &user, nil
}

func (fur *FakeUserRepository) FindByID(userID int) (*model.User, error) {
	user := model.User{
		ID:               1,
		CreatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		UpdatedAt:        time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local),
		Password:         "testtest",
		Name:             "",
		Email:            "test@test.co.jp",
		Age:              0,
		Gender:           0,
		ImageUrl:         "",
		FavoriteMusicAge: 0,
		FavoriteArtist:   "",
		Comment:          "",
	}

	return &user, nil
}

func (fur *FakeUserRepository) Update(userID int, p model.User) error {
	return nil
}

//idで指定したユーザーの情報を返すハンドラのテスト
func TestGetUserHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/user/1", nil)

	//headerをセット
	req.Header.Set("Content-Type", "application/json")

	//リクエストユーザー作成
	user := model.User{Email: "a@test.co.jp", Password: "aaaaaa"}

	//トークン作成
	token, err := createToken(user)
	if err != nil {
		log.Println("err:", err)
	}

	jointToken := "Bearer" + " " + token

	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成
	res := httptest.NewRecorder()

	fakeUserController := &UserController{
		UserInteractor: usecases.UserInteractor{
			UserRepository: &FakeUserRepository{},
		},
	}

	r := mux.NewRouter()
	r.Handle("/api/user/{id}", http.HandlerFunc(fakeUserController.GetUserHandler)).Methods("GET")
	r.ServeHTTP(res, req)

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}

	//log.Printf("res: %v", res)
	//log.Printf("res.Body: %v, %t", res.Body, res.Body)

	var p model.User
	//json.Marshal(res.Body, &p)
	//b, err := json.Marshal(res.Body)
	//if err != nil {
	//	log.Println(err)
	//}
	//if err := json.Unmarshal(test1, &p); err != nil{
	// fmt.Println(err)
	//}
	//fmt.Println(test1.Name)
	//if err := json.Unmarshal(b, &p); err != nil {
	//	log.Println(err)
	//}

	//if err := json.Unmarshal(res, &p); err != nil {
	//	log.Println(err)
	//}
	//err := json.Unmarshal([]byte(res.Body), p)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//log.Printf("b: %v", b)
	log.Printf("p: %v", p)

	//var p model.User
	//if err := json.Unmarshal(res.Body, &p); err != nil {
	//	fmt.Println(err)
	//}
	//log.Println(p)

	expected := `{"id":1,"createdAt":"2020-06-01T09:00:00+09:00","updatedAt":"2020-06-01T09:00:00+09:00","deletedAt":null,"name":"","email":"test@test.co.jp","age":0,"gender":0,"imageUrl":"","favoriteMusicAge":0,"favoriteArtist":"","comment":"","bookmarkings":null,"followings":null}`
	var u model.User
	var u2 model.User
	if err := json.Unmarshal([]byte(expected), u); err != nil {
		log.Println(err)
	}
	log.Println("u:", u)

	//if err := json.Unmarshal([]byte(res.Body.String()), u2); err != nil {
	if err := json.Unmarshal([]byte(res.Body.String()), u2); err != nil {
		log.Fatal(err)
	}

	log.Println("u2:", u2)

	//if res.Body.String() != expected {
	//	t.Errorf("handler returned unexpected body: got %v want %v",
	//		res.Body.String(), expected)
	//}
	//if res.Body != u {
	//	t.Errorf("handler returned unexpected body: got %v want %v",
	//		res.Body.String(), expected)
	//}

}

//リクエストユーザーの情報を返すハンドラのテスト
func TestUserHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/user", nil)

	//headerをセット
	req.Header.Set("Content-Type", "application/json")

	//リクエストユーザー作成
	user := model.User{Email: "a@test.co.jp", Password: "aaaaaa"}

	//トークン作成
	token, err := createToken(user)
	if err != nil {
		log.Println("err:", err)
	}

	jointToken := "Bearer" + " " + token

	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成
	res := httptest.NewRecorder()

	f := &UserController{UserInteractor: usecases.UserInteractor{
		UserRepository: &FakeUserRepository{},
	}}
	f.UserHandler(res, req)

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}

	// レスポンスのボディが期待通りか確認
	expected := `{"id":1,"createdAt":"2020-06-01T09:00:00+09:00","updatedAt":"2020-06-01T09:00:00+09:00","deletedAt":null,"name":"","email":"a@test.co.jp","age":0,"gender":0,"imageUrl":"","favoriteMusicAge":0,"favoriteArtist":"","comment":"","bookmarkings":null,"followings":null}`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			res.Body.String(), expected)
	}
}

//全てのユーザーの情報を返すハンドラのテスト
func TestAllUsersHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/users", nil)

	//headerをセット
	req.Header.Set("Content-Type", "application/json")

	//リクエストユーザー作成
	user := model.User{Email: "a@test.co.jp", Password: "aaaaaa"}

	//トークン作成
	token, err := createToken(user)
	if err != nil {
		log.Println("err:", err)
	}

	jointToken := "Bearer" + " " + token

	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成
	res := httptest.NewRecorder()

	f := &UserController{UserInteractor: usecases.UserInteractor{
		UserRepository: &FakeUserRepository{},
	}}
	f.AllUsersHandler(res, req)

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}

	// レスポンスのボディが期待通りか確認
	expected := `[{"id":1,"createdAt":"2020-06-01T09:00:00+09:00","updatedAt":"2020-06-01T09:00:00+09:00","deletedAt":null,"name":"","email":"a@test.co.jp","age":0,"gender":0,"imageUrl":"","favoriteMusicAge":0,"favoriteArtist":"","comment":"","bookmarkings":null,"followings":null},{"id":2,"createdAt":"2020-06-01T09:00:00+09:00","updatedAt":"2020-06-01T09:00:00+09:00","deletedAt":null,"name":"","email":"i@test.co.jp","age":0,"gender":0,"imageUrl":"","favoriteMusicAge":0,"favoriteArtist":"","comment":"","bookmarkings":null,"followings":null}]`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			res.Body.String(), expected)
	}
}

//idで指定したユーザーの情報を更新するハンドラのテスト
func TestUpdateUserHandler(t *testing.T) {
	// テスト用の JSON ボディ作成
	b, err := json.Marshal(model.User{Email: "hoge@test.co.jp", Name: "hogehoge", Age: 0, Gender: 0, FavoriteMusicAge: 0, FavoriteArtist: "椎名林檎", Comment: "テストユーザーです。"})
	if err != nil {
		t.Fatal(err)
	}

	// テスト用のリクエスト作成
	req := httptest.NewRequest("PUT", "/api/user/2/update", bytes.NewBuffer(b))

	//headerをセット
	req.Header.Set("Content-Type", "application/json")

	//リクエストユーザー作成
	user := model.User{Email: "a@test.co.jp", Password: "aaaaaa"}

	//トークン作成
	token, err := createToken(user)
	if err != nil {
		log.Println("err:", err)
	}

	jointToken := "Bearer" + " " + token

	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成
	res := httptest.NewRecorder()

	fakeUserController := &UserController{
		UserInteractor: usecases.UserInteractor{
			UserRepository: &FakeUserRepository{},
		},
	}

	r := mux.NewRouter()
	r.Handle("/api/user/{id}/update", http.HandlerFunc(fakeUserController.UpdateUserHandler)).Methods("PUT")
	r.ServeHTTP(res, req)

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}
}
