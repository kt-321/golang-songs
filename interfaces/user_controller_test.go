package interfaces

import (
	"golang-songs/infrastructure"
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
	var users []model.User

	var user1 model.User

	user1.ID = 1
	user1.CreatedAt = time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local)
	user1.UpdatedAt = time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local)
	user1.Password = "aaaaaa"
	user1.Name = ""
	user1.Email = "a@test.co.jp"
	user1.Age = 0
	user1.Gender = 0
	user1.ImageUrl = ""
	user1.FavoriteMusicAge = 0
	user1.FavoriteArtist = ""
	user1.Comment = ""

	var user2 model.User

	user2.ID = 2
	user2.CreatedAt = time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local)
	user2.UpdatedAt = time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local)
	user2.Password = "iiiiii"
	user2.Name = ""
	user2.Email = "i@test.co.jp"
	user2.Age = 0
	user2.Gender = 0
	user2.ImageUrl = ""
	user2.FavoriteMusicAge = 0
	user2.FavoriteArtist = ""
	user2.Comment = ""

	users = append(users, user1, user2)

	return &users, nil
}

func (fur *FakeUserRepository) GetUser(userEmail string) (*model.User, error) {
	var user model.User

	user.ID = 1
	user.CreatedAt = time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local)
	user.UpdatedAt = time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local)
	//user.DeletedAt = null
	user.Password = "aaaaaa"
	user.Name = ""
	user.Email = "a@test.co.jp"
	user.Age = 0
	user.Gender = 0
	user.ImageUrl = ""
	user.FavoriteMusicAge = 0
	user.FavoriteArtist = ""
	user.Comment = ""

	return &user, nil
}

func (fur *FakeUserRepository) FindByID(userID int) (*model.User, error) {
	var user model.User
	user.ID = 1
	user.CreatedAt = time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local)
	user.UpdatedAt = time.Date(2020, 6, 1, 9, 0, 0, 0, time.Local)
	user.Password = "uuuuuu"
	user.Name = ""
	user.Email = "u@u"
	user.Age = 0
	user.Gender = 0
	user.ImageUrl = ""
	user.FavoriteMusicAge = 0
	user.FavoriteArtist = ""
	user.Comment = ""
	return &user, nil
}

func (fur *FakeUserRepository) Update(userID int, p model.User) error {
	return nil
}

func TestGetUserHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/user/1", nil)

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

	log.Println("++++++++++++")
	log.Println("req")
	log.Println(req)
	log.Println("++++++++++++")

	// テスト用のレスポンス作成
	res := httptest.NewRecorder()

	//f := &UserController{UserInteractor: usecases.UserInteractor{
	//
	//	UserRepository: &FakeUserRepository{},
	//}}
	//f.GetUserHandler(res, req)

	//レシーバ無しの場合
	r := mux.NewRouter()
	r.Handle("/api/user/{id}", infrastructure.JwtMiddleware.Handler(http.HandlerFunc(infrastructure.UserController.GetUserHandler))).Methods("GET")
	//r.Handle("/api/user/{id}", infrastructure.JwtMiddleware.Handler(http.HandlerFunc(infrastucture.userController.GetUserHandler))).Methods("GET")
	//r.Handle("/api/user/{id}", &UserController{UserInteractor: usecases.UserInteractor{
	//	UserRepository: &FakeUserRepository{},
	//}})
	r.GetUserHandler(res, req)

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}

	expected := `{"id":1,"createdAt":"2020-06-01T09:00:00+09:00","updatedAt":"2020-06-01T09:00:00+09:00","deletedAt":null,"name":"","email":"u@u","age":0,"gender":0,"imageUrl":"","favoriteMusicAge":0,"favoriteArtist":"","comment":"","bookmarkings":null,"followings":null}`
	if res.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			res.Body.String(), expected)
	}

}

func TestUserHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/user", nil)

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

	jointToken := "Bearer" + " " + token

	req.Header.Set("Authorization", jointToken)

	// テスト用のレスポンス作成
	res := httptest.NewRecorder()

	//レシーバ付きの場合
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

func TestAllUsersHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/users", nil)

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

//func TestUpdateUserHandler(t *testing.T) {
//	// テスト用の JSON ボディ作成
//	b, err := json.Marshal(model.User{Email: "hoge@test.co.jp", Name: "hogehoge", Age: 0, Gender: 0, FavoriteMusicAge: 0, FavoriteArtist: "椎名林檎", Comment: "テストユーザーです。"})
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// テスト用のリクエスト作成
//	req := httptest.NewRequest("PUT", "/api/user/1/update", bytes.NewBuffer(b))
//
//	//headerをセット
//	req.Header.Set("Content-Type", "application/json")
//
//	var user model.User
//	user.Email = "a@test.co.jp"
//	user.Password = "aaaaaa"
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
//	f := &UserController{UserInteractor: usecases.UserInteractor{
//		UserRepository: &FakeUserRepository{},
//	}}
//	f.UpdateUserHandler(res, req)
//
//	log.Printf("req: %v", req)
//	log.Printf("res: %v", res)
//
//	// レスポンスのステータスコードのテスト
//	if res.Code != http.StatusOK {
//		t.Errorf("invalid code: %d", res.Code)
//	}
//
//}
