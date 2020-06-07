package usecases

type UserFollowInteractor struct {
	UserFollowRepository UserFollowRepository
}

func (ufi *UserFollowInteractor) Follow(requestUserEmail string, tagertUserID int) error {
	err := ufi.UserFollowRepository.Follow(requestUserEmail, tagertUserID)

	return err
}

func (ufi *UserFollowInteractor) Unfollow(requestUserEmail string, tagertUserID int) error {
	err := ufi.UserFollowRepository.Unfollow(requestUserEmail, tagertUserID)

	return err
}
