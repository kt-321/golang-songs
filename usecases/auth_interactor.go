package usecases

import "golang-songs/model"

type AuthInteractor struct {
	AuthRepository AuthRepository
}

func (ai *AuthInteractor) SignUp(p model.Form) error {
	return ai.AuthRepository.SignUp(p)
}

func (ai *AuthInteractor) Login(p model.Form) (*model.User, error) {
	return ai.AuthRepository.Login(p)
}
