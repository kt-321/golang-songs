package interfaces

import (
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

// 曲をお気に入りに登録.
func (bc *BookmarkController) BookmarkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		errorInResponse(w, http.StatusBadRequest, 13)

		return
	}

	songID, err := strconv.Atoi(id)

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

	userEmail := parsedToken.Email

	if err := bc.BookmarkInteractor.Bookmark(userEmail, songID); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 32)

		return
	}

	// 201 Created
	w.WriteHeader(201)
}

// 曲をお気に入り登録から解除.
func (bc *BookmarkController) RemoveBookmarkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		errorInResponse(w, http.StatusBadRequest, 13)

		return
	}

	songID, err := strconv.Atoi(id)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 14)

		return
	}

	headerHoge := r.Header.Get("Authorization")
	bearerToken := strings.Split(headerHoge, " ")
	authToken := bearerToken[1]

	parsedToken, err := Parse(authToken)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 18)

		return
	}

	userEmail := parsedToken.Email

	if err := bc.BookmarkInteractor.RemoveBookmark(userEmail, songID); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 33)

		return
	}

	// 201 Created
	w.WriteHeader(201)
}
