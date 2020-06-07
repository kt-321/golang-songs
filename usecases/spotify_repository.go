package usecases

import "golang-songs/model"

type SpotifyRepository interface {
	GetTracks(token string, title string) (*model.Response, error)
}
