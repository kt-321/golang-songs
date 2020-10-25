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
	DB           *gorm.DB
	Redis        redis.Conn
	SidecarRedis redis.Conn
}

//Redisにキャッシュが存在するか確認
func ExistsSongByID(songID int, rc redis.Conn) (int, error) {
	exists, err := redis.Int(rc.Do("EXISTS", fmt.Sprintf("song:%d", songID)))
	if err != nil {
		return 0, err
	}
	return exists, nil
}

//Redisから該当する曲を取得
func GetSongByID(songID int, rc redis.Conn) (map[string]string, error) {
	t, err := redis.StringMap(rc.Do("HGETALL", fmt.Sprintf("song:%d", songID)))
	if err != nil {
		return nil, err
	}
	return t, nil
}

//Redisに曲を保存
func SetSongByID(songID int, t map[string]string, rc redis.Conn) error {
	_, err := rc.Do("HMSET", fmt.Sprintf("song:%d", songID), "ID", t["ID"], "CreatedAt", t["CreatedAt"], "UpdatedAt", t["UpdatedAt"], "DeletedAt", t["DeletedAt"], "Title", t["Title"], "Artist", t["Artist"], "MusicAge", t["MusicAge"], "Image", t["Image"], "Video", t["Video"], "Album", t["Album"], "Description", t["Description"], "SpotifyTrackId", t["SpotifyTrackId"], "UserID", t["UserID"])
	if err != nil {
		return err
	}
	return nil
}

//Redisから該当する曲を削除
func DeleteSongByID(songID int, rc redis.Conn) error {
	_, err := rc.Do("DEL", fmt.Sprintf("song:%d", songID))
	if err != nil {
		return err
	}
	return nil
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

	// サイドカーコンテナのRedisにキャッシュがあるか確認
	exists, err := ExistsSongByID(songID, sr.SidecarRedis)
	if err != nil {
		return nil, err
	}
	// サイドカーコンテナのRedisにキャッシュが存在する場合
	if exists > 0 {
		// サイドカーのRedisのキャッシュを取得
		t, err := GetSongByID(songID, sr.SidecarRedis)
		if err != nil {
			return nil, err
		}

		//予めtime.Localにタイムゾーンの設定情報を入れておく
		time.Local = time.FixedZone("Local", 9*60*60)
		//ロケーションを指定して、パース
		jst, err := time.LoadLocation("Local")
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
	} else {
		// サイドカーコンテナのRedisにキャッシュが存在しない場合、リモートのRedisにキャッシュがあるか確認
		exists, err := ExistsSongByID(songID, sr.Redis)
		if err != nil {
			return nil, err
		}

		//リモートのRedisにキャッシュが存在する場合
		if exists > 0 {
			// リモートのRedisのキャッシュを取得
			t, err := GetSongByID(songID, sr.Redis)
			if err != nil {
				return nil, err
			}

			//予めtime.Localにタイムゾーンの設定情報を入れておく
			time.Local = time.FixedZone("Local", 9*60*60)
			//ロケーションを指定して、パース
			jst, err := time.LoadLocation("Local")
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

			//サイドカーコンテナのRedisに保存
			err = SetSongByID(songID, t, sr.SidecarRedis)
			if err != nil {
				return nil, err
			}

			//キャッシュのTTLを1800秒(30分)に設定
			_, err = sr.SidecarRedis.Do("EXPIRE", fmt.Sprintf("song:%d", songID), "1800")
			if err != nil {
				return nil, err
			}
		} else {
			//リモートのRedisにキャッシュが存在しない場合RDSに値を取りに行く。
			result := sr.DB.Where("id = ?", songID).Find(&song)

			//RDSからの値取得に成功した場合
			if result.Error == nil {
				t := map[string]string{
					"ID":             strconv.Itoa(songID),
					"CreatedAt":      song.CreatedAt.Format("2006年01月02日 15時04分05秒"),
					"UpdatedAt":      song.UpdatedAt.Format("2006年01月02日 15時04分05秒"),
					"DeletedAt":      "",
					"Title":          song.Title,
					"Artist":         song.Artist,
					"MusicAge":       strconv.Itoa(song.MusicAge),
					"Image":          song.Image,
					"Video":          song.Video,
					"Album":          song.Album,
					"Description":    song.Description,
					"SpotifyTrackId": song.SpotifyTrackId,
					"UserID":         strconv.Itoa(int(song.UserID)),
				}

				//リモートのRedisに保存
				err = SetSongByID(songID, t, sr.Redis)
				if err != nil {
					return nil, err
				}
				//キャッシュのTTLを1800秒(30分)に設定
				_, err = sr.Redis.Do("EXPIRE", fmt.Sprintf("song:%d", songID), "1800")
				if err != nil {
					return nil, err
				}

				//サイドカーコンテナのRedisに保存
				err = SetSongByID(songID, t, sr.SidecarRedis)
				if err != nil {
					return nil, err
				}
				//キャッシュのTTLを1800秒(30分)に設定
				_, err = sr.SidecarRedis.Do("EXPIRE", fmt.Sprintf("song:%d", songID), "1800")
				if err != nil {
					return nil, err
				}
			} else {
				//RDSからの値取得に失敗した場合
				return nil, result.Error
			}

		}

	}
	return &song, nil
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

	//RDSへのInsertに成功した場合
	if result.Error == nil {
		var song model.Song

		scanResult := result.Scan(&song)
		if scanResult.Error != nil {
			return scanResult.Error
		}

		//mapに変換
		t := map[string]string{
			"ID":             strconv.Itoa(int(song.ID)),
			"CreatedAt":      song.CreatedAt.Format("2006年01月02日 15時04分05秒"),
			"UpdatedAt":      song.UpdatedAt.Format("2006年01月02日 15時04分05秒"),
			"DeletedAt":      "",
			"Title":          song.Title,
			"Artist":         song.Artist,
			"MusicAge":       strconv.Itoa(song.MusicAge),
			"Image":          song.Image,
			"Video":          song.Video,
			"Album":          song.Album,
			"Description":    song.Description,
			"SpotifyTrackId": song.SpotifyTrackId,
			"UserID":         strconv.Itoa(int(song.UserID)),
		}

		//リモートのRedisに入れる
		err := SetSongByID(int(song.ID), t, sr.Redis)
		if err != nil {
			return err
		}
		//キャッシュのTTLを1800秒(30分)に設定
		_, err = sr.Redis.Do("EXPIRE", fmt.Sprintf("song:%d", song.ID), "1800")
		if err != nil {
			return err
		}

		//サイドカーコンテナのRedisに入れる
		err = SetSongByID(int(song.ID), t, sr.SidecarRedis)
		if err != nil {
			return err
		}
		//キャッシュのTTLを1800秒(30分)に設定
		_, err = sr.SidecarRedis.Do("EXPIRE", fmt.Sprintf("song:%d", song.ID), "1800")
		if err != nil {
			return err
		}

		return nil
	} else {
		//RDBへのInsertに失敗した場合
		return result.Error
	}
}

