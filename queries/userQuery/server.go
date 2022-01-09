package userQuery

import (
	"encoding/json"
	"golang-songs/interfaces"

	//"golang-songs/queries/userQuery"
	//"golang-songs/usecases"

	"net/http"

	"github.com/jmoiron/sqlx"
)

type userQueryServer struct {
	usecase usecase
}

func NewUserQueryServer(DB *sqlx.DB) *userQueryServer {
	return &userQueryServer{
		usecase: usecase{
			da: &dataAccessor{
				DB: DB,
			},
		},
	}
}

//TODO interface

// 全てのユーザーを返す.
func (uqs *userQueryServer) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	allUsers, err := uqs.usecase.GetAllUsers()

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
func (uqs *userQueryServer) GetAuthUser(w http.ResponseWriter, r *http.Request) {
	// リクエストユーザーのメアドを取得.
	userEmail, errorSet := interfaces.GetEmail(r)

	if errorSet != nil {
		interfaces.ErrorInResponse(w, errorSet.StatusCode, errorSet.Message)

		return
	}

	user, err := uqs.usecase.FindUserByEmail(userEmail)

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
func (uqs *userQueryServer) GetUser(w http.ResponseWriter, r *http.Request) {
	// 対象のユーザーidを取得.
	userID, errorSet := interfaces.GetId(r)

	if errorSet != nil {
		interfaces.ErrorInResponse(w, errorSet.StatusCode, errorSet.Message)

		return
	}

	user, err := uqs.usecase.FindUserByID(userID)

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
