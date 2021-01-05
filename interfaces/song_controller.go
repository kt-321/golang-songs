package interfaces

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"golang-songs/model"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"

	"golang-songs/usecases"
	"net/http"
)

type SongController struct {
	SongInteractor usecases.SongInteractor
}

func NewSongController(DB *gorm.DB, Redis redis.Conn, SidecarRedis redis.Conn) *SongController {
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
		errorInResponse(w, http.StatusInternalServerError, 11)

		return
	}

	v, err := json.Marshal(songs)

	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 12)

		return
	}

	if _, err := w.Write(v); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 12)

		return
	}
}

// idで指定した曲を返す.
func (sc *SongController) GetSongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		errorInResponse(w, http.StatusBadRequest, 13)

		return
	}

	songID, err := strconv.Atoi(id)
	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 14)

		return
	}

	var song *model.Song

	song, err = sc.SongInteractor.Show(songID)
	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 15)

		return
	}

	v, err := json.Marshal(song)
	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 6)

		return
	}

	if _, err := w.Write(v); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 16)

		return
	}
}

// 曲を追加.
func (sc *SongController) CreateSongHandler(w http.ResponseWriter, r *http.Request) {
	userEmail, errorSet := GetEmail(r)

	if errorSet != nil {
		errorInResponse(w, errorSet.StatusCode, errorSet.MessageNumber)

		return
	}

	var d model.Song

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 17)

		return
	}

	if err := sc.SongInteractor.Store(userEmail, d); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 19)

		return
	}

	// 201 Created
	w.WriteHeader(201)
}

// idで指定した曲の情報を更新.
func (sc *SongController) UpdateSongHandler(w http.ResponseWriter, r *http.Request) {
	userEmail, errorSet := GetEmail(r)

	if errorSet != nil {
		errorInResponse(w, errorSet.StatusCode, errorSet.MessageNumber)

		return
	}

	var d model.Song

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 17)

		return
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		errorInResponse(w, http.StatusBadRequest, 13)

		return
	}

	songID, err := strconv.Atoi(id)
	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 14)

		return
	}

	if err := sc.SongInteractor.Update(userEmail, songID, d); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 20)

		return
	}

	// 204 No Content
	w.WriteHeader(204)
}

// idで指定した曲を削除.
func (sc *SongController) DeleteSongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		errorInResponse(w, http.StatusBadRequest, 13)

		return
	}

	songID, err := strconv.Atoi(id)
	if err != nil {
		errorInResponse(w, http.StatusInternalServerError, 14)

		return
	}

	if err := sc.SongInteractor.Destroy(songID); err != nil {
		errorInResponse(w, http.StatusInternalServerError, 21)

		return
	}

	// 204 No Content
	w.WriteHeader(204)
}
