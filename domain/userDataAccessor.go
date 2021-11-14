package domain

import (
	"golang-songs/model"

	"github.com/jinzhu/gorm"
)

type UserDataAccessor struct {
	DB *gorm.DB
}

func (uda *UserDataAccessor) updateUser(userID int, p model.User) error {
	var user model.User
	if err := uda.DB.Model(&user).Where("id = ?", userID).Update(model.User{Email: p.Email, Name: p.Name, Age: p.Age, Gender: p.Gender, FavoriteMusicAge: p.FavoriteMusicAge, FavoriteArtist: p.FavoriteArtist, Comment: p.Comment}).Error; err != nil {
		return err
	}
	return nil
}
