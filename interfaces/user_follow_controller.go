package interfaces

import (
	"golang-songs/model"
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

//idで指定したユーザーをフォローする
func (ufc *UserFollowController) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		var error model.Error
		error.Message = "ユーザーのidを取得できません。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	targetUserID, err := strconv.Atoi(id)
	if err != nil {
		var error model.Error
		error.Message = "idのint型への型変換に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	headerAuthorization := r.Header.Get("Authorization")
	if len(headerAuthorization) == 0 {
		var error model.Error
		error.Message = "認証トークンの取得に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	bearerToken := strings.Split(headerAuthorization, " ")
	if len(bearerToken) < 2 {
		var error model.Error
		error.Message = "bearerトークンの取得に失敗しました。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	authToken := bearerToken[1]

	parsedToken, err := Parse(authToken)
	if err != nil {
		var error model.Error
		error.Message = "認証コードのパースに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	requestUserEmail := parsedToken.Email

	if err := ufc.UserFollowInteractor.Follow(requestUserEmail, targetUserID); err != nil {
		var error model.Error
		error.Message = "ユーザーのフォローに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	//201 Created
	w.WriteHeader(201)
}

//idで指定したユーザーのフォローを解除する
func (ufc *UserFollowController) UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		var error model.Error
		error.Message = "ユーザーのidを取得できません。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	targetUserID, err := strconv.Atoi(id)
	if err != nil {
		var error model.Error
		error.Message = "idのint型への型変換に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	headerAuthorization := r.Header.Get("Authorization")
	if len(headerAuthorization) == 0 {
		var error model.Error
		error.Message = "認証トークンの取得に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	bearerToken := strings.Split(headerAuthorization, " ")
	if len(bearerToken) < 2 {
		var error model.Error
		error.Message = "bearerトークンの取得に失敗しました。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	authToken := bearerToken[1]

	parsedToken, err := Parse(authToken)
	if err != nil {
		var error model.Error
		error.Message = "認証コードのパースに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	requestUserEmail := parsedToken.Email

	if err := ufc.UserFollowInteractor.Unfollow(requestUserEmail, targetUserID); err != nil {
		var error model.Error
		error.Message = "ユーザーのフォロー解除に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	//201 Created
	w.WriteHeader(201)
}
