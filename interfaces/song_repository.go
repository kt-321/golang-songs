package interfaces

import (
	"fmt"
	"golang-songs/model"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

type SongRepository struct {
	DB    *gorm.DB
	Redis redis.Conn
}

func (sr *SongRepository) FindAll() (*[]model.Song, error) {
	var songs []model.Song

	if err := sr.DB.Find(&songs).Error; gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &songs, nil
}

func (sr *SongRepository) FindByID(songID int) (*model.Song, error) {
	var song model.Song

	exists, err := redis.Int(sr.Redis.Do("EXISTS", fmt.Sprintf("song:%d", songID)))
	if err != nil {
		return nil, err
	}
	//キャッシュが存在する場合
	if exists > 0 {
		//log.Print("キャッシュが存在する")
		t, err := redis.StringMap(sr.Redis.Do("HGETALL", fmt.Sprintf("song:%d", songID)))
		if err != nil {
			return nil, err
		}

		//ロケーションを指定して、パース
		jst, err := time.LoadLocation("Asia/Tokyo")
		if err != nil {
			return nil, err
		}
		CreatedAt, err := time.ParseInLocation("2006年01月02日 15時04分05秒", t["CreatedAt"], jst)
		if err != nil {
			return nil, err
		}
		UpdatedAt, err := time.ParseInLocation("2006年01月02日 15時04分05秒", t["UpdatedAt"], jst)
		if err != nil {
			return nil, err
		}

		//IDとUserIDをstringからunitに変換
		intID, err := strconv.Atoi(t["ID"])
		if err != nil {
			return nil, err
		}
		uintID := uint(intID)
		if err != nil {
			return nil, err
		}
		intUserId, err := strconv.Atoi(t["UserID"])

		uintUserID := uint(intUserId)

		//MusicAgeをstringからintに変換
		MusicAge, err := strconv.Atoi(t["MusicAge"])
		if err != nil {
			return nil, err
		}

		song = model.Song{
			ID:             uintID,
			Title:          t["Title"],
			Artist:         t["Artist"],
			MusicAge:       MusicAge,
			Image:          t["Image"],
			Video:          t["Video"],
			Album:          t["Album"],
			Description:    t["Description"],
			SpotifyTrackId: t["SpotifyTrackId"],
			UserID:         uintUserID,
			CreatedAt:      CreatedAt,
			UpdatedAt:      UpdatedAt,
			DeletedAt:      nil,
		}

		//キャッシュのTTLを1800秒(30分)にリセット
		_, err = sr.Redis.Do("EXPIRE", fmt.Sprintf("song:%d", songID), "1800")
		if err != nil {
			return nil, err
		}

		return &song, nil
	} else {
		//キャッシュが存在しない場合、DBに取りに行き、取得した値はRedisに保存する
		result := sr.DB.Where("id = ?", songID).Find(&song)

		if result.Error == nil {
			//Redisに入れる前にformatする
			formattedCreatedAt := song.CreatedAt.Format("2006年01月02日 15時04分05秒")
			formattedUpdatedAt := song.UpdatedAt.Format("2006年01月02日 15時04分05秒")

			//Redisに保存
			_, err := redis.String(sr.Redis.Do("HMSET", fmt.Sprintf("song:%d", songID), "ID", songID, "CreatedAt", formattedCreatedAt, "UpdatedAt", formattedUpdatedAt, "DeletedAt", nil, "Title", song.Title, "Artist", song.Artist, "MusicAge", song.MusicAge, "Image", song.Image, "Video", song.Video, "Album", song.Album, "Description", song.Description, "SpotifyTrackId", song.SpotifyTrackId, "UserID", song.UserID))
			if err != nil {
				return nil, err
			}

			//キャッシュのTTLを1800秒(30分)に設定
			_, err = sr.Redis.Do("EXPIRE", fmt.Sprintf("song:%d", songID), "1800")
			if err != nil {
				return nil, err
			}

			return &song, nil
		} else {
			return nil, result.Error
		}
	}
}

func (sr *SongRepository) Save(userEmail string, p model.Song) error {

	var user model.User
	if err := sr.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}

	result := sr.DB.Create(&model.Song{
		Title:          p.Title,
		Artist:         p.Artist,
		MusicAge:       p.MusicAge,
		Image:          p.Image,
		Video:          p.Video,
		Album:          p.Album,
		Description:    p.Description,
		SpotifyTrackId: p.SpotifyTrackId,
		UserID:         user.ID})

	if result.Error == nil {
		var song model.Song

		scanResult := result.Scan(&song)
		if scanResult.Error != nil {
			return scanResult.Error
		}

		//Redisに入れる前にformatする
		formattedCreatedAt := song.CreatedAt.Format("2006年01月02日 15時04分05秒")
		formattedUpdatedAt := song.UpdatedAt.Format("2006年01月02日 15時04分05秒")

		//Redisに入れる
		_, err := redis.String(sr.Redis.Do("HMSET", fmt.Sprintf("song:%d", song.ID), "ID", song.ID, "CreatedAt", formattedCreatedAt, "UpdatedAt", formattedUpdatedAt, "DeletedAt", nil, "Title", song.Title, "Artist", song.Artist, "MusicAge", song.MusicAge, "Image", song.Image, "Video", song.Video, "Album", song.Album, "Description", song.Description, "SpotifyTrackId", song.SpotifyTrackId, "UserID", song.UserID))
		if err != nil {
			return err
		}

		//キャッシュのTTLを1800秒(30分)に設定
		_, err = sr.Redis.Do("EXPIRE", fmt.Sprintf("song:%d", song.ID), "1800")
		if err != nil {
			return err
		}

		return nil
	} else {
		//RDBにInsertするのに失敗した場合
		return result.Error
	}
}

