package interfaces

import (
	"golang-songs/usecases"
	"net/http"

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
	// リクエストユーザーのメアドと対象の曲idを取得.
	userEmail, songID, errorSet := GetEmailAndId(r)

	if errorSet != nil {
		ErrorInResponse(w, errorSet.StatusCode, errorSet.Message)

		return
	}

	if err := bc.BookmarkInteractor.Bookmark(userEmail, songID); err != nil {
		ErrorInResponse(w, http.StatusInternalServerError, BookmarkSongError)

		return
	}

	// 201 Created
	w.WriteHeader(201)
}

// 曲をお気に入り登録から解除.
func (bc *BookmarkController) RemoveBookmarkHandler(w http.ResponseWriter, r *http.Request) {
	// リクエストユーザーのメアドと対象の曲idを取得.
	userEmail, songID, errorSet := GetEmailAndId(r)

	if errorSet != nil {
		ErrorInResponse(w, errorSet.StatusCode, errorSet.Message)

		return
	}

	if err := bc.BookmarkInteractor.RemoveBookmark(userEmail, songID); err != nil {
		ErrorInResponse(w, http.StatusInternalServerError, RemoveBookmarkError)

		return
	}

	// 201 Created
	w.WriteHeader(201)
}
