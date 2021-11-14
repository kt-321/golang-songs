package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"golang-songs/interfaces"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	migrate "github.com/rubenv/sql-migrate"

	stats_api "github.com/fukata/golang-stats-api-handler"

	"github.com/jinzhu/gorm"
	"golang.org/x/sync/errgroup"

	"github.com/garyburd/redigo/redis"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"

	"golang-songs/queries/usersQuery"

	"github.com/gorilla/mux"
)

func Dispatch(DB *gorm.DB, Redis redis.Conn, SidecarRedis redis.Conn) {
	authController := interfaces.NewAuthController(DB)
	userController := interfaces.NewUserController(DB)
	usersController := usersQuery.NewUserController(DB)
	songController := interfaces.NewSongController(DB, Redis, SidecarRedis)
	bookmarkController := interfaces.NewBookmarkController(DB)
	userFollowController := interfaces.NewUserFollowController(DB)
	spotifyController := interfaces.NewSpotifyController(DB)

	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	s.HandleFunc("/signup", authController.SignUpHandler).Methods("POST")
	s.HandleFunc("/login", authController.LoginHandler).Methods("POST")

	s.Handle("/user", JwtMiddleware.Handler(http.HandlerFunc(usersController.UserHandler))).Methods("GET")
	s.Handle("/users", JwtMiddleware.Handler(http.HandlerFunc(usersController.AllUsersHandler))).Methods("GET")
	s.Handle("/user/{id}", JwtMiddleware.Handler(http.HandlerFunc(usersController.GetUserHandler))).Methods("GET")
	//s.Handle("/user", JwtMiddleware.Handler(http.HandlerFunc(userController.UserHandler))).Methods("GET")
	//s.Handle("/users", JwtMiddleware.Handler(http.HandlerFunc(userController.AllUsersHandler))).Methods("GET")
	//s.Handle("/user/{id}", JwtMiddleware.Handler(http.HandlerFunc(userController.GetUserHandler))).Methods("GET")

	//TODO
	s.Handle("/user/{id}/update", JwtMiddleware.Handler(http.HandlerFunc(userController.UpdateUserHandler))).Methods("PUT")

	s.Handle("/songs", JwtMiddleware.Handler(http.HandlerFunc(songController.AllSongsHandler))).Methods("GET")
	s.Handle("/song/{id}", JwtMiddleware.Handler(http.HandlerFunc(songController.GetSongHandler))).Methods("GET")
	s.Handle("/song", JwtMiddleware.Handler(http.HandlerFunc(songController.CreateSongHandler))).Methods("POST")
	s.Handle("/song/{id}", JwtMiddleware.Handler(http.HandlerFunc(songController.UpdateSongHandler))).Methods("PUT")
	s.Handle("/song/{id}", JwtMiddleware.Handler(http.HandlerFunc(songController.DeleteSongHandler))).Methods("DELETE")

	s.Handle("/get-redirect-url", JwtMiddleware.Handler(http.HandlerFunc(spotifyController.GetRedirectURLHandler))).Methods("GET")
	s.Handle("/get-token", JwtMiddleware.Handler(http.HandlerFunc(spotifyController.GetTokenHandler))).Methods("POST")
	s.Handle("/tracks", JwtMiddleware.Handler(http.HandlerFunc(spotifyController.GetTracksHandler))).Methods("POST")

	s.Handle("/song/{id}/bookmark", JwtMiddleware.Handler(http.HandlerFunc(bookmarkController.BookmarkHandler))).Methods("POST")
	s.Handle("/song/{id}/remove-bookmark", JwtMiddleware.Handler(http.HandlerFunc(bookmarkController.RemoveBookmarkHandler))).Methods("POST")

	s.Handle("/user/{id}/follow", JwtMiddleware.Handler(http.HandlerFunc(userFollowController.FollowUserHandler))).Methods("POST")
	s.Handle("/user/{id}/unfollow", JwtMiddleware.Handler(http.HandlerFunc(userFollowController.UnfollowUserHandler))).Methods("POST")

	s.HandleFunc("/stats", stats_api.Handler)

	r.HandleFunc("/", healthzHandler).Methods("GET")

	//マイグレーション実行
	migrations := &migrate.FileMigrationSource{
		Dir: os.Getenv("migrationDir"),
	}

	db, err := sql.Open("mysql", os.Getenv("mysqlConfig"))
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	appliedCount, err := migrate.Exec(db, "mysql", migrations, migrate.Up)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Applied %v migrations", appliedCount)

	os.Exit(run(context.Background(), r))
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

// ELBのヘルスチェック用のハンドラ.
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func run(ctx context.Context, r *mux.Router) int {
	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)

	eg.Go(func() error {
		return runServer(ctx, r)
	})
	eg.Go(func() error {
		return acceptSignal(ctx)
	})
	eg.Go(func() error {
		<-ctx.Done()

		return ctx.Err()
	})

	if err := eg.Wait(); err != nil {
		log.Println(err)

		return 1
	}

	return 0
}

func acceptSignal(ctx context.Context) error {
	sigCh := make(chan os.Signal, 1)
	// signalをハンドリングする。SIGINTとSIGTERMの両方.
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		signal.Reset()

		return ctx.Err()
	case sig := <-sigCh:
		return fmt.Errorf("signal received: %v", sig.String())
	}
}

func runServer(ctx context.Context, r *mux.Router) error {
	s := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	errCh := make(chan error)

	go func() {
		defer close(errCh)

		if err := s.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		return s.Shutdown(ctx)
	case err := <-errCh:
		return err
	}
}
