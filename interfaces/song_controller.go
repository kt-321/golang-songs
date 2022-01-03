package interfaces

import (
	"encoding/json"
	"golang-songs/model"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/garyburd/redigo/redis"

	"golang-songs/usecases"
	"net/http"
)

type SongController struct {
	SongInteractor usecases.SongInteractor
}

func NewSongController(DB *sqlx.DB, Redis redis.Conn, SidecarRedis redis.Conn) *SongController {
	return &SongController{
		SongInteractor: usecases.SongInteractor{
			SongRepository: &SongRepository{
				DB:           DB,
				Redis:        Redis,
				SidecarRedis: SidecarRedis,
			},
		},
	}
}

// 全ての曲を返す.
func (sc *SongController) AllSongsHandler(w http.ResponseWriter, r *http.Request) {
	songs, err := sc.SongInteractor.Index()
	if err != nil {
		ErrorInResponse(w, http.StatusInternalServerError, GetSongError)

		return
	}

	v, err := json.Marshal(songs)

	if err != nil {
		ErrorInResponse(w, http.StatusInternalServerError, JsonEncodeError)

		return
	}

	if _, err := w.Write(v); err != nil {
		ErrorInResponse(w, http.StatusInternalServerError, GetSongsListError)

		return
	}
}

// idで指定した曲を返す.
func (sc *SongController) GetSongHandler(w http.ResponseWriter, r *http.Request) {
	// 対象の曲idを取得.
	songID, errorSet := GetId(r)

	if errorSet != nil {
		ErrorInResponse(w, errorSet.StatusCode, errorSet.Message)
		log.Printf("%v", errorSet.Err)

		return
	}

	song, err := sc.SongInteractor.Show(songID)

	if err != nil {
		ErrorInResponse(w, http.StatusInternalServerError, GetSongError)
		log.Printf("%v", err)

		return
	}

	v, err := json.Marshal(song)

	if err != nil {
		ErrorInResponse(w, http.StatusInternalServerError, JsonEncodeError)
		log.Printf("%v", errors.WithStack(err))

		return
	}

	if _, err := w.Write(v); err != nil {
		ErrorInResponse(w, http.StatusInternalServerError, GetSongDetailError)
		log.Printf("%v", errors.WithStack(err))

		return
	}
}

// 曲を追加.
func (sc *SongController) CreateSongHandler(w http.ResponseWriter, r *http.Request) {
	// リクエストユーザーのメアドを取得.
	userEmail, errorSet := GetEmail(r)

	if errorSet != nil {
		ErrorInResponse(w, errorSet.StatusCode, errorSet.Message)

		return
	}

	var d model.Song

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		ErrorInResponse(w, http.StatusInternalServerError, DecodeError)

		return
	}

	if err := sc.SongInteractor.Store(userEmail, d); err != nil {
		ErrorInResponse(w, http.StatusInternalServerError, PostSongError)

		return
	}

	// 201 Created
	w.WriteHeader(201)
}

// idで指定した曲の情報を更新.
func (sc *SongController) UpdateSongHandler(w http.ResponseWriter, r *http.Request) {
	// リクエストユーザーのメアドと対象の曲idを取得.
	userEmail, songID, errorSet := GetEmailAndId(r)

	if errorSet != nil {
		ErrorInResponse(w, errorSet.StatusCode, errorSet.Message)

		return
	}

	var d model.Song

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		ErrorInResponse(w, http.StatusInternalServerError, DecodeError)

		return
	}

	if err := sc.SongInteractor.Update(userEmail, songID, d); err != nil {
		ErrorInResponse(w, http.StatusInternalServerError, UpdateSongError)

		return
	}

	// 204 No Content
	w.WriteHeader(204)
}

// idで指定した曲を削除.
func (sc *SongController) DeleteSongHandler(w http.ResponseWriter, r *http.Request) {
	// 対象の曲idを取得.
	songID, errorSet := GetId(r)

	if errorSet != nil {
		ErrorInResponse(w, errorSet.StatusCode, errorSet.Message)

		return
	}

	if err := sc.SongInteractor.Destroy(songID); err != nil {
		ErrorInResponse(w, http.StatusInternalServerError, DeleteSongError)

		return
	}

	// 204 No Content
	w.WriteHeader(204)
}
