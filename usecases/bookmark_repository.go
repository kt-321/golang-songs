package usecases

type BookmarkRepositoryInterface interface {
	Bookmark(string, int) error
	RemoveBookmark(string, int) error
}
