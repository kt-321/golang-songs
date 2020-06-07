package interfaces

import (
	"golang-songs/model"

	"github.com/jinzhu/gorm"
)

type AuthRepository struct {
	DB *gorm.DB
}

func (ar *AuthRepository) SignUp(form model.Form) error {
	if err := ar.DB.Create(&model.User{Email: form.Email, Password: form.Password}).Error; err != nil {
		//var error model.Error
		//error.Message = "アカウントの作成に失敗しました"
		//errorInResponse(w, http.StatusUnauthorized, error)
		return err
	}
	return nil
}

func (ar *AuthRepository) Login(form model.Form) (*model.User, error) {
	var user model.User

	if err := ar.DB.Where("email = ?", form.Email).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		//var error model.Error
		//error.Message = "該当するアカウントが見つかりません。"
		//errorInResponse(w, http.StatusUnauthorized, error)
		return nil, err
	}
	return &user, nil
}
