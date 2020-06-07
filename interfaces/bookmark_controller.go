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

type BookmarkController struct {
	BookmarkInteractor usecases.BookmarkInteractor
}

func NewBookmarkController(DB *gorm.DB) *BookmarkController {
	return &BookmarkController{
		BookmarkInteractor: usecases.BookmarkInteractor{
			BookmarkRepository: &BookmarkRepository{
				DB: DB,
			},
		},
	}
}

func (bc *BookmarkController) BookmarkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		var error model.Error
		error.Message = "idの取得に失敗しました"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}
	songID, err := strconv.Atoi(id)
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

	userEmail := parsedToken.Email

	if err := bc.BookmarkInteractor.Bookmark(userEmail, songID); err != nil {
		var error model.Error
		error.Message = "曲のお気に入り登録に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}

func (bc *BookmarkController) RemoveBookmarkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		var error model.Error
		error.Message = "idの取得に失敗しました"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}
	songID, err := strconv.Atoi(id)
	if err != nil {
		var error model.Error
		error.Message = "idのint型への型変換に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	header_hoge := r.Header.Get("Authorization")
	bearerToken := strings.Split(header_hoge, " ")
	authToken := bearerToken[1]

	parsedToken, err := Parse(authToken)
	if err != nil {
		var error model.Error
		error.Message = "認証コードのパースに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	userEmail := parsedToken.Email

	if err := bc.BookmarkInteractor.RemoveBookmark(userEmail, songID); err != nil {
		var error model.Error
		error.Message = "曲のお気に入り解除に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}
