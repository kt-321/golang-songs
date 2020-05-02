package controller

import (
	"encoding/json"
	"fmt"
	"golang-songs/service"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"golang.org/x/oauth2"
)

//追加した
var config oauth2.Config

type Code struct {
	Code string
}

type Title struct {
	Title string
}

//gin使うか
//func Login(c *gin.Context) {
//	c.Redirect(http.StatusTemporaryRedirect, service.GetRedirectURL())
//}

//コンテキスト使うのか
//func OAuth(w http.ResponseWriter, r *http.Request) {
//func OAuth(ctx context.Context) {

//code := ctx.Value("code")

//gin使わない形に書き換えるか

var OAuth = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//func OAuth(c *gin.Context) {

	//codeにClient IDとClient SecretをBase64でエンコードした値を入れる

	dec := json.NewDecoder(r.Body)
	var d Code
	dec.Decode(&d)

	log.Println("d:", d)
	code := d.Code
	log.Println("code:", code)
	//code := c.Query("code")

	//e := r.ParseForm()
	//log.Println("e:", e)
	//code := r.FormValue("code")
	//fmt.Println("code:", code)

	//一旦は直打ち Client IDとClient SecretをBase64でエンコードした値
	//code = "ZjNmN2YxZWUxY2EwNGU3NzllYWRlZDBkMTM3NTJiNWY6ZmFlZDk1NzYwYWJhNGZhNGFhYTc1NjI2MGIxMjMxOTU="

	err := godotenv.Load()
	if err != nil {
		log.Println(".envファイルの読み込み失敗")
	}

	config = oauth2.Config{
		ClientID:     os.Getenv("client_id"),
		ClientSecret: os.Getenv("client_secret"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},

		//RedirectURL: "https://musi-app.now.sh/oauth",
		RedirectURL: "http://localhost:8081/api/oauth",
		//RedirectURL: "http://localhost:3000/dashboard",
		//今回はリダイレクトしない
		//RedirectURL: "urn:ietf:wg:oauth:2.0:oob",
		//Scopes: []string{"playlist-modify", "user-read-private", "user-library-read"},
		//Scopes: []string{"http://localhost:8080"},
		Scopes: []string{},
	}

	//追加した 必要かわからん
	//url := config.AuthCodeURL("test")

	//認証のURL
	//url := config.AuthCodeURL("state")
	url := config.AuthCodeURL("test")
	fmt.Println("url:", url)

	//以下の処理をやっているのか
	//$ echo -n f3f7f1ee1ca04e779eaded0d13752b5f:faed95760aba4fa4aaa756260b123195 | base64
	//やってみるとClient IDとClient SecretをBase64でエンコードした値を取得した以下の値で、
	//ZjNmN2YxZWUxY2EwNGU3NzllYWRlZDBkMTM3NTJiNWY6ZmFlZDk1NzYwYWJhNGZhNGFhYTc1NjI2MGIxMjMxOTU=
	//→ これがcodeか
	//アクセストークンを取得しにいく
	//$ curl -X "POST" -H "Authorization: Basic ZjNmN2YxZWUxY2EwNGU3NzllYWRlZDBkMTM3NTJiNWY6ZmFlZDk1NzYwYWJhNGZhNGFhYTc1NjI2MGIxMjMxOTU=" -d grant_type=client_credentials https://accounts.spotify.com/api/token

	log.Println("config:", config)

	//認証URLをクリック後に表示された（今回はURL欄に）code
	testcode := "AQAmqBWKde1Ml2Q-O5XmGpF5PlJ2ymH10T2jH63j3XI9fw96ERl52Hu65w1xbr8cFmZnOtNnvXjxPPKWRobixPWQg82EKtjvJfRZx9LJO6cnc44iBoH8c7qCG9dTaPM0GviQR0WJTpkxPcmHVor25mXr_NtlZm5Np3IAX0VYiD6Mek2h8qY1PLF02BOu72_2F4p5n54"

	//token, err := config.Exchange(oauth2.NoContext, code)
	token, err := config.Exchange(oauth2.NoContext, testcode)

	//→よくわからんができたかもしれない https://kido0617.github.io/go/2016-07-18-oauth2/

	//token, err := config.Exchange(oauth2.Endpoint.TokenURL, code)

	//log.Println("oauth2.NoContext:", oauth2.NoContext)
	//log.Println("context.TODO:", context.TODO)

	if err != nil {
		log.Fatal(err)
		log.Fatal("controller//token取得失敗")
	}

	log.Println("token:", token)

	//帰ってきた例
	//token: &{BQChNb6GKt6oM-qey9Yg8G94OqsY53hS13YTO621EmFWTZaF8A4OYP8QenmmxlUFavLpTaqU79qkb32rGGNJy-R6jYSJf9DODz-XEE70uEUdDkPIeQZzP1pDYzMy8htRCkLPw8JV2nhXTOvkYdAeKI7qdxPUTMmR-Q Bearer AQDU17AbxZo5Sk2L8ntUbW-G2-0UYkCbv78OIssVchiWldM7tx-rxgbYC8Qn3sdQVTc-1PHtqdWNCAw-OAA1sDk3TL58zb6bQB9dPX1slwdMKn_JRKuDGg9S7rp8-lBFE-c 2020-04-30 10:15:58.37369 +0900 JST m=+3606.898030009 map[access_token:BQChNb6GKt6oM-qey9Yg8G94OqsY53hS13YTO621EmFWTZaF8A4OYP8QenmmxlUFavLpTaqU79qkb32rGGNJy-R6jYSJf9DODz-XEE70uEUdDkPIeQZzP1pDYzMy8htRCkLPw8JV2nhXTOvkYdAeKI7qdxPUTMmR-Q expires_in:3600 refresh_token:AQDU17AbxZo5Sk2L8ntUbW-G2-0UYkCbv78OIssVchiWldM7tx-rxgbYC8Qn3sdQVTc-1PHtqdWNCAw-OAA1sDk3TL58zb6bQB9dPX1slwdMKn_JRKuDGg9S7rp8-lBFE-c scope: token_type:Bearer]}

	//追記した
	//client := config.Client(oauth2.NoContext, token) //httpクライアントを取得
	//log.Println(client)

	//if err != nil {
	//	//c.AbortWithError(http.StatusInternalServerError, err)
	//	//return
	//	log.Fatal("token取得失敗")
	//}
	log.Println("token.AccessToken:", token.AccessToken)

	//別サイト参考に作った
	//cookie := &http.Cookie{
	//	Name:  "spotify-token",
	//	Value: token.AccessToken,
	//}

	//log.Println(cookie)
	//http.SetCookie(w, cookie)

	//fmt.Fprintf(w, "Cookieの設定ができたよ")

	//要修正か
	//c.SetCookie("spotify-token", token.AccessToken, 1000*60*60*24*7, "/", "https://musi-app.now.sh", false, false)

	//これはそのアプリケーションでの話か
	//c.SetCookie("spotify-token", token.AccessToken, 1000*60*60*24*7, "/", "http://localhost:8081", false, false)
	//c.Redirect(http.StatusTemporaryRedirect, "/")
})

