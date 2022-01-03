package interfaces

import (
	"golang-songs/model"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
)

type UserFollowRepository struct {
	DB *sqlx.DB
}

func (ufr *UserFollowRepository) Follow(requestUserEmail string, targetUserID int) error {
	// リクエストユーザーを取得.
	var requestUser model.User
	if err := ufr.DB.Where("email = ?", requestUserEmail).Find(&requestUser).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}

	// フォローする対象のユーザーを取得.
	var targetUser model.User
	if err := ufr.DB.Where("id = ?", targetUserID).Find(&targetUser).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}

	// deleted_atがnullであるuser_followsレコードがある時はレコード追加せず、該当レコードがない時はレコード追加する.
	if err := ufr.DB.
		Where("user_id = ?", requestUser.ID).
		Where("follow_id = ?", targetUser.ID).
		FirstOrCreate(&model.UserFollow{
			UserID:   requestUser.ID,
			FollowID: targetUser.ID}).
		Error; err != nil {
		return err
	}

	return nil
}

func (ufr *UserFollowRepository) Unfollow(requestUserEmail string, targetUserID int) error {
	var userFollow model.UserFollow

	// user_followsテーブルをusersテーブルと内部結合して、該当するレコードを取得する.
	if err := ufr.DB.Debug().Table("user_follows").
		Where("user_follows.deleted_at is null").
		Joins("INNER JOIN users as user1 ON user1.id = user_follows.user_id AND user1.email = ? AND user1.deleted_at is null", requestUserEmail).
		Joins("INNER JOIN users as user2 ON user2.id = user_follows.follow_id AND user2.id = ? AND user2.deleted_at is null", targetUserID).
		Scan(&userFollow).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}

	// user_followsレコードを論理削除.
	if err := ufr.DB.Debug().Delete(&userFollow).Error; err != nil {
		return err
	}

	return nil
}
