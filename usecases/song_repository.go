package usecases

import (
	//"golang-songs/domain"
	"golang-songs/model"
)

type SongRepositoryInterface interface {
	FindAll() (*[]model.Song, error)
	FindByID(int) (*model.Song, error)
	Save(string, model.Song) error
	UpdateByID(string, int, model.Song) error
	DeleteByID(int) error
}
