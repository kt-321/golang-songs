package interfaces

import (
	"golang-songs/model"
	"net/http"

	//"time"

	"github.com/jinzhu/gorm"
)

// A UserRepository belong to the inteface layer
type UserRepository struct {
	DB *gorm.DB
}

// FindAll is returns the number of entities.
//名前付き戻り値は良くない
//func (ur *UserRepository) FindAll() (users domain.Users, err error) {
func (ur *UserRepository) FindAll() (*model.Users, error) {
	//var users []model.User
	var users model.Users
	if err := ur.DB.Find(&users).Error; gorm.IsRecordNotFoundError(err) {
		var error model.Error
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
	return users, err
}

// FindByID is returns the entity identified by the given id.
func (ur *UserRepository) FindByID(userID int) (*model.User, error) {
	var user model.User

	if err := ur.DB.Where("id = ?", userID).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		var error model.Error
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	var bookmarkings []model.Song

	if err := ur.DB.Preload("Bookmarkings").Find(&user).Error; err != nil {
		var error model.Error
		error.Message = "該当する参照が見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	if ur.DB.Model(&user).Related(&bookmarkings, "Bookmarikings").RecordNotFound() {
		error := model.Error{}
		error.Message = "レコードが見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	var followings []model.User
	if err := ur.DB.Preload("Followings").Find(&user).Error; err != nil {
		var error model.Error
		error.Message = "該当する参照が見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
	if ur.DB.Model(&user).Related(&followings, "Followings").RecordNotFound() {
		var error model.Error
		error.Message = "レコードが見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
	return user, err
}

func (ur *UserRepository) Update(userID int) (*model.User, error) {
	var user model.User

	if err := ur.DB.Where("id = ?", userID).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		var error model.Error
		error.Message = "該当するアカウントが見つかりません。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	var bookmarkings []model.Song

	if err := ur.DB.Preload("Bookmarkings").Find(&user).Error; err != nil {
		var error model.Error
		error.Message = "該当する参照が見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	if ur.DB.Model(&user).Related(&bookmarkings, "Bookmarikings").RecordNotFound() {
		error := model.Error{}
		error.Message = "レコードが見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	var followings []model.User
	if err := ur.DB.Preload("Followings").Find(&user).Error; err != nil {
		var error model.Error
		error.Message = "該当する参照が見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
	if ur.DB.Model(&user).Related(&followings, "Followings").RecordNotFound() {
		var error model.Error
		error.Message = "レコードが見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
	return user, err
}
