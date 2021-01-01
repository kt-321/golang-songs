package interfaces

import (
	"golang-songs/model"

	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (ur *UserRepository) FindAll() (*[]model.User, error) {
	var users []model.User
	if err := ur.DB.Find(&users).Error; gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	return &users, nil
}

func (ur *UserRepository) GetUser(userEmail string) (*model.User, error) {
	var user model.User
	if err := ur.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	if err := ur.DB.Preload("Bookmarkings", "bookmarks.deleted_at is null").Preload("Followings", "user_follows.deleted_at is null").Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) FindByID(userID int) (*model.User, error) {
	var user model.User

	if err := ur.DB.Where("id = ?", userID).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	if err := ur.DB.Preload("Bookmarkings", "bookmarks.deleted_at is null").Preload("Followings", "user_follows.deleted_at is null").Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) Update(userID int, p model.User) error {
	var user model.User

	if err := ur.DB.Model(&user).Where("id = ?", userID).Update(model.User{Email: p.Email, Name: p.Name, Age: p.Age, Gender: p.Gender, FavoriteMusicAge: p.FavoriteMusicAge, FavoriteArtist: p.FavoriteArtist, Comment: p.Comment}).Error; err != nil {
		return err
	}

	return nil
}
