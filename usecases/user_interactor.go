package usecases

import (
	"golang-songs/model"
)

type UserInteractor struct {
	UserRepository UserRepository
}

func (ui *UserInteractor) Index() (*model.Users, error) {
	users, err := ui.UserRepository.FindAll()
	return users, err
}

func (ui *UserInteractor) User(userEmail string) (*model.User, error) {
	user, err := ui.UserRepository.GetUser(userEmail)

	return user, err
}

func (ui *UserInteractor) Show(userID int) (*model.User, error) {
	user, err := ui.UserRepository.FindByID(userID)

	return user, err
}

func (ui *UserInteractor) Update(userID int, p model.User) error {
	err := ui.UserRepository.Update(userID, p)

	return err
}
