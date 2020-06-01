package interfaces

import (
	"encoding/json"

	"github.com/jinzhu/gorm"

	//"golang-songs/domain"
	"golang-songs/domain"
	"golang-songs/usecases"
	"net/http"
	"strconv"
)

type SongController struct {
	SongInteractor usecases.SongInteractor
	Logger         usecases.Logger
}

//func NewSongController(sqlHandler SQLHandler, logger usecases.Logger) *SongController {
func NewSongController(DB *gorm.DB) *SongController {
	return &SongController{
		SongInteractor: usecases.SongInteractor{
			SongRepository: &SongRepository{
				DB: DB,
			},
		},
		//Logger: logger,
	}
}

// Index is display a listing of the resource.
func (pc *SongController) Index(w http.ResponseWriter, r *http.Request) {
	pc.Logger.LogAccess("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	songs, err := pc.SongInteractor.Index()
	if err != nil {
		pc.Logger.LogError("%s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

// Store is stora a newly created resource in storage.
func (pc *SongController) Store(w http.ResponseWriter, r *http.Request) {
	pc.Logger.LogAccess("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	p := domain.Song{}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		pc.Logger.LogError("%s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
	}

	song, err := pc.SongInteractor.Store(p)
	if err != nil {
		pc.Logger.LogError("%s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

// Destroy is remove the specified resource from storage.
func (pc *SongController) Destroy(w http.ResponseWriter, r *http.Request) {
	pc.Logger.LogAccess("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	songID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	err := pc.SongInteractor.Destroy(songID)
	if err != nil {
		pc.Logger.LogError("%s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
}
