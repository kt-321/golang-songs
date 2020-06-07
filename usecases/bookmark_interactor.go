package usecases

type BookmarkInteractor struct {
	BookmarkRepository BookmarkRepository
}

func (bi *BookmarkInteractor) Bookmark(userEmail string, songID int) error {
	err := bi.BookmarkRepository.Bookmark(userEmail, songID)

	return err
}

func (bi *BookmarkInteractor) RemoveBookmark(userEmail string, songID int) error {
	err := bi.BookmarkRepository.RemoveBookmark(userEmail, songID)

	return err
}
