package interfaces

import (
	"errors"
	"fmt"
	"golang-songs/model"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"golang.org/x/sync/singleflight"
)

type SongRepository struct {
	DB           *gorm.DB
	Redis        redis.Conn
	SidecarRedis redis.Conn
}

// Redisにキャッシュが存在するか確認.
func ExistsSongByID(songID int, rc redis.Conn) (int, error) {
	exists, err := redis.Int(rc.Do("EXISTS", fmt.Sprintf("song:%d", songID)))

	if err != nil {
		return 0, fmt.Errorf("failed to check existing cache in Redis: %v", err)
	}

	return exists, nil
}

// Redisから該当する曲を取得.
func GetSongByID(songID int, rc redis.Conn) (map[string]string, error) {
	t, err := redis.StringMap(rc.Do("HGETALL", fmt.Sprintf("song:%d", songID)))

	if err != nil {
		return nil, err
	}

	return t, nil
}

// Redisに曲を保存.
func SetSongByID(songID int, t map[string]string, rc redis.Conn, ttl int) error {
	_, err := rc.Do("HMSET", fmt.Sprintf("song:%d", songID), "ID", t["ID"], "CreatedAt", t["CreatedAt"], "UpdatedAt", t["UpdatedAt"], "DeletedAt", t["DeletedAt"], "Title", t["Title"], "Artist", t["Artist"], "MusicAge", t["MusicAge"], "Image", t["Image"], "Video", t["Video"], "Album", t["Album"], "Description", t["Description"], "SpotifyTrackId", t["SpotifyTrackId"], "UserID", t["UserID"])

	if err != nil {
		return fmt.Errorf("failed to save a song to Redis: %v", err)
	}

	if ttl > 0 {
		// キャッシュのTTLを設定
		_, err = rc.Do("EXPIRE", fmt.Sprintf("song:%d", songID), ttl)

		if err != nil {
			return fmt.Errorf("failed to set ttl in SetSongByID: %v", err)
		}
	}

	return nil
}

// Redisから該当する曲を削除.
func DeleteSongByID(songID int, rc redis.Conn) error {
	_, err := rc.Do("DEL", fmt.Sprintf("song:%d", songID))
	if err != nil {
		return fmt.Errorf("failed to delete song by ID: %v", err)
	}

	return nil
}

// mapから構造体Songへと変換.
func MapToSong(t map[string]string) (*model.Song, error) {
	// ロケーションを指定して、パース.
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

	// IDとUserIDをstringからunitに変換.
	intID, err := strconv.Atoi(t["ID"])
	if err != nil {
		return nil, err
	}

	uintID := uint(intID)
	intUserID, err := strconv.Atoi(t["UserID"])

	if err != nil {
		return nil, err
	}

	// MusicAgeをstringからintに変換.
	MusicAge, err := strconv.Atoi(t["MusicAge"])
	if err != nil {
		return nil, err
	}

	song := model.Song{
		ID:             uintID,
		Title:          t["Title"],
		Artist:         t["Artist"],
		MusicAge:       MusicAge,
		Image:          t["Image"],
		Video:          t["Video"],
		Album:          t["Album"],
		Description:    t["Description"],
		SpotifyTrackId: t["SpotifyTrackId"],
		UserID:         uint(intUserID),
		CreatedAt:      CreatedAt,
		UpdatedAt:      UpdatedAt,
		DeletedAt:      nil,
	}

	return &song, nil
}

