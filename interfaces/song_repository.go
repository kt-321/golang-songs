package interfaces

import (
	"fmt"
	"golang-songs/model"
	"log"
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

//キャッシュが存在するか確認
func ExistsSongByID(songID int, rc redis.Conn) (int, error) {
	exists, err := redis.Int(rc.Do("EXISTS", fmt.Sprintf("song:%d", songID)))
	if err != nil {
		return 0, err
	}
	return exists, nil
}

//func GetSongByID(songID int, rc redis.Conn) (*model.Song, error) {
func GetSongByID(songID int, rc redis.Conn) (map[string]string, error) {
	t, err := redis.StringMap(rc.Do("HGETALL", fmt.Sprintf("song:%d", songID)))
	if err != nil {
		log.Println("キャッシュ取得失敗GetSongByID")
		return nil, err
	}
	log.Println("キャッシュ取得成功GetSongByID")

	return t, nil
}

func SetSongByID(songID int, t map[string]string, rc redis.Conn) error {
	log.Println("SetSongByID")
	log.Println(t, &rc)
	_, err := rc.Do("HMSET", fmt.Sprintf("song:%d", songID), "ID", t["ID"], "CreatedAt", t["CreatedAt"], "UpdatedAt", t["UpdatedAt"], "DeletedAt", t["DeletedAt"], "Title", t["Title"], "Artist", t["Artist"], "MusicAge", t["MusicAge"], "Image", t["Image"], "Video", t["Video"], "Album", t["Album"], "Description", t["Description"], "SpotifyTrackId", t["SpotifyTrackId"], "UserID", t["UserID"])
	if err != nil {
		log.Println("キャッシュ保存失敗SetSongByID")
		log.Print(err)
		return err
	}
	log.Println("キャッシュ保存成功SetSongByID")

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
	//exists, err := redis.Int(sr.Redis.Do("EXISTS", fmt.Sprintf("song:%d", songID)))
	if err != nil {
		return nil, err
	}
	// サイドカーコンテナのRedisにキャッシュが存在する場合
	if exists > 0 {
		log.Println("キャッシュあり")

		// サイドカーのRedisのキャッシュを取得
		t, err := GetSongByID(songID, sr.SidecarRedis)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		log.Println("キャッシュ取得成功2")

		//return song, nil

		//予めtime.Localにタイムゾーンの設定情報を入れておく
		time.Local = time.FixedZone("Local", 9*60*60)
		//ロケーションを指定して、パース
		jst, err := time.LoadLocation("Local")
		if err != nil {
			log.Println(err)
			return nil, err
		}

		CreatedAt, err := time.ParseInLocation("2006年01月02日 15時04分05秒", t["CreatedAt"], jst)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		UpdatedAt, err := time.ParseInLocation("2006年01月02日 15時04分05秒", t["UpdatedAt"], jst)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		//IDとUserIDをstringからunitに変換
		intID, err := strconv.Atoi(t["ID"])
		if err != nil {
			log.Println(err)
			return nil, err
		}
		uintID := uint(intID)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		log.Printf("%v:%T", t, t)
		log.Printf("%v:%T", t["UserID"], t["UserID"])
		intUserId, err := strconv.Atoi(t["UserID"])
		log.Printf("%v:%T", intUserId, intUserId)
		uintUserID := uint(intUserId)
		log.Printf("%v:%T", uintUserID, uintUserID)

		//MusicAgeをstringからintに変換
		MusicAge, err := strconv.Atoi(t["MusicAge"])
		if err != nil {
			log.Println(err)
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
			log.Println("リモートにキャッシュあり")
			// リモートのRedisのキャッシュを取得
			t, err := GetSongByID(songID, sr.Redis)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			//予めtime.Localにタイムゾーンの設定情報を入れておく
			time.Local = time.FixedZone("Local", 9*60*60)
			//ロケーションを指定して、パース
			jst, err := time.LoadLocation("Local")
			if err != nil {
				log.Println(err)
				return nil, err
			}

			CreatedAt, err := time.ParseInLocation("2006年01月02日 15時04分05秒", t["CreatedAt"], jst)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			UpdatedAt, err := time.ParseInLocation("2006年01月02日 15時04分05秒", t["UpdatedAt"], jst)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			//IDとUserIDをstringからunitに変換
			intID, err := strconv.Atoi(t["ID"])
			if err != nil {
				log.Println(err)
				return nil, err
			}
			uintID := uint(intID)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			intUserId, err := strconv.Atoi(t["UserID"])

			uintUserID := uint(intUserId)

			//MusicAgeをstringからintに変換
			MusicAge, err := strconv.Atoi(t["MusicAge"])
			if err != nil {
				log.Println(err)
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
			//_, err := redis.String(sr.SidecarRedis.Do("HMSET", fmt.Sprintf("song:%d", songID), "ID", songID, "CreatedAt", formattedCreatedAt, "UpdatedAt", formattedUpdatedAt, "DeletedAt", nil, "Title", song.Title, "Artist", song.Artist, "MusicAge", song.MusicAge, "Image", song.Image, "Video", song.Video, "Album", song.Album, "Description", song.Description, "SpotifyTrackId", song.SpotifyTrackId, "UserID", song.UserID))
			if err != nil {
				return nil, err
			}

			//キャッシュのTTLを1800秒(30分)に設定
			_, err = sr.SidecarRedis.Do("EXPIRE", fmt.Sprintf("song:%d", songID), "1800")
			if err != nil {
				return nil, err
			}

			log.Println("リモートから取ってきてサイドカーに保存完了")

			//return song, nil
		} else {
			log.Println("リモートにもサイドカーにも値なし")
			//リモートのRedisにキャッシュが存在しない場合RDSに値を取りに行く。
			result := sr.DB.Where("id = ?", songID).Find(&song)

			//RDSからの値取得に成功した場合
			if result.Error == nil {
				//Redisに保存する前にformatする
				formattedCreatedAt := song.CreatedAt.Format("2006年01月02日 15時04分05秒")
				formattedUpdatedAt := song.UpdatedAt.Format("2006年01月02日 15時04分05秒")
				log.Printf("%v:%T", songID, songID)
				log.Printf("%v:%T", song.MusicAge, song.MusicAge)
				t := map[string]string{
					"ID":             strconv.Itoa(songID),
					"CreatedAt":      formattedCreatedAt,
					"UpdatedAt":      formattedUpdatedAt,
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
				log.Printf("%v", t)

				//リモートのRedisに保存
				err = SetSongByID(songID, t, sr.Redis)
				//_, err := redis.String(sr.Redis.Do("HMSET", fmt.Sprintf("song:%d", songID), "ID", songID, "CreatedAt", formattedCreatedAt, "UpdatedAt", formattedUpdatedAt, "DeletedAt", nil, "Title", song.Title, "Artist", song.Artist, "MusicAge", song.MusicAge, "Image", song.Image, "Video", song.Video, "Album", song.Album, "Description", song.Description, "SpotifyTrackId", song.SpotifyTrackId, "UserID", song.UserID))
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
				//_, err = redis.String(sr.SidecarRedis.Do("HMSET", fmt.Sprintf("song:%d", songID), t))
				//_, err = redis.String(sr.SidecarRedis.Do("HMSET", fmt.Sprintf("song:%d", songID), "ID", songID, "CreatedAt", formattedCreatedAt, "UpdatedAt", formattedUpdatedAt, "DeletedAt", nil, "Title", song.Title, "Artist", song.Artist, "MusicAge", song.MusicAge, "Image", song.Image, "Video", song.Video, "Album", song.Album, "Description", song.Description, "SpotifyTrackId", song.SpotifyTrackId, "UserID", song.UserID))
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
				log.Println("RDSからの値取得に失敗")
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
