package infrastructure

import (
	"golang-songs/interfaces"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"

	"github.com/garyburd/redigo/redis"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"
)

func Dispatch(DB *gorm.DB, Redis redis.Conn) {
	authController := interfaces.NewAuthController(DB)
	userController := interfaces.NewUserController(DB)
	songController := interfaces.NewSongController(DB, Redis)
	bookmarkController := interfaces.NewBookmarkController(DB)
	userFollowController := interfaces.NewUserFollowController(DB)
	spotifyController := interfaces.NewSpotifyController(DB)

	r := mux.NewRouter()

	r.HandleFunc("/api/signup", authController.SignUpHandler).Methods("POST")
	r.HandleFunc("/api/login", authController.LoginHandler).Methods("POST")

	r.Handle("/api/user", JwtMiddleware.Handler(http.HandlerFunc(userController.UserHandler))).Methods("GET")
	r.Handle("/api/users", JwtMiddleware.Handler(http.HandlerFunc(userController.AllUsersHandler))).Methods("GET")
	r.Handle("/api/user/{id}", JwtMiddleware.Handler(http.HandlerFunc(userController.GetUserHandler))).Methods("GET")
	r.Handle("/api/user/{id}/update", JwtMiddleware.Handler(http.HandlerFunc(userController.UpdateUserHandler))).Methods("PUT")

	r.Handle("/api/songs", JwtMiddleware.Handler(http.HandlerFunc(songController.AllSongsHandler))).Methods("GET")
	r.Handle("/api/song/{id}", JwtMiddleware.Handler(http.HandlerFunc(songController.GetSongHandler))).Methods("GET")
	r.Handle("/api/song", JwtMiddleware.Handler(http.HandlerFunc(songController.CreateSongHandler))).Methods("POST")
	r.Handle("/api/song/{id}", JwtMiddleware.Handler(http.HandlerFunc(songController.UpdateSongHandler))).Methods("PUT")
	r.Handle("/api/song/{id}", JwtMiddleware.Handler(http.HandlerFunc(songController.DeleteSongHandler))).Methods("DELETE")

	r.Handle("/api/get-redirect-url", JwtMiddleware.Handler(http.HandlerFunc(spotifyController.GetRedirectURLHandler))).Methods("GET")
	r.Handle("/api/get-token", JwtMiddleware.Handler(http.HandlerFunc(spotifyController.GetTokenHandler))).Methods("POST")
	r.Handle("/api/tracks", JwtMiddleware.Handler(http.HandlerFunc(spotifyController.GetTracksHandler))).Methods("POST")

	r.Handle("/api/song/{id}/bookmark", JwtMiddleware.Handler(http.HandlerFunc(bookmarkController.BookmarkHandler))).Methods("POST")
	r.Handle("/api/song/{id}/remove-bookmark", JwtMiddleware.Handler(http.HandlerFunc(bookmarkController.RemoveBookmarkHandler))).Methods("POST")

	r.Handle("/api/user/{id}/follow", JwtMiddleware.Handler(http.HandlerFunc(userFollowController.FollowUserHandler))).Methods("POST")
	r.Handle("/api/user/{id}/unfollow", JwtMiddleware.Handler(http.HandlerFunc(userFollowController.UnfollowUserHandler))).Methods("POST")

	r.HandleFunc("/", healthzHandler).Methods("GET")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Println(err)
	}
}

var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		secret := os.Getenv("SIGNINGKEY")

		if os.Getenv("SIGNINGKEY") == "" {
			panic("環境変数SIGNINGKEYが存在しません。")
		}
		return []byte(secret), nil
	},
	Debug:         true,
	SigningMethod: jwt.SigningMethodHS256,
})

//ELBのヘルスチェック用のハンドラ
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
