package interfaces

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"golang-songs/model"
	"golang-songs/usecases"
	"net/http"
	"strconv"

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

// 全てのユーザーを返す.
func (uc *UserController) AllUsersHandler(w http.ResponseWriter, r *http.Request) {
	allUsers, err := uc.UserInteractor.Index()

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, GetUserError)

		return
	}

	v, err := json.Marshal(allUsers)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, JsonEncodeError)

		return
	}

	if _, err := w.Write(v); err != nil {
		errorInResponse(w, http.StatusInternalServerError, GetUsersListError)

		return
	}
}

// リクエストユーザーの情報を返す.
func (uc *UserController) UserHandler(w http.ResponseWriter, r *http.Request) {
	userEmail, errorSet := GetEmail(r)

	if errorSet != nil {
		errorInResponse(w, errorSet.StatusCode, errorSet.Message)

		return
	}

	var user *model.User

	user, err := uc.UserInteractor.User(userEmail)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, GetAccountError)

		return
	}

	v, err := json.Marshal(user)
	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, JsonEncodeError)

		return
	}

	if _, err := w.Write(v); err != nil {
		errorInResponse(w, http.StatusInternalServerError, GetUserDetailError)

		return
	}
}

// idで指定したユーザーの情報を返す.
func (uc *UserController) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		errorInResponse(w, http.StatusBadRequest, GetUserIdError)

		return
	}

	userID, err := strconv.Atoi(id)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, ConvertIdToIntError)

		return
	}

	var user *model.User

	user, err = uc.UserInteractor.Show(userID)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, GetAccountError)

		return
	}

	v, err := json.Marshal(user)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, JsonEncodeError)

		return
	}

	if _, err := w.Write(v); err != nil {
		errorInResponse(w, http.StatusInternalServerError, GetUserDetailError)

		return
	}
}

// idで指定したユーザーの情報を更新する.
func (uc *UserController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		errorInResponse(w, http.StatusBadRequest, GetUserIdError)

		return
	}

	userID, err := strconv.Atoi(id)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, ConvertIdToIntError)

		return
	}

	var d model.User

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		errorInResponse(w, http.StatusInternalServerError, DecodeError)

		return
	}

	if err := uc.UserInteractor.Update(userID, d); err != nil {
		errorInResponse(w, http.StatusInternalServerError, UpdateUserError)

		return
	}

	// 204 No Content.
	w.WriteHeader(204)
}