func (sr *SongRepository) UpdateByID(userEmail string, songID int, p model.Song) error {
	var user model.User
	if err := sr.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}

	var song model.Song

	result := sr.DB.Model(&song).Where("id = ?", songID).Update(model.Song{
		Title:          p.Title,
		Artist:         p.Artist,
		MusicAge:       p.MusicAge,
		Image:          p.Image,
		Video:          p.Video,
		Album:          p.Album,
		Description:    p.Description,
		SpotifyTrackId: p.SpotifyTrackId})

	if result.Error == nil {
		var song model.Song

		scanResult := result.Scan(&song)
		if scanResult.Error != nil {
			return scanResult.Error
		}

		//Redisに入れる前にformatする
		formattedCreatedAt := song.CreatedAt.Format("2006年01月02日 15時04分05秒")
		formattedUpdatedAt := song.UpdatedAt.Format("2006年01月02日 15時04分05秒")

		//Redisに保存
		_, err := redis.String(sr.Redis.Do("HMSET", fmt.Sprintf("song:%d", song.ID), "ID", song.ID, "CreatedAt", formattedCreatedAt, "UpdatedAt", formattedUpdatedAt, "DeletedAt", nil, "Title", song.Title, "Artist", song.Artist, "MusicAge", song.MusicAge, "Image", song.Image, "Video", song.Video, "Album", song.Album, "Description", song.Description, "SpotifyTrackId", song.SpotifyTrackId, "UserID", song.UserID))
		if err != nil {
			return err
		}

		//キャッシュのTTLを1800秒(30分)に設定
		_, err = sr.Redis.Do("EXPIRE", fmt.Sprintf("song:%d", songID), "1800")
		if err != nil {
			return err
		}

		return nil
	} else {
		//RDBにInsertするのに失敗した場合
		return result.Error
	}
}

func (sr *SongRepository) DeleteByID(songID int) error {
	var song model.Song

	if err := sr.DB.Where("id = ?", songID).Delete(&song).Error; err != nil {
		return err
	}

	//Redisに入っている場合のみに削除
	exists, err := redis.Int(sr.Redis.Do("EXISTS", fmt.Sprintf("song:%d", songID)))
	if err != nil {
		return err
	}
	if exists > 0 {
		_, err := sr.Redis.Do("DEL", fmt.Sprintf("song:%d", songID))
		if err != nil {
			return err
		}
	}

	return nil
}
