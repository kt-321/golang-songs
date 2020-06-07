package usecases

type BookmarkInteractor struct {
	BookmarkRepository BookmarkRepository
}

func (bi *BookmarkInteractor) Bookmark(userEmail string, songID int) error {
	return bi.BookmarkRepository.Bookmark(userEmail, songID)
}

func (bi *BookmarkInteractor) RemoveBookmark(userEmail string, songID int) error {
	return bi.BookmarkRepository.RemoveBookmark(userEmail, songID)
}
