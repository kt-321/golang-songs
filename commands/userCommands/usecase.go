package userCommands

import (
	"golang-songs/domain"
	"golang-songs/model"
)

type Usecase struct {
	UserRepository domain.UserRepositoryInterface
}

func (ui *Usecase) Update(userID int, p model.User) error {
	return ui.UserRepository.Update(userID, p)
}
