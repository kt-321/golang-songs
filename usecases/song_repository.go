package usecases

import (
	//"golang-songs/domain"
	"golang-songs/model"
)

// A PostRepository belong to the usecases layer.
type SongRepository interface {
	//FindAll() (model.Songs, error)
	FindAll(string) (*model.Songs, error)
	FindByID(string, int)
	Save(string, model.Song) (int64, error)
	//UpdateByID(int) error
	UpdateByID(string, int, model.Song) error
	//DeleteByID(int) error
	DeleteByID(string, int) error
}
