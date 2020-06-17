package interfaces

import (
	"encoding/json"
	"golang-songs/model"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"

	"golang-songs/usecases"
	"net/http"
)

type SongController struct {
	SongInteractor usecases.SongInteractor
}

func NewSongController(DB *gorm.DB) *SongController {
	return &SongController{
		SongInteractor: usecases.SongInteractor{
			SongRepository: &SongRepository{
				DB: DB,
			},
		},
	}
}

//全ての曲を返す
func (sc *SongController) AllSongsHandler(w http.ResponseWriter, r *http.Request) {
	songs, err := sc.SongInteractor.Index()
	if err != nil {
		var error model.Error
		error.Message = "曲が見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
	v, err := json.Marshal(songs)
	if err != nil {
		var error model.Error
		error.Message = "曲一覧の取得に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	if _, err := w.Write(v); err != nil {
		var error model.Error
		error.Message = "曲一覧の取得に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}

//idで指定した曲を返す
func (sc *SongController) GetSongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		var error model.Error
		error.Message = "曲のidを取得できません。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	songID, err := strconv.Atoi(id)
	if err != nil {
		var error model.Error
		error.Message = "idのint型への型変換に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	var song *model.Song

	song, err = sc.SongInteractor.Show(songID)
	if err != nil {
		var error model.Error
		error.Message = "該当する曲が見つかりません。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	v, err := json.Marshal(song)
	if err != nil {
		var error model.Error
		error.Message = "JSONへの変換に失敗しました"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	if _, err := w.Write(v); err != nil {
		var error model.Error
		error.Message = "曲情報の取得に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}

//曲を追加
func (sc *SongController) CreateSongHandler(w http.ResponseWriter, r *http.Request) {
	var d model.Song

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		var error model.Error
		error.Message = "リクエストボディのデコードに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	header_hoge := r.Header.Get("Authorization")
	bearerToken := strings.Split(header_hoge, " ")
	authToken := bearerToken[1]

	parsedToken, err := Parse(authToken)
	if err != nil {
		var error model.Error
		error.Message = "認証コードのパースに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	userEmail := parsedToken.Email

	if err := sc.SongInteractor.Store(userEmail, d); err != nil {
		var error model.Error
		error.Message = "曲の投稿に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}

//idで指定した曲の情報を更新
func (sc *SongController) UpdateSongHandler(w http.ResponseWriter, r *http.Request) {
	var d model.Song

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		var error model.Error
		error.Message = "リクエストボディのデコードに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	header_hoge := r.Header.Get("Authorization")
	bearerToken := strings.Split(header_hoge, " ")
	authToken := bearerToken[1]

	parsedToken, err := Parse(authToken)
	if err != nil {
		var error model.Error
		error.Message = "認証コードのパースに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	userEmail := parsedToken.Email

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		var error model.Error
		error.Message = "曲のidを取得できません。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}
	songID, err := strconv.Atoi(id)
	if err != nil {
		var error model.Error
		error.Message = "idのint型への型変換に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	if err := sc.SongInteractor.Update(userEmail, songID, d); err != nil {
		var error model.Error
		error.Message = "曲の更新に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}

//idで指定した曲を削除
func (sc *SongController) DeleteSongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		var error model.Error
		error.Message = "曲のidを取得できません。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}
	songID, err := strconv.Atoi(id)
	if err != nil {
		var error model.Error
		error.Message = "idのint型への型変換に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	if err := sc.SongInteractor.Destroy(songID); err != nil {
		var error model.Error
		error.Message = "曲の削除に失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}
