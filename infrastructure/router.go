package infrastructure

import (
	"fmt"
	"golang-songs/interfaces"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"
)

// Dispatch is handle routing
//func Dispatch(logger usecases.Logger, sqlHandler interfaces.SQLHandler) {
//func Dispatch() {
func Dispatch(DB *gorm.DB) {
	//userController := interfaces.NewUserController(sqlHandler, logger)
	//userController := interfaces.NewUserController()
	userController := interfaces.NewUserController(DB)
	songController := interfaces.NewSongController(DB)

	r := mux.NewRouter()

	// http.HandlerFunc にキャスト
	// https: //gist.github.com/d-kuro/1b462a8b0544965046d63ea6ef5a0ec7
	//r.HandleFunc("/users", JwtMiddleware.Handler(http.HandlerFunc(userController.Index))).Methods("GET")
	r.Handle("/api/user", JwtMiddleware.Handler(http.HandlerFunc(userController.UserHandler))).Methods("GET")
	r.Handle("/api/users", JwtMiddleware.Handler(http.HandlerFunc(userController.Index))).Methods("GET")
	r.Handle("/api/user/{id}", JwtMiddleware.Handler(http.HandlerFunc(userController.Show))).Methods("GET")
	r.Handle("/api/user/{id}/update", JwtMiddleware.Handler(http.HandlerFunc(userController.Update))).Methods("PUT")

	r.Handle("/api/songs", JwtMiddleware.Handler(http.HandlerFunc(songController.Index))).Methods("GET")
	r.Handle("/api/song/{id}", JwtMiddleware.Handler(http.HandlerFunc(songController.Show))).Methods("GET")
	r.Handle("/api/song", JwtMiddleware.Handler(http.HandlerFunc(songController.Store))).Methods("POST")
	r.Handle("/api/song/{id}", JwtMiddleware.Handler(http.HandlerFunc(songController.Update))).Methods("PUT")
	r.Handle("/api/song/{id}", JwtMiddleware.Handler(http.HandlerFunc(songController.Destroy))).Methods("DELETE")

	r.Handle("/api/signup", JwtMiddleware.Handler(http.HandlerFunc(userController.SignUpHandler))).Methods("POST")
	r.Handle("/api/login", JwtMiddleware.Handler(http.HandlerFunc(userController.LoginHandler))).Methods("POST")

	r.Handle("/api/get-redirect-url", JwtMiddleware.Handler(http.HandlerFunc(songController.GetRedirectURL))).Methods("GET")
	r.Handle("/api/get-token", JwtMiddleware.Handler(http.HandlerFunc(songController.GetToken))).Methods("POST")
	r.Handle("/api/tracks", JwtMiddleware.Handler(http.HandlerFunc(songController.GetTracks))).Methods("POST")

	r.Handle("/api/song/{id}/bookmark", JwtMiddleware.Handler(http.HandlerFunc(songController.BookmarkHandler))).Methods("POST")
	r.Handle("/api/song/{id}/remove-bookmark", JwtMiddleware.Handler(http.HandlerFunc(songController.RemoveBookmarkHandler))).Methods("POST")

	r.Handle("/api/user/{id}/follow", JwtMiddleware.Handler(http.HandlerFunc(userController.FollowUserHandler))).Methods("POST")
	r.Handle("/api/user/{id}/unfollow", JwtMiddleware.Handler(http.HandlerFunc(userController.UnfollowUserHandler))).Methods("POST")

	r.HandleFunc("/", healthzHandler).Methods("GET")

	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Println(err)
	}

	//if err := http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), r); err != nil {
	//	logger.LogError("%s", err)
	//}
}

var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		secret := os.Getenv("SIGNINGKEY")
		return []byte(secret), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

//ELBのヘルスチェック用のハンドラ
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}
