package usecases

import "golang-songs/domain"

// A PostInteractor belong to the usecases layer.
type SongInteractor struct {
	SongRepository SongRepository
}

// Index is display a listing of the resource.
func (pi *SongInteractor) Index() (posts domain.Songs, err error) {
	posts, err = pi.SongRepository.FindAll()

	return
}

// Store is store a newly created resource in storage.
func (pi *SongInteractor) Store(p domain.Post) (id int64, err error) {
	id, err = pi.SongRepository.Save(p)

	return
}

// Destroy is remove the specified resource from storage.
func (pi *SongInteractor) Destroy(postID int) (err error) {
	err = pi.SongRepository.DeleteByID(postID)

	return
}
