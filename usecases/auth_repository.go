package usecases

import "golang-songs/model"

type AuthRepositoryInterface interface {
	SignUp(model.Form) error
	Login(model.Form) (*model.User, error)
}
