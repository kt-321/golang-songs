package interfaces

import (
	"golang-songs/model"
	"net/http"

	"github.com/jinzhu/gorm"
)

type SongRepository struct {
	//SQLHandler SQLHandler
	DB *gorm.DB
}

//func (sr *SongRepository) FindAll() (songs *model.Songs, err error) {
func (sr *SongRepository) FindAll(userEmail string) (*model.Songs, error) {
	//var user model.User
	////userEmailを引数で渡されるべきか
	//if err := sr.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
	//	error := model.Error{}
	//	error.Message = "該当するアカウントが見つかりません。"
	//	errorInResponse(w, http.StatusUnauthorized, error)
	//	return
	//}

	songs := []model.Song{}

	if err := sr.DB.Find(&songs).Error; gorm.IsRecordNotFoundError(err) {
		var error model.Error
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
	return songs, err
}

// FindByID is returns the entity identified by the given id.
func (sr *SongRepository) FindByID(userEmail string, songID int) (*model.Song, error) {
	var user model.User
	//userEmailを引数で渡されるべきか
	if err := sr.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	var song model.Song

	if err := sr.DB.Where("id = ?", songID).Find(&song).Error; gorm.IsRecordNotFoundError(err) {
		var error model.Error
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	return song, err
}

// Save is saves the given entity.
//func (sr *SongRepository) Save(p model.Song) (id int64, err error) {
func (sr *SongRepository) Save(userEmail string, p model.Song) (int64, error) {

	var user model.User
	//userEmailを引数で渡されるべきか
	if err := sr.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
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
		var error model.Error
		error.Message = "曲の追加に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
	//result, err := tx.Exec(query, p.Title, p.Artist, p.MusicAge, p.Image, p.Video, p.Album, p.Description, p.SpotifyTrackId, p.UserID)
	//if err != nil {
	//	_ = tx.Rollback()
	//	return
	//}
	//
	//if err = tx.Commit(); err != nil {
	//	return
	//}
	//
	//id, err = result.LastInsertId()
	//if err != nil {
	//	return id, nil
	//}

	return id, nil
}

// DeleteByID is deletes the entity identified by the given id.
func (sr *SongRepository) UpdateByID(userEmail string, songID int, p model.Song) (err error) {
	var user model.User
	//userEmailを引数で渡されるべきか
	if err := sr.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
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
		var error model.Error
		error.Message = "曲の更新に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	return
}

// DeleteByID is deletes the entity identified by the given id.
//func (sr *SongRepository) DeleteByID(songID int) (err error) {
func (sr *SongRepository) DeleteByID(userEmail string, songID int) (err error) {
	var user model.User
	//userEmailを引数で渡されるべきか
	if err := sr.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		error := model.Error{}
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}
	var song model.Song

	if err := sr.DB.Where("id = ?", songID).Delete(&song).Error; err != nil {
		var error model.Error
		error.Message = "曲の削除に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	return
}
