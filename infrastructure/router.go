package infrastructure

import (
	"golang-songs/interfaces"
	"log"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"
)

// Dispatch is handle routing
//func Dispatch(logger usecases.Logger, sqlHandler interfaces.SQLHandler) {
func Dispatch() {
	//userController := interfaces.NewUserController(sqlHandler, logger)
	userController := interfaces.NewUserController(DB)
	//userController := interfaces.NewUserController(DB)
	//songController := interfaces.NewSongController(sqlHandler, logger)
	songController := interfaces.NewSongController()

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

	//r.Handle("/api/signup", &SignUpHandler{DB: db}).Methods("POST")
	//r.Handle("/api/login", &LoginHandler{DB: db}).Methods("POST")

	//r.HandleFunc("/api/get-redirect-url", controller.GetRedirectURL).Methods("GET")
	//r.HandleFunc("/api/get-token", controller.GetToken).Methods("POST")
	//r.HandleFunc("/api/tracks", controller.GetTracks).Methods("POST")

	//r.Handle("/api/song/{id}/bookmark", JwtMiddleware.Handler(&BookmarkHandler{DB: db})).Methods("POST")
	//r.Handle("/api/song/{id}/remove-bookmark", JwtMiddleware.Handler(&RemoveBookmarkHandler{DB: db})).Methods("POST")

	//r.Handle("/api/user/{id}/follow", JwtMiddleware.Handler(&FollowUserHandler{DB: db})).Methods("POST")
	//r.Handle("/api/user/{id}/unfollow", JwtMiddleware.Handler(&UnfollowUserHandler{DB: db})).Methods("POST")
	//
	//r.HandleFunc("/", healthzHandler).Methods("GET")

	//r.Handle("/api/user", JwtMiddleware.Handler(&UserHandler{DB: db})).Methods("GET")
	//r.Handle("/api/user/{id}", JwtMiddleware.Handler(&GetUserHandler{DB: db})).Methods("GET")
	//r.Handle("/api/users", JwtMiddleware.Handler(&AllUsersHandler{DB: db})).Methods("GET")
	//r.Handle("/api/user/{id}/update", JwtMiddleware.Handler(&UpdateUserHandler{DB: db})).Methods("PUT")
	//
	//r.Handle("/api/song", JwtMiddleware.Handler(&CreateSongHandler{DB: db})).Methods("POST")
	//r.Handle("/api/song/{id}", JwtMiddleware.Handler(&GetSongHandler{DB: db})).Methods("GET")
	//r.Handle("/api/songs", JwtMiddleware.Handler(&AllSongsHandler{DB: db})).Methods("GET")
	//r.Handle("/api/song/{id}", JwtMiddleware.Handler(&UpdateSongHandler{DB: db})).Methods("PUT")
	//r.Handle("/api/song/{id}", JwtMiddleware.Handler(&DeleteSongHandler{DB: db})).Methods("DELETE")
	//
	//r.HandleFunc("/api/get-redirect-url", controller.GetRedirectURL).Methods("GET")
	//r.HandleFunc("/api/get-token", controller.GetToken).Methods("POST")
	//r.HandleFunc("/api/tracks", controller.GetTracks).Methods("POST")
	//
	//r.Handle("/api/song/{id}/bookmark", JwtMiddleware.Handler(&BookmarkHandler{DB: db})).Methods("POST")
	//r.Handle("/api/song/{id}/remove-bookmark", JwtMiddleware.Handler(&RemoveBookmarkHandler{DB: db})).Methods("POST")
	//
	//r.Handle("/api/user/{id}/follow", JwtMiddleware.Handler(&FollowUserHandler{DB: db})).Methods("POST")
	//r.Handle("/api/user/{id}/unfollow", JwtMiddleware.Handler(&UnfollowUserHandler{DB: db})).Methods("POST")
	//
	//r.HandleFunc("/", healthzHandler).Methods("GET")

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
