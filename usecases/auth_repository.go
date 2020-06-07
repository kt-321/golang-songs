package usecases

import "golang-songs/model"

type AuthRepository interface {
	SignUp(model.Form) error
	Login(model.Form) (*model.User, error)
}
