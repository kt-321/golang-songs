package interfaces

import (
	"github.com/jmoiron/sqlx"
)

type BookmarkRepository struct {
	DB *sqlx.DB
}

// 曲をお気に入り登録.
func (br *BookmarkRepository) Bookmark(userEmail string, songID int) error {
	//// リクエストユーザーを取得.
	//var user model.User
	//if err := br.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
	//	return err
	//}
	//
	//// お気に入り登録する対象の曲を取得.
	//var song model.Song
	//if err := br.DB.Where("id = ?", songID).Find(&song).Error; gorm.IsRecordNotFoundError(err) {
	//	return err
	//}
	//
	//// deleted_atがnullであるbookmarksレコードがある時はレコード追加せず、該当レコードがない時はレコード追加する.
	//if err := br.DB.
	//	Where("user_id = ?", user.ID).
	//	Where("song_id = ?", song.ID).
	//	FirstOrCreate(&model.Bookmark{
	//		UserID: user.ID,
	//		SongID: song.ID}).
	//	Error; err != nil {
	//	return err
	//}

	//TODO
	return nil
}

// 曲をお気に入り登録から解除.
func (br *BookmarkRepository) RemoveBookmark(userEmail string, songID int) error {
	//var bookmark model.Bookmark
	//
	//// bookmarksテーブルをusersテーブルやsongsテーブルと内部結合して、該当するレコードを取得する.
	//if err := br.DB.Table("bookmarks").
	//	Where("bookmarks.deleted_at is null").
	//	Joins("INNER JOIN users ON users.id = bookmarks.user_id AND users.email = ? AND users.deleted_at is null", userEmail).
	//	Joins("INNER JOIN songs ON songs.id = bookmarks.song_id AND songs.id = ? AND songs.deleted_at is null", songID).
	//	Scan(&bookmark).Error; gorm.IsRecordNotFoundError(err) {
	//	return err
	//}
	//
	//// bookmarksレコードを論理削除.
	//if err := br.DB.Delete(&bookmark).Error; err != nil {
	//	return err
	//}

	//TODO
	return nil
}
