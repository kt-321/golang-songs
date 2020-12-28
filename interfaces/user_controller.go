package interfaces

import (
	"encoding/json"
	"golang-songs/model"
	"golang-songs/usecases"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

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
		errorInResponse(w, http.StatusInternalServerError, 22)

		return
	}

	v, err := json.Marshal(allUsers)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 23)

		return
	}

	if _, err := w.Write(v); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 23)

		return
	}
}

// リクエストユーザーの情報を返す.
func (uc *UserController) UserHandler(w http.ResponseWriter, r *http.Request) {
	headerHoge := r.Header.Get("Authorization")
	bearerToken := strings.Split(headerHoge, " ")
	authToken := bearerToken[1]

	parsedToken, err := Parse(authToken)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 18)

		return
	}

	userEmail := parsedToken.Email

	var user *model.User

	user, err = uc.UserInteractor.User(userEmail)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 24)

		return
	}

	v, err := json.Marshal(user)
	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 6)

		return
	}

	if _, err := w.Write(v); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 25)

		return
	}
}

// idで指定したユーザーの情報を返す.
func (uc *UserController) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		errorInResponse(w, http.StatusBadRequest, 26)

		return
	}

	userID, err := strconv.Atoi(id)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 14)

		return
	}

	var user *model.User

	user, err = uc.UserInteractor.Show(userID)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 24)

		return
	}

	v, err := json.Marshal(user)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 6)

		return
	}

	if _, err := w.Write(v); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 25)

		return
	}
}

// idで指定したユーザーの情報を更新する.
func (uc *UserController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		errorInResponse(w, http.StatusBadRequest, 26)

		return
	}

	userID, err := strconv.Atoi(id)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 14)

		return
	}

	var d model.User

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 17)

		return
	}

	if err := uc.UserInteractor.Update(userID, d); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 27)

		return
	}

	// 204 No Content.
	w.WriteHeader(204)
}
