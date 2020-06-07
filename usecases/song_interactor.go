package usecases

import (
	"golang-songs/model"
)

type SongInteractor struct {
	SongRepository SongRepository
}

func (si *SongInteractor) Index() (*[]model.Song, error) {
	songs, err := si.SongRepository.FindAll()

	return songs, err
}

func (si *SongInteractor) Show(songID int) (*model.Song, error) {
	song, err := si.SongRepository.FindByID(songID)

	return song, err
}

func (si *SongInteractor) Store(userEmail string, p model.Song) error {
	err := si.SongRepository.Save(userEmail, p)

	return err
}

func (si *SongInteractor) Update(userEmail string, songID int, p model.Song) error {
	err := si.SongRepository.UpdateByID(userEmail, songID, p)

	return err
}

func (si *SongInteractor) Destroy(songID int) error {
	err := si.SongRepository.DeleteByID(songID)

	return err
}