func (sr *SongRepository) UpdateByID(userEmail string, songID int, p model.Song) error {
	var user model.User
	//RDSからリクエストユーザーの情報を取得
	if err := sr.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}

	var song model.Song

	//RDSの該当する曲のデータを更新
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
		var updatedSong model.Song

		scanResult := result.Scan(&updatedSong)
		if scanResult.Error != nil {
			return scanResult.Error
		}

		//mapに変換
		t := map[string]string{
			"ID":             strconv.Itoa(int(updatedSong.ID)),
			"CreatedAt":      updatedSong.CreatedAt.Format("2006年01月02日 15時04分05秒"),
			"UpdatedAt":      updatedSong.UpdatedAt.Format("2006年01月02日 15時04分05秒"),
			"DeletedAt":      "",
			"Title":          updatedSong.Title,
			"Artist":         updatedSong.Artist,
			"MusicAge":       strconv.Itoa(updatedSong.MusicAge),
			"Image":          updatedSong.Image,
			"Video":          updatedSong.Video,
			"Album":          updatedSong.Album,
			"Description":    updatedSong.Description,
			"SpotifyTrackId": updatedSong.SpotifyTrackId,
			"UserID":         strconv.Itoa(int(updatedSong.UserID)),
		}

		//リモートのRedisに保存
		err := SetSongByID(int(updatedSong.ID), t, sr.Redis)
		if err != nil {
			return err
		}
		//キャッシュのTTLを1800秒(30分)に設定
		_, err = sr.Redis.Do("EXPIRE", fmt.Sprintf("song:%d", updatedSong.ID), "1800")
		if err != nil {
			return err
		}

		//サイドカーコンテナのRedisに保存
		err = SetSongByID(int(updatedSong.ID), t, sr.SidecarRedis)
		if err != nil {
			return err
		}
		//キャッシュのTTLを1800秒(30分)に設定
		_, err = sr.SidecarRedis.Do("EXPIRE", fmt.Sprintf("song:%d", updatedSong.ID), "1800")
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

	//リモートのRedisに入っている場合のみに削除
	remoteExists, err := redis.Int(sr.Redis.Do("EXISTS", fmt.Sprintf("song:%d", songID)))
	if err != nil {
		return err
	}
	if remoteExists > 0 {
		err := DeleteSongByID(int(songID), sr.Redis)
		if err != nil {
			return err
		}
	}

	//サイドカーコンテナのRedisに入っている場合のみに削除
	sidecarExists, err := redis.Int(sr.SidecarRedis.Do("EXISTS", fmt.Sprintf("song:%d", songID)))
	if err != nil {
		return err
	}
	if sidecarExists > 0 {
		err := DeleteSongByID(int(songID), sr.SidecarRedis)
		if err != nil {
			return err
		}
	}
	return nil
}
