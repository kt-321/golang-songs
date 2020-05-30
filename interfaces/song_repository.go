package interfaces

import (
	"golang-songs/domain"
)

type SongRepository struct {
	SQLHandler SQLHandler
}

func (pr *SongRepository) FindAll() (posts domain.Posts, err error) {
	const query = `
	SELECT
		id,
		user_id,
		body
	FROM
		posts
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
		var userID uint
		var body string
		if err = rows.Scan(&id, &userID, &body); err != nil {
			return
		}
		song := domain.Song{
			//ID:     id,
			//UserID: userID,
			//Body:   body,
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
		posts = append(songs, song)
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}

// Save is saves the given entity.
func (pr *SongRepository) Save(p domain.Post) (id int64, err error) {
	// NOTE: this is a transaction example.
	tx, err := pr.SQLHandler.Begin()
	if err != nil {
		return
	}

	const query = `
		INSERT INTO
			posts(user_id, body)
		VALUES
			(?, ?)
	`

	result, err := tx.Exec(query, p.UserID, p.Body)
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
func (pr *SongRepository) DeleteByID(postID int) (err error) {
	const query = `
		DELETE
		FROM
			posts
		WHERE
			id = ?
	`

	_, err = pr.SQLHandler.Exec(query, postID)
	if err != nil {
		return
	}

	return
}
