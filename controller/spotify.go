package controller

import (
	"bytes"
	"encoding/json"
	"golang-songs/model"
	"net/http"
	"os"

	"github.com/konojunya/musi/service"

	"github.com/jinzhu/gorm"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var config oauth2.Config

type Code struct {
	Code string
}

type Title struct {
	Title string
	Token string
}

//var OAuth = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	//func OAuth(c *gin.Context) {
//
//	//codeにClient IDとClient SecretをBase64でエンコードした値を入れる
//	dec := json.NewDecoder(r.Body)
//	var d Code
//	if err := dec.Decode(&d); err != nil {
//		var error model.Error
//		error.Message = "リクエストボディのデコードに失敗しました。"
//		errorInResponse(w, http.StatusInternalServerError, error)
//		return
//	}
//
//	code := d.Code
//
//	err := godotenv.Load()
//	if err != nil {
//		var error model.Error
//		error.Message = ".envファイルの読み込み失敗"
//		errorInResponse(w, http.StatusUnauthorized, error)
//		return
//	}
//
//	config = oauth2.Config{
//		ClientID:     os.Getenv("client_id"),
//		ClientSecret: os.Getenv("client_secret"),
//		Endpoint: oauth2.Endpoint{
//			AuthURL:  "https://accounts.spotify.com/authorize",
//			TokenURL: "https://accounts.spotify.com/api/token",
//		},
//		RedirectURL: "http://localhost:3000/spotify/songs",
//		Scopes:      []string{},
//	}
//
//	//認証のURL
//	//GetRedirectURL()の中でしたい
//	url := config.AuthCodeURL("test")
//
//	token, err := config.Exchange(oauth2.NoContext, code)
//
//	if err != nil {
//		log.Println(err)
//	}
//})

type GetToken struct {
	DB *gorm.DB
}

func (f *GetToken) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var d Code
	if err := dec.Decode(&d); err != nil {
		var error model.Error
		error.Message = "リクエストボディのデコードに失敗しました。"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	code := d.Code

	err := godotenv.Load()
	if err != nil {
		var error model.Error
		error.Message = ".envの読み込みに失敗しました。"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	config = oauth2.Config{
		ClientID:     os.Getenv("client_id"),
		ClientSecret: os.Getenv("client_secret"),
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://accounts.spotify.com/api/token",
		},
		RedirectURL: "http://localhost:3000/spotify/songs",
		Scopes:      []string{},
	}

	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		var error model.Error
		error.Message = "トークンの取得に失敗しました"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	v, err := json.Marshal(token.AccessToken)
	if err != nil {
		var error model.Error
		error.Message = "JSONへの変換に失敗しました"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	if _, err := w.Write(v); err != nil {
		var error model.Error
		error.Message = "URLの取得に失敗しました。"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}

type GetRedirectURL struct {
	DB *gorm.DB
}

func (f *GetRedirectURL) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		var error model.Error
		error.Message = ".envの読み込みに失敗しました。"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	config = oauth2.Config{
		ClientID:     os.Getenv("client_id"),
		ClientSecret: os.Getenv("client_secret"),
		Endpoint: oauth2.Endpoint{
			AuthURL: "https://accounts.spotify.com/authorize",
			//TokenURL: "https://accounts.spotify.com/api/token",
		},

		RedirectURL: "http://localhost:3000/spotify/songs",
		Scopes:      []string{},
	}

	url := config.AuthCodeURL("state")

	w.Header().Set("Content-Type", "application/json")

	if _, err := json.Marshal(url); err != nil {
		var error model.Error
		error.Message = "JSONへの変換に失敗しました。"
		//errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	v, err := JSONSafeMarshal(url, true)
	if err != nil {
		var error model.Error
		error.Message = "JSONへの変換に失敗しました"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	if _, err := w.Write(v); err != nil {
		var error model.Error
		error.Message = "URLの取得に失敗しました。"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}

func JSONSafeMarshal(v interface{}, safeEncoding bool) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		var error model.Error
		error.Message = "JSONへの変換に失敗しました。"
		//errorInResponse(w, http.StatusUnauthorized, error)
		//return
	}

	if safeEncoding {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}

type GetTracks struct {
	DB *gorm.DB
}

func (f *GetTracks) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//var GetTracks = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//一旦は直打ち OAuthで帰ってきたtoken: &{BQChNb6GK...の最初の部分
	//cookie := "AQBuZP2SIhzr3mTQrjKceOoK60cOGSHlVCpWRPH6pS_n13c8st-Oq18qcUuBMvuhKMKgtcKSxE4i01Vytt3M7y_78nTQqVcvpmzU1MjSJLSMrfLeF3p3yvxoYIe9Uf5PctWkjetxWD5WT11NptJ9ksmxo4AmuzIcRWN3IjyDiYnaBA7665016xVUlkQJFhHlKfXvpXSn3gCvJ_0"
	dec := json.NewDecoder(r.Body)
	var d Title
	if err := dec.Decode(&d); err != nil {
		var error model.Error
		error.Message = "リクエストボディのデコードに失敗しました。"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	cookie := d.Token
	title := d.Title

	if cookie == "" {
		var error model.Error
		error.Message = "アクセストークンの取得に失敗しました"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	//トラック（曲）検索
	//tracks, err := service.GetTracks(cookie)
	tracks, err := service.GetTracks(cookie, title)
	if err != nil {
		var error model.Error
		error.Message = "トラックの取得に失敗しました"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	v, err := json.Marshal(tracks)
	if err != nil {
		var error model.Error
		error.Message = "JSONへの変換に失敗しました"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
	w.Write(v)
}
