package interfaces

import (
	"golang-songs/model"

	"github.com/jinzhu/gorm"
)

type BookmarkRepository struct {
	DB *gorm.DB
}

func (br *BookmarkRepository) Bookmark(userEmail string, songID int) error {

	var user model.User
	if err := br.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		//var error model.Error
		//error.Message = "該当するアカウントが見つかりません。"
		//errorInResponse(w, http.StatusUnauthorized, error)
		return err
	}

	var song model.Song
	if err := br.DB.Where("id = ?", songID).Find(&song).Error; gorm.IsRecordNotFoundError(err) {
		//error := model.Error{}
		//error.Message = "該当する曲が見つかりません。"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return err
	}

	if err := br.DB.Create(&model.Bookmark{
		UserID: user.ID,
		SongID: song.ID}).Error; err != nil {
		//var error model.Error
		//error.Message = "曲のお気に入り登録に失敗しました"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return err
	}

	if err := br.DB.Preload("Bookmarkings").Find(&user).Error; err != nil {
		//var error model.Error
		//error.Message = "該当する参照が見つかりません。"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return err
	}

	if err := br.DB.Model(&user).Association("Bookmarkings").Append(&song).Error; err != nil {
		//error := model.Error{}
		//error.Message = "参照の追加に失敗しました。"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return err
	}

	return nil
}

func (br *BookmarkRepository) RemoveBookmark(userEmail string, songID int) error {
	var user model.User
	if err := br.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}

	var song model.Song

	if err := br.DB.Where("id = ?", songID).Find(&song).Error; gorm.IsRecordNotFoundError(err) {
		//error := model.Error{}
		//error.Message = "該当する曲が見つかりません。"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return err
	}

	if err := br.DB.Preload("Bookmarkings").Find(&user).Error; err != nil {
		//var error model.Error
		//error.Message = "該当する参照が見つかりません。"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return err
	}

	if err := br.DB.Model(&user).Association("Bookmarkings").Delete(&song).Error; err != nil {
		//error := model.Error{}
		//error.Message = "参照の削除に失敗しました。"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return err
	}

	return nil
}
