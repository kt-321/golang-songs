package interfaces

import (
	"golang-songs/model"

	"github.com/jinzhu/gorm"
)

type BookmarkRepository struct {
	DB *gorm.DB
}

//曲をお気に入り登録
func (br *BookmarkRepository) Bookmark(userEmail string, songID int) error {

	var user model.User
	if err := br.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}

	var song model.Song
	if err := br.DB.Where("id = ?", songID).Find(&song).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}

	if err := br.DB.Create(&model.Bookmark{
		UserID: user.ID,
		SongID: song.ID}).Error; err != nil {
		return err
	}

	if err := br.DB.Preload("Bookmarkings").Find(&user).Error; err != nil {
		return err
	}

	if err := br.DB.Model(&user).Association("Bookmarkings").Append(&song).Error; err != nil {
		return err
	}

	return nil
}

//曲をお気に入り登録から解除
func (br *BookmarkRepository) RemoveBookmark(userEmail string, songID int) error {
	var user model.User
	if err := br.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}

	var song model.Song

	if err := br.DB.Where("id = ?", songID).Find(&song).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}

	if err := br.DB.Preload("Bookmarkings").Find(&user).Error; err != nil {
		return err
	}

	if err := br.DB.Model(&user).Association("Bookmarkings").Delete(&song).Error; err != nil {
		return err
	}

	return nil
}
