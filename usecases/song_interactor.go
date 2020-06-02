package usecases

import (
	"golang-songs/model"
)

// A PostInteractor belong to the usecases layer.
type SongInteractor struct {
	SongRepository SongRepository
}

// Index is display a listing of the resource.
func (si *SongInteractor) Index() (*model.Songs, error) {
	//func (si *SongInteractor) Index() {
	//songs, err = si.SongRepository.FindAll()
	songs, err := si.SongRepository.FindAll()

	return
}

// Store is store a newly created resource in storage.
//func (si *SongInteractor) Store(p model.Song) (id int64, err error) {
func (si *SongInteractor) Store(p model.Song) (int64, error) {
	id, err = si.SongRepository.Save(p)

	return
}

//func (pi *SongInteractor) Update(songID int) (err error) {
func (si *SongInteractor) Update(songID int) error {
	err = si.SongRepository.UpdateByID(songID)

	return
}

// Destroy is remove the specified resource from storage.
func (si *SongInteractor) Destroy(songID int) (err error) {
	err = si.SongRepository.DeleteByID(songID)

	return
}
