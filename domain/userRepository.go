package domain

import (
	"golang-songs/model"
)

type UserRepository struct {
	Da UserDataAccessor
}

type UserRepositoryInterface interface {
	Update(int, model.User) error
}

func (ur *UserRepository) Update(userID int, p model.User) error {
	if err := ur.Da.updateUser(userID, p); err != nil {
		return err
	}

	return nil
}
