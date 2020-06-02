package interfaces

import (
	"encoding/json"
	"golang-songs/model"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"

	//"golang-songs/domain"
	"golang-songs/usecases"
	"net/http"
	//"strconv"
)

type SongController struct {
	SongInteractor usecases.SongInteractor
	//Logger         usecases.Logger
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
func (sc *SongController) Index(w http.ResponseWriter, r *http.Request) {
	//pc.Logger.LogAccess("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	songs, err := sc.SongInteractor.Index()
	if err != nil {
		//pc.Logger.LogError("%s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

// Store is stora a newly created resource in storage.
func (sc *SongController) Store(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var d model.Song

	if err := dec.Decode(&d); err != nil {
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

	//userEmailを引数として渡すように（返すように）かえるか
	//var user model.User
	//
	//if err := f.DB.Where("email = ?", userEmail).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
	//	error := model.Error{}
	//	error.Message = "該当するアカウントが見つかりません。"
	//	errorInResponse(w, http.StatusUnauthorized, error)
	//	return
	//}

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
}

// Store is stora a newly created resource in storage.
func (sc *SongController) Update(w http.ResponseWriter, r *http.Request) {
	//pc.Logger.LogAccess("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	//p := domain.Song{}
	p := model.Song{}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		//pc.Logger.LogError("%s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
	}

	song, err := sc.SongInteractor.Update(p)
	if err != nil {
		//pc.Logger.LogError("%s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

// Destroy is remove the specified resource from storage.
func (sc *SongController) Destroy(w http.ResponseWriter, r *http.Request) {
	//pc.Logger.LogAccess("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	//songID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	vars := mux.Vars(r)
	songID, ok := vars["id"]
	if !ok {
		var error model.Error
		error.Message = "曲のidを取得できません。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	err := sc.SongInteractor.Destroy(songID)
	if err != nil {
		//pc.Logger.LogError("%s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
}

// Parse は jwt トークンから元になった認証情報を取り出す。
func Parse(signedString string) (*model.Auth, error) {
	//追加
	secret := os.Getenv("SIGNINGKEY")

	token, err := jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.Errorf("unexpected signing method: %v", token.Header)
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.Errorf("not found claims in %s", signedString)
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.Errorf("not found %s in %s", email, signedString)
	}

	return &model.Auth{Email: email}, nil
}
