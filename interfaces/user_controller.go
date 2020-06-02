package interfaces

import (
	"encoding/json"
	"golang-songs/model"
	"golang-songs/usecases"
	"net/http"

	"github.com/gorilla/mux"

	//"strconv"

	"github.com/jinzhu/gorm"
)

// A UserController belong to the interface layer.
type UserController struct {
	UserInteractor usecases.UserInteractor
	//Logger         usecases.Logger
}

// NewUserController returns the resource of users.
//func NewUserController(sqlHandler SQLHandler, logger usecases.Logger) *UserController {
func NewUserController(DB *gorm.DB) *UserController {
	return &UserController{
		UserInteractor: usecases.UserInteractor{
			UserRepository: &UserRepository{
				DB: DB,
			},
		},
		//Logger: logger,
	}
}

// Index return response which contain a listing of the resource of users.
func (uc *UserController) Index(w http.ResponseWriter, r *http.Request) {
	allUsers, err := uc.UserInteractor.Index()

	//以下は元のコード
	v, err := json.Marshal(allUsers)
	if err != nil {
		//var error model.Error
		//error.Message = "ユーザー一覧の取得に失敗しました"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
	if _, err := w.Write(v); err != nil {
		//var error model.Error
		//error.Message = "ユーザー一覧の取得に失敗しました。"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

}

// Show return response which contain the specified resource of a user.
func (uc *UserController) Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["id"]
	if !ok {
		var error model.Error
		error.Message = "ユーザーのidを取得できません。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	var user model.User

	user, err := uc.UserInteractor.Show(userID)

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

// Show return response which contain the specified resource of a user.
func (uc *UserController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["id"]
	if !ok {
		var error model.Error
		error.Message = "ユーザーのidを取得できません。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	dec := json.NewDecoder(r.Body)
	var d model.Song
	if err := dec.Decode(&d); err != nil {
		var error model.Error
		error.Message = "リクエストボディのデコードに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	var user model.User

	user, err := uc.UserInteractor.Update(userID)

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
