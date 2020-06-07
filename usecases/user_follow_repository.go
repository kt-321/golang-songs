package usecases

type UserFollowRepository interface {
	Follow(string, int) error
	Unfollow(string, int) error
}
