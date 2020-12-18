package usecases

import (
	"golang-songs/model"
)

type UserInteractor struct {
	UserRepository UserRepositoryInterface
}

func (ui *UserInteractor) Index() (*[]model.User, error) {
	return ui.UserRepository.FindAll()
}

func (ui *UserInteractor) User(userEmail string) (*model.User, error) {
	return ui.UserRepository.GetUser(userEmail)
}

func (ui *UserInteractor) Show(userID int) (*model.User, error) {
	return ui.UserRepository.FindByID(userID)
}

func (ui *UserInteractor) Update(userID int, p model.User) error {
	return ui.UserRepository.Update(userID, p)
}
