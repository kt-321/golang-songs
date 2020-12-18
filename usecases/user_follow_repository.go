package usecases

type UserFollowRepositoryInterface interface {
	Follow(string, int) error
	Unfollow(string, int) error
}
