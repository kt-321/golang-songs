package usecases

import "golang-songs/model"

type UserRepository interface {
	FindAll() (*[]model.User, error)
	GetUser(string) (*model.User, error)
	FindByID(int) (*model.User, error)
	Update(int, model.User) error
}
