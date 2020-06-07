package interfaces

import (
	"golang-songs/model"

	"github.com/jinzhu/gorm"
)

type SongRepository struct {
	DB *gorm.DB
}

func (sr *SongRepository) FindAll() (*model.Songs, error) {
	var songs model.Songs

	if err := sr.DB.Find(&songs).Error; gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &songs, nil
}

func (sr *SongRepository) FindByID(songID int) (*model.Song, error) {
	var song model.Song

	if err := sr.DB.Where("id = ?", songID).Find(&song).Error; gorm.IsRecordNotFoundError(err) {
		//var error model.Error
		//error.Message = "該当するアカウントが見つかりません。"
		//errorInResponse(w, http.StatusUnauthorized, error)
		return nil, err
	}

	return &song, nil
}

func (sr *SongRepository) Save(userEmail string, p model.Song) error {

	var user model.User
	if err := sr.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		//error := model.Error{}
		//error.Message = "該当するアカウントが見つかりません。"
		//errorInResponse(w, http.StatusUnauthorized, error)
		return err
	}

	if err := sr.DB.Create(&model.Song{
		Title:          p.Title,
		Artist:         p.Artist,
		MusicAge:       p.MusicAge,
		Image:          p.Image,
		Video:          p.Video,
		Album:          p.Album,
		Description:    p.Description,
		SpotifyTrackId: p.SpotifyTrackId,
		UserID:         user.ID}).Error; err != nil {
		//var error model.Error
		//error.Message = "曲の追加に失敗しました"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return err
	}

	return nil
}

func (sr *SongRepository) UpdateByID(userEmail string, songID int, p model.Song) error {
	var user model.User
	if err := sr.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		//error := model.Error{}
		//error.Message = "該当するアカウントが見つかりません。"
		//errorInResponse(w, http.StatusUnauthorized, error)
		return err
	}

	var song model.Song

	if err := sr.DB.Model(&song).Where("id = ?", songID).Update(model.Song{
		Title:          p.Title,
		Artist:         p.Artist,
		MusicAge:       p.MusicAge,
		Image:          p.Image,
		Video:          p.Video,
		Album:          p.Album,
		Description:    p.Description,
		SpotifyTrackId: p.SpotifyTrackId}).Error; err != nil {
		//var error model.Error
		//error.Message = "曲の更新に失敗しました"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return err
	}

	return nil
}

func (sr *SongRepository) DeleteByID(songID int) error {
	var song model.Song

	if err := sr.DB.Where("id = ?", songID).Delete(&song).Error; err != nil {
		//var error model.Error
		//error.Message = "曲の削除に失敗しました"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return err
	}

	return nil
}
