package infrastructure

import (
	"golang-songs/interfaces"
	"golang-songs/usecases"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Dispatch is handle routing
func Dispatch(logger usecases.Logger, sqlHandler interfaces.SQLHandler) {
	userController := interfaces.NewUserController(sqlHandler, logger)
	songController := interfaces.NewSongController(sqlHandler, logger)

	//r := chi.NewRouter()
	r := mux.NewRouter()

	//r.Get("/users", userController.Index)
	r.HandleFunc("/users", userController.Index).Methods("GET")
	r.HandleFunc("/user", userController.Show).Methods("GET")
	//r.HandleFunc("/api/user/{id}", JwtMiddleware.userController.Show).Methods("GET")
	r.HandleFunc("/songs", songController.Index).Methods("GET")
	r.HandleFunc("/song", songController.Store).Methods("POST")
	r.HandleFunc("/song", songController.Destroy).Methods("DELETE")

	//r.Handle("/api/signup", &SignUpHandler{DB: db}).Methods("POST")
	//r.Handle("/api/login", &LoginHandler{DB: db}).Methods("POST")
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
