package interfaces

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"golang-songs/usecases"
	"net/http"
	"strconv"
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
	userEmail, errorSet := GetEmail(r)

	if errorSet != nil {
		errorInResponse(w, errorSet.StatusCode, errorSet.Message)

		return
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		errorInResponse(w, http.StatusBadRequest, GetSongIdError)

		return
	}

	songID, err := strconv.Atoi(id)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, ConvertIdToIntError)

		return
	}

	if err := bc.BookmarkInteractor.Bookmark(userEmail, songID); err != nil {
		errorInResponse(w, http.StatusInternalServerError, BookmarkSongError)

		return
	}

	// 201 Created
	w.WriteHeader(201)
}

// 曲をお気に入り登録から解除.
func (bc *BookmarkController) RemoveBookmarkHandler(w http.ResponseWriter, r *http.Request) {
	userEmail, errorSet := GetEmail(r)

	if errorSet != nil {
		errorInResponse(w, errorSet.StatusCode, errorSet.Message)

		return
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		errorInResponse(w, http.StatusBadRequest, GetSongIdError)

		return
	}

	songID, err := strconv.Atoi(id)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, ConvertIdToIntError)

		return
	}

	if err := bc.BookmarkInteractor.RemoveBookmark(userEmail, songID); err != nil {
		errorInResponse(w, http.StatusInternalServerError, RemoveBookmarkError)

		return
	}

	// 201 Created
	w.WriteHeader(201)
}
