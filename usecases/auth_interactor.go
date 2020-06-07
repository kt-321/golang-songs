package usecases

import "golang-songs/model"

type AuthInteractor struct {
	AuthRepository AuthRepository
}

func (ai *AuthInteractor) SignUp(p model.Form) error {
	err := ai.AuthRepository.SignUp(p)

	return err
}

func (ai *AuthInteractor) Login(p model.Form) (*model.User, error) {
	user, err := ai.AuthRepository.Login(p)

	return user, err
}
