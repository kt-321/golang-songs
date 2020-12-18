package usecases

import "golang-songs/model"

type SpotifyRepositoryInterface interface {
	GetTracks(token string, title string) (*model.Response, error)
}
