package usecases

import "golang-songs/model"

// インタフェースAuthRepositoryInterfaceを満たしている
type AuthInteractor struct {
	AuthRepository AuthRepositoryInterface
}

func (ai *AuthInteractor) SignUp(p model.Form) error {
	return ai.AuthRepository.SignUp(p)
}

func (ai *AuthInteractor) Login(p model.Form) (*model.User, error) {
	return ai.AuthRepository.Login(p)
}
