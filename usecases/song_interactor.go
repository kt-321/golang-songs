package usecases

import (
	"golang-songs/model"
)

type SongInteractor struct {
	SongRepository SongRepository
}

func (si *SongInteractor) Index() (*[]model.Song, error) {
	return si.SongRepository.FindAll()
}

func (si *SongInteractor) Show(songID int) (*model.Song, error) {
	return si.SongRepository.FindByID(songID)
}

func (si *SongInteractor) Store(userEmail string, p model.Song) error {
	return si.SongRepository.Save(userEmail, p)
}

func (si *SongInteractor) Update(userEmail string, songID int, p model.Song) error {
	return si.SongRepository.UpdateByID(userEmail, songID, p)
}

func (si *SongInteractor) Destroy(songID int) error {
	return si.SongRepository.DeleteByID(songID)
}