//func GetArtist(c *gin.Context) {
//
//	//一旦はここでやる
//	//c.SetCookie("spotify-token", token.AccessToken, 1000*60*60*24*7, "/", "https://localhost:8080", false, false)
//	//c.SetCookie("spotify-token", "BQAQ8IXL2MkfLMrNWK0Bs-b_jRdVIVfFMZvTSk_eCifwXedDodKrxch5hCRAZQXDyWdCa1HQv7F9RZJc2mw", 1000*60*60*24*7, "/", "http://localhost:8080", false, false)
//	c.SetCookie("spotify-token", "BQBVguyVRUe0r7dFQbxOkfQqHx7JE-O2Nk-y0F1h6Iy8DMFEVsBU0VdZAcnGWAt6tHIWcs_EA-cPRYPixPs", 1000*60*60*24*7, "/", "http://localhost:8080", false, false)
//
//	//gin使う場合
//	//spotify-tokenはアクセストークン
//	cookie, _ := c.Cookie("spotify-token")
//
//	if cookie == "" {
//		//gin使わない場合
//		//if cookie.Value == "" {
//
//		c.AbortWithStatus(http.StatusUnauthorized)
//		return
//		//log.Println("cookieが空白")
//	}
//	log.Println("cookie:", cookie)
//
//	//playlist, err := service.GetTracks(cookie, location)
//
//	//gin使う場合
//
//	//アーティスト検索
//	artists, err := service.SearchMusicArtists(cookie)
//
//	//gin使わない場合
//	//playlist, err := service.GetTracks(cookie.Value)
//	log.Println("artists:", artists)
//
//	if err != nil {
//		c.AbortWithStatus(http.StatusInternalServerError)
//		return
//		log.Fatal("アーティスト取得失敗")
//	}
//
//	//gin使う場合
//	c.JSON(http.StatusOK, artists)
//}

//アーティストの詳細をアーティストのIDでとるAPI（曲情報保存時にあらかじめ取っておく）

var GetTracks = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//func GetTracks(c *gin.Context) {
	//一旦はここでやる)

	//これらはフロントに保管しておくべきか
	//c.SetCookie("spotify-token", "BQAEBctWCsrvSH4UgnEnVbJsl3RWWpVKvSIW-uRKiQlamOo_Zw1yIRiqYMx3mCHFmVNksbAWTdw_DMEj_Nk", 1000*60*60*24*7, "/", "http://localhost:8080", false, false)
	//
	////gin使う場合
	//cookie, _ := c.Cookie("spotify-token")

	//一旦は直打ち OAuthで帰ってきたtoken: &{BQChNb6GK...の最初の部分
	cookie := "BQCV3weGc5JWT5hr3czciP-SatskSs49J9zGSmBwcgDy-LvSQltD4RDc3gdgGuBdt0lzBAyb9j5KEKa8mw6PaCxQ2V_NvgG_AEt3Lvs3UYxlWLcLxAimQOf3iqQmZlcZG3o_BaIHfMv7Z5XL7YR7myn3wSuvxRTscg"

	if cookie == "" {
		////gin使わない場合
		//c.AbortWithStatus(http.StatusUnauthorized)
		//return
		log.Println("cookieが空白")
	}
	log.Println("cookie:", cookie)

	//title := c.Query("title")

	//e := r.ParseForm()
	//log.Println(e)
	//title := r.FormValue("title")
	//fmt.Println("title", title)

	dec := json.NewDecoder(r.Body)
	var d Title
	dec.Decode(&d)

	log.Println("d:", d)
	title := d.Title
	log.Println("title:", title)

	//トラック（曲）検索
	//tracks, err := service.GetTracks(cookie)
	tracks, err := service.GetTracks(cookie, title)

	//log.Println("tracks:", tracks)
	log.Print("tracks:", tracks)

	if err != nil {
		//c.AbortWithStatus(http.StatusInternalServerError)
		//return
		log.Println("トラック取得失敗")
	}

	//gin使う場合
	//c.JSON(http.StatusOK, tracks)

	w.Header().Set("Content-Type", "application/json")

	v, err := json.Marshal(tracks)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(v)
})
