package usecases

import "golang-songs/domain"

// A PostRepository belong to the usecases layer.
type PostRepository interface {
	FindAll() (domain.Posts, error)
	Save(domain.Post) (int64, error)
	DeleteByID(int) error
}
