package interfaces

import (
	"golang-songs/usecases"
	"net/http"

	"github.com/jinzhu/gorm"
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
	// リクエストユーザーのメアドと対象のユーザーidを取得.
	requestUserEmail, targetUserID, errorSet := GetEmailAndId(r)

	if errorSet != nil {
		ErrorInResponse(w, errorSet.StatusCode, errorSet.Message)

		return
	}

	if err := ufc.UserFollowInteractor.Follow(requestUserEmail, targetUserID); err != nil {
		ErrorInResponse(w, http.StatusInternalServerError, FollowUserError)

		return
	}

	// 201 Created.
	w.WriteHeader(201)
}

// idで指定したユーザーのフォローを解除する.
func (ufc *UserFollowController) UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	// リクエストユーザーのメアドと対象のユーザーidを取得.
	requestUserEmail, targetUserID, errorSet := GetEmailAndId(r)

	if errorSet != nil {
		ErrorInResponse(w, errorSet.StatusCode, errorSet.Message)

		return
	}

	if err := ufc.UserFollowInteractor.Unfollow(requestUserEmail, targetUserID); err != nil {
		ErrorInResponse(w, http.StatusInternalServerError, UnfollowUserError)

		return
	}

	// 201 Created.
	w.WriteHeader(201)
}
