package usecases

type BookmarkInteractor struct {
	BookmarkRepository BookmarkRepositoryInterface
}

func (bi *BookmarkInteractor) Bookmark(userEmail string, songID int) error {
	return bi.BookmarkRepository.Bookmark(userEmail, songID)
}

func (bi *BookmarkInteractor) RemoveBookmark(userEmail string, songID int) error {
	return bi.BookmarkRepository.RemoveBookmark(userEmail, songID)
}
