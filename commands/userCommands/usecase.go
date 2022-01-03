package userCommands

import (
	"golang-songs/domain"
	"golang-songs/model"
)

type usecase struct {
	UserRepository domain.UserRepositoryInterface
}

type Usecase interface {
	Update(int, model.User) error
}

func (ui *usecase) Update(userID int, p model.User) error {
	return ui.UserRepository.Update(userID, p)
}
