package usecases

import "golang-songs/domain"

// A PostRepository belong to the usecases layer.
type SongRepository interface {
	FindAll() (domain.Songs, error)
	Save(domain.Song) (int64, error)
	DeleteByID(int) error
}
