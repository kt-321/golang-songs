package interfaces

import (
	"golang-songs/usecases"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
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
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		errorInResponse(w, http.StatusBadRequest, 26)

		return
	}

	targetUserID, err := strconv.Atoi(id)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 14)

		return
	}

	headerAuthorization := r.Header.Get("Authorization")

	if len(headerAuthorization) == 0 {
		errorInResponse(w, http.StatusInternalServerError, 28)

		return
	}

	bearerToken := strings.Split(headerAuthorization, " ")

	if len(bearerToken) < 2 {
		errorInResponse(w, http.StatusUnauthorized, 29)

		return
	}

	authToken := bearerToken[1]

	parsedToken, err := Parse(authToken)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 18)

		return
	}

	requestUserEmail := parsedToken.Email

	if err := ufc.UserFollowInteractor.Follow(requestUserEmail, targetUserID); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 30)

		return
	}

	// 201 Created.
	w.WriteHeader(201)
}

// idで指定したユーザーのフォローを解除する.
func (ufc *UserFollowController) UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		errorInResponse(w, http.StatusBadRequest, 26)

		return
	}

	targetUserID, err := strconv.Atoi(id)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 14)

		return
	}

	headerAuthorization := r.Header.Get("Authorization")

	if len(headerAuthorization) == 0 {
		errorInResponse(w, http.StatusInternalServerError, 28)

		return
	}

	bearerToken := strings.Split(headerAuthorization, " ")

	if len(bearerToken) < 2 {
		errorInResponse(w, http.StatusUnauthorized, 29)

		return
	}

	authToken := bearerToken[1]

	parsedToken, err := Parse(authToken)
	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 18)

		return
	}

	requestUserEmail := parsedToken.Email

	if err := ufc.UserFollowInteractor.Unfollow(requestUserEmail, targetUserID); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 31)

		return
	}

	// 201 Created.
	w.WriteHeader(201)
}
