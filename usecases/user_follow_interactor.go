package usecases

type UserFollowInteractor struct {
	UserFollowRepository UserFollowRepository
}

func (ufi *UserFollowInteractor) Follow(requestUserEmail string, tagertUserID int) error {
	return ufi.UserFollowRepository.Follow(requestUserEmail, tagertUserID)
}

func (ufi *UserFollowInteractor) Unfollow(requestUserEmail string, tagertUserID int) error {
	return ufi.UserFollowRepository.Unfollow(requestUserEmail, tagertUserID)
}
