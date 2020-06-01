package interfaces

import (
	"golang-songs/domain"

	"github.com/jinzhu/gorm"
)

type SongRepository struct {
	//SQLHandler SQLHandler
	DB *gorm.DB
}

func (pr *SongRepository) FindAll() (songs domain.Songs, err error) {
	const query = `
	SELECT
		id,
		user_id,
		body
	FROM
		songs
	`

	rows, err := pr.SQLHandler.Query(query)

	defer rows.Close()

	if err != nil {
		return
	}

	for rows.Next() {
		//var id int
		var id uint
		//var userID int
		var title string
		var artist string
		var musicAge int
		var image string
		var video string
		var album string
		var description string
		var spotifyTrackId string
		var userID uint
		if err = rows.Scan(&id, &title, &artist, &musicAge, &image, &video, &album, &description, &spotifyTrackId, &userID); err != nil {
			return
		}
		song := domain.Song{
			ID:             id,
			Title:          title,
			Artist:         artist,
			MusicAge:       musicAge,
			Image:          image,
			Video:          video,
			Album:          album,
			Description:    description,
			SpotifyTrackId: spotifyTrackId,
			UserID:         userID,
		}
		//
		//if err := f.DB.Create(&model.Song{
		//	Title:          d.Title,
		//	Artist:         d.Artist,
		//	MusicAge:       d.MusicAge,
		//	Image:          d.Image,
		//	Video:          d.Video,
		//	Album:          d.Album,
		//	Description:    d.Description,
		//	SpotifyTrackId: d.SpotifyTrackId,
		//	UserID:         user.ID}).Error; err != nil {
		//	var error model.Error
		//	error.Message = "曲の追加に失敗しました"
		//	errorInResponse(w, http.StatusInternalServerError, error)
		//	return
		//}
		songs = append(songs, song)
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}

// Save is saves the given entity.
func (pr *SongRepository) Save(p domain.Song) (id int64, err error) {
	// NOTE: this is a transaction example.
	tx, err := pr.SQLHandler.Begin()
	if err != nil {
		return
	}

	const query = `
		INSERT INTO
			songs(user_id, body)
		VALUES
			(?, ?)
	`

	result, err := tx.Exec(query, p.Title, p.Artist, p.MusicAge, p.Image, p.Video, p.Album, p.Description, p.SpotifyTrackId, p.UserID)
	if err != nil {
		_ = tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	id, err = result.LastInsertId()
	if err != nil {
		return id, nil
	}

	return
}

// DeleteByID is deletes the entity identified by the given id.
func (pr *SongRepository) DeleteByID(songID int) (err error) {
	const query = `
		DELETE
		FROM
			songs
		WHERE
			id = ?
	`

	_, err = pr.SQLHandler.Exec(query, songID)
	if err != nil {
		return
	}

	return
}
