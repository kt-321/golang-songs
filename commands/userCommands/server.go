package userCommands

import (
	"encoding/json"
	"golang-songs/domain"
	"golang-songs/interfaces"
	"golang-songs/model"
	"net/http"

	"github.com/jinzhu/gorm"
)

//type UserController struct {
//	UserInteractor usecases.UserInteractor
//}

type userCommandsServer struct {
	usecase usecase
}

func NewUserController(DB *gorm.DB) *userCommandsServer {
	return &userCommandsServer{
		usecase: usecase{
			UserRepository: &domain.UserRepository{
				Da: domain.UserDataAccessor{DB: DB},
			},
		},
	}
}

// idで指定したユーザーの情報を更新する.
func (uc *userCommandsServer) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// 対象のユーザーidを取得.
	userID, errorSet := interfaces.GetId(r)

	if errorSet != nil {
		interfaces.ErrorInResponse(w, errorSet.StatusCode, errorSet.Message)

		return
	}

	var d model.User

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		interfaces.ErrorInResponse(w, http.StatusInternalServerError, interfaces.DecodeError)

		return
	}

	if err := uc.usecase.UserRepository.Update(userID, d); err != nil {
		interfaces.ErrorInResponse(w, http.StatusInternalServerError, interfaces.UpdateUserError)

		return
	}

	// 204 No Content.
	w.WriteHeader(204)
}
