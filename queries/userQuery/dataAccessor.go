package userQuery

import (
	"golang-songs/model"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
)

type dataAccessor struct {
	DB *sqlx.DB
}

func (ur *dataAccessor) FindAll() (*[]model.User, error) {
	var users []model.User
	if err := ur.DB.Find(&users).Error; gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	return &users, nil
}

func (ur *dataAccessor) GetUser(userEmail string) (*model.User, error) {
	var user model.User
	if err := ur.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	if err := ur.DB.Preload("Bookmarkings", "bookmarks.deleted_at is null").Preload("Followings", "user_follows.deleted_at is null").Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *dataAccessor) FindByID(userID int) (*model.User, error) {
	var user model.User

	if err := ur.DB.Where("id = ?", userID).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	if err := ur.DB.Preload("Bookmarkings", "bookmarks.deleted_at is null").Preload("Followings", "user_follows.deleted_at is null").Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