// RDSから取得した曲情報をmapへと変換.
func rdsSongToMap(song model.Song) map[string]string {
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

	return t
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

	// singleflightで同時関数呼び出しを1度に抑える.
	var g singleflight.Group
	v, err, _ := g.Do("key", func() (interface{}, error) {
		// サイドカーコンテナのRedisにキャッシュがあるか確認
		exists, err := ExistsSongByID(songID, sr.SidecarRedis)

		if err != nil {
			return nil, err
		}

		// サイドカーコンテナのRedisにキャッシュが存在する場合.
		if exists > 0 {
			// サイドカーのRedisのキャッシュを取得.
			t, err := GetSongByID(songID, sr.SidecarRedis)

			if err != nil {
				return nil, err
			}

			// mapから構造体Songへと変換.
			value, err := MapToSong(t)

			if err != nil {
				return nil, err
			}
			song = *value
		} else {
			// サイドカーコンテナのRedisにキャッシュが存在しない場合、リモートのRedisにキャッシュがあるか確認.
			exists, err := ExistsSongByID(songID, sr.Redis)

			if err != nil {
				return nil, err
			}

			// リモートのRedisにキャッシュが存在する場合.
			if exists > 0 {
				// リモートのRedisのキャッシュを取得.
				t, err := GetSongByID(songID, sr.Redis)

				if err != nil {
					return nil, err
				}

				// mapから構造体Songへと変換.
				value, err := MapToSong(t)

				if err != nil {
					return nil, err
				}

				song = *value

				// サイドカーコンテナのRedisに保存。キャッシュのTTLは1800秒(30分).
				err = SetSongByID(songID, t, sr.SidecarRedis, 1800)

				if err != nil {
					return nil, err
				}
			} else {
				// リモートのRedisにキャッシュが存在しない場合RDSに値を取りに行く.
				result := sr.DB.Where("id = ?", songID).Find(&song)

				// RDSからの値取得に成功した場合
				if result.Error == nil {
					// RDSから取得した曲情報をmapへと変換.
					t := rdsSongToMap(song)

					// リモートのRedisに保存。キャッシュのTTLは1800秒(30分).
					err = SetSongByID(songID, t, sr.Redis, 1800)

					if err != nil {
						return nil, err
					}

					// サイドカーコンテナのRedisに保存。キャッシュのTTLは1800秒(30分).
					err = SetSongByID(songID, t, sr.SidecarRedis, 1800)

					if err != nil {
						return nil, err
					}
				} else {
					// RDSからの値取得に失敗した場合.
					return nil, result.Error
				}
			}
		}

		return song, nil
	})

	if err != nil {
		return nil, err
	}

	// model.Song型に戻す.
	responseSong, ok := v.(model.Song)

	if !ok {
		return nil, errors.New("型変換に失敗")
	}

	return &responseSong, nil
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

	// RDSへのInsertに成功した場合.
	if result.Error == nil {
		var song model.Song

		scanResult := result.Scan(&song)

		if scanResult.Error != nil {
			return scanResult.Error
		}

		// 曲情報をmapへと変換.
		t := rdsSongToMap(song)

		// リモートのRedisに入れる。キャッシュのTTLは1800秒(30分).
		err := SetSongByID(int(song.ID), t, sr.Redis, 1800)

		if err != nil {
			return err
		}

		// サイドカーコンテナのRedisに入れる。キャッシュのTTLは1800秒(30分).
		err = SetSongByID(int(song.ID), t, sr.SidecarRedis, 1800)

		if err != nil {
			return err
		}

		return nil
	}
	// RDBへのInsertに失敗した場合.
	return result.Error
}

func (sr *SongRepository) UpdateByID(userEmail string, songID int, p model.Song) error {
	var user model.User
	// RDSからリクエストユーザーの情報を取得.
	if err := sr.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}

	var song model.Song

	// RDSの該当する曲のデータを更新.
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

		// 曲情報をmapへと変換.
		t := rdsSongToMap(updatedSong)

		// リモートのRedisに保存。キャッシュのTTLは1800秒(30分).
		err := SetSongByID(int(updatedSong.ID), t, sr.Redis, 1800)

		if err != nil {
			return err
		}

		// サイドカーコンテナのRedisに保存。キャッシュのTTLは1800秒(30分).
		err = SetSongByID(int(updatedSong.ID), t, sr.SidecarRedis, 1800)
		if err != nil {
			return err
		}

		return nil
	}
	// RDBにInsertするのに失敗した場合.
	return result.Error
}

func (sr *SongRepository) DeleteByID(songID int) error {
	var song model.Song

	if err := sr.DB.Where("id = ?", songID).Delete(&song).Error; err != nil {
		return err
	}

	// リモートのRedisに入っている場合のみに削除.
	remoteExists, err := redis.Int(sr.Redis.Do("EXISTS", fmt.Sprintf("song:%d", songID)))

	if err != nil {
		return err
	}

	if remoteExists > 0 {
		err := DeleteSongByID(songID, sr.Redis)
		if err != nil {
			return err
		}
	}

	// サイドカーコンテナのRedisに入っている場合のみに削除.
	sidecarExists, err := redis.Int(sr.SidecarRedis.Do("EXISTS", fmt.Sprintf("song:%d", songID)))

	if err != nil {
		return err
	}

	if sidecarExists > 0 {
		err := DeleteSongByID(songID, sr.SidecarRedis)
		if err != nil {
			return err
		}
	}

	return nil
}
