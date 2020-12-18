package usecases

import "golang-songs/model"

type SpotifyInteractor struct {
	SpotifyRepository SpotifyRepositoryInterface
}

func (spi *SpotifyInteractor) GetTracks(token string, title string) (*model.Response, error) {
	return spi.SpotifyRepository.GetTracks(token, title)
}
