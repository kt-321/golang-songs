package usecases

type BookmarkRepository interface {
	Bookmark(string, int) error
	RemoveBookmark(string, int) error
}
