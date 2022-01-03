package userQuery

import (
	"encoding/json"
	"golang-songs/interfaces"
	"golang-songs/model"

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

// 全てのユーザーを返す.
func (uqs *userQueryServer) AllUsers(w http.ResponseWriter, r *http.Request) {
	allUsers, err := uqs.usecase.Index()

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
func (uqs *userQueryServer) User(w http.ResponseWriter, r *http.Request) {
	// リクエストユーザーのメアドを取得.
	userEmail, errorSet := interfaces.GetEmail(r)

	if errorSet != nil {
		interfaces.ErrorInResponse(w, errorSet.StatusCode, errorSet.Message)

		return
	}

	var user *model.User

	user, err := uqs.usecase.User(userEmail)

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

	user, err := uqs.usecase.Show(userID)

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
