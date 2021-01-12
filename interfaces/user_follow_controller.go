package interfaces

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"golang-songs/usecases"
	"net/http"
	"strconv"
)

type UserFollowController struct {
	UserFollowInteractor usecases.UserFollowInteractor
}

func NewUserFollowController(DB *gorm.DB) *UserFollowController {
	return &UserFollowController{
		UserFollowInteractor: usecases.UserFollowInteractor{
			UserFollowRepository: &UserFollowRepository{
				DB: DB,
			},
		},
	}
}

// idで指定したユーザーをフォローする.
func (ufc *UserFollowController) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	requestUserEmail, errorSet := GetEmail(r)

	if errorSet != nil {
		errorInResponse(w, errorSet.StatusCode, errorSet.Message)

		return
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		errorInResponse(w, http.StatusBadRequest, GetUserIdError)

		return
	}

	targetUserID, err := strconv.Atoi(id)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, ConvertIdToIntError)

		return
	}

	if err := ufc.UserFollowInteractor.Follow(requestUserEmail, targetUserID); err != nil {
		errorInResponse(w, http.StatusInternalServerError, FollowUserError)

		return
	}

	// 201 Created.
	w.WriteHeader(201)
}

// idで指定したユーザーのフォローを解除する.
func (ufc *UserFollowController) UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	requestUserEmail, errorSet := GetEmail(r)

	if errorSet != nil {
		errorInResponse(w, errorSet.StatusCode, errorSet.Message)

		return
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		errorInResponse(w, http.StatusBadRequest, GetUserIdError)

		return
	}

	targetUserID, err := strconv.Atoi(id)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, ConvertIdToIntError)

		return
	}

	if err := ufc.UserFollowInteractor.Unfollow(requestUserEmail, targetUserID); err != nil {
		errorInResponse(w, http.StatusInternalServerError, UnfollowUserError)

		return
	}

	// 201 Created.
	w.WriteHeader(201)
}
