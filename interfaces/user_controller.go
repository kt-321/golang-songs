package interfaces

import (
	"encoding/json"
	"golang-songs/model"
	"golang-songs/usecases"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	//"strconv"

	"github.com/jinzhu/gorm"
)

type UserController struct {
	UserInteractor usecases.UserInteractor
}

func NewUserController(DB *gorm.DB) *UserController {
	return &UserController{
		UserInteractor: usecases.UserInteractor{
			UserRepository: &UserRepository{
				DB: DB,
			},
		},
	}
}

//全てのユーザーを返す
func (uc *UserController) AllUsersHandler(w http.ResponseWriter, r *http.Request) {
	allUsers, err := uc.UserInteractor.Index()
	if err != nil {
		var error model.Error
		error.Message = "曲が見つかりません。"
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
		//var error model.Error
		//error.Message = "ユーザー一覧の取得に失敗しました。"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

}

//リクエストユーザーの情報を返す
func (uc *UserController) UserHandler(w http.ResponseWriter, r *http.Request) {
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

	var user *model.User

	user, err = uc.UserInteractor.User(userEmail)
	if err != nil {
		var error model.Error
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
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

//idで指定したユーザーの情報を返す
func (uc *UserController) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		var error model.Error
		error.Message = "ユーザーのidを取得できません。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	userID, err := strconv.Atoi(id)
	if err != nil {
		var error model.Error
		error.Message = "idのint型への型変換に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	var user *model.User

	user, err = uc.UserInteractor.Show(userID)
	if err != nil {
		var error model.Error
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
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

//idで指定したユーザーの情報を更新する
func (uc *UserController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		var error model.Error
		error.Message = "ユーザーのidを取得できません。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}
	userID, err := strconv.Atoi(id)
	if err != nil {
		var error model.Error
		error.Message = "idのint型への型変換に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	var d model.User
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		var error model.Error
		error.Message = "リクエストボディのデコードに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	if err := uc.UserInteractor.Update(userID, d); err != nil {
		var error model.Error
		error.Message = "曲の更新に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	//204 No Content
	w.WriteHeader(204)
	return
}
