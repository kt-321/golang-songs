package interfaces

import (
	"golang-songs/model"

	"github.com/jinzhu/gorm"
)

type UserFollowRepository struct {
	DB *gorm.DB
}

func (ufr *UserFollowRepository) Follow(requestUserEmail string, targetUserID int) error {
	var requestUser model.User

	if err := ufr.DB.Where("email = ?", requestUserEmail).Find(&requestUser).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}

	var targetUser model.User

	if err := ufr.DB.Where("id = ?", targetUserID).Find(&targetUser).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}

	if err := ufr.DB.Create(&model.UserFollow{
		UserID:   requestUser.ID,
		FollowID: targetUser.ID}).Error; err != nil {
		return err
	}

	if err := ufr.DB.Preload("Followings").Find(&requestUser).Error; err != nil {
		return err
	}

	if err := ufr.DB.Model(&requestUser).Association("Followings").Append(&targetUser).Error; err != nil {
		return err
	}

	return nil
}

func (ufr *UserFollowRepository) Unfollow(requestUserEmail string, targetUserID int) error {
	var requestUser model.User

	if err := ufr.DB.Where("email = ?", requestUserEmail).Find(&requestUser).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}

	var targetUser model.User

	if err := ufr.DB.Where("id = ?", targetUserID).Find(&targetUser).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}

	if err := ufr.DB.Preload("Followings").Find(&requestUser).Error; err != nil {
		return err
	}

	if err := ufr.DB.Model(&requestUser).Association("Followings").Delete(&targetUser).Error; err != nil {
		return err
	}

	return nil
}
