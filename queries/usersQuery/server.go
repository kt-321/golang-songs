package usersQuery

import (
	"encoding/json"
	"golang-songs/interfaces"
	"golang-songs/model"

	//"golang-songs/queries/usersQuery"
	//"golang-songs/usecases"

	"net/http"

	"github.com/jinzhu/gorm"
)

type usersQueryServer struct {
	usecase usecase
}

func NewUserQueryServer(DB *gorm.DB) *usersQueryServer {
	return &usersQueryServer{
		usecase: usecase{
			da: dataAccessor{
				DB: DB,
			},
		},
	}
}

// 全てのユーザーを返す.
func (uc *usersQueryServer) AllUsers(w http.ResponseWriter, r *http.Request) {
	allUsers, err := uc.usecase.Index()

	if err != nil {
		interfaces.ErrorInResponse(w, http.StatusInternalServerError, interfaces.GetUserError)

		return
	}

	v, err := json.Marshal(allUsers)

	if err != nil {
		interfaces.ErrorInResponse(w, http.StatusInternalServerError, interfaces.JsonEncodeError)

		return
	}

	if _, err := w.Write(v); err != nil {
		interfaces.ErrorInResponse(w, http.StatusInternalServerError, interfaces.GetUsersListError)

		return
	}
}

// リクエストユーザーの情報を返す.
func (uc *usersQueryServer) User(w http.ResponseWriter, r *http.Request) {
	// リクエストユーザーのメアドを取得.
	userEmail, errorSet := interfaces.GetEmail(r)

	if errorSet != nil {
		interfaces.ErrorInResponse(w, errorSet.StatusCode, errorSet.Message)

		return
	}

	var user *model.User

	user, err := uc.usecase.User(userEmail)

	if err != nil {
		interfaces.ErrorInResponse(w, http.StatusInternalServerError, interfaces.GetAccountError)

		return
	}

	v, err := json.Marshal(user)

	if err != nil {
		interfaces.ErrorInResponse(w, http.StatusInternalServerError, interfaces.JsonEncodeError)

		return
	}

	if _, err := w.Write(v); err != nil {
		interfaces.ErrorInResponse(w, http.StatusInternalServerError, interfaces.GetUserDetailError)

		return
	}
}

// idで指定したユーザーの情報を返す.
func (uc *usersQueryServer) GetUser(w http.ResponseWriter, r *http.Request) {
	// 対象のユーザーidを取得.
	userID, errorSet := interfaces.GetId(r)

	if errorSet != nil {
		interfaces.ErrorInResponse(w, errorSet.StatusCode, errorSet.Message)

		return
	}

	user, err := uc.usecase.Show(userID)

	if err != nil {
		interfaces.ErrorInResponse(w, http.StatusInternalServerError, interfaces.GetAccountError)

		return
	}

	v, err := json.Marshal(user)

	if err != nil {
		interfaces.ErrorInResponse(w, http.StatusInternalServerError, interfaces.JsonEncodeError)

		return
	}

	if _, err := w.Write(v); err != nil {
		interfaces.ErrorInResponse(w, http.StatusInternalServerError, interfaces.GetUserDetailError)

		return
	}
}

// idで指定したユーザーの情報を更新する.
//func (uc *UserController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
//	// 対象のユーザーidを取得.
//	userID, errorSet := GetId(r)
//
//	if errorSet != nil {
//		errorInResponse(w, errorSet.StatusCode, errorSet.Message)
//
//		return
//	}
//
//	var d model.User
//
//	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
//		errorInResponse(w, http.StatusInternalServerError, DecodeError)
//
//		return
//	}
//
//	if err := uc.UserInteractor.Update(userID, d); err != nil {
//		errorInResponse(w, http.StatusInternalServerError, UpdateUserError)
//
//		return
//	}
//
//	// 204 No Content.
//	w.WriteHeader(204)
//}
