package interfaces

import (
	"golang-songs/model"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// インタフェースAuthRepositoryInterfaceを満たす
type AuthRepository struct {
	DB *sqlx.DB
}

func (ar *AuthRepository) SignUp(form model.Form) error {
	//if err := ar.DB.Create(&model.User{Email: form.Email, Password: form.Password}).Error; err != nil {
	//	return err
	//}

	// https://codehex.hateblo.jp/entry/2018/05/21/100000
	tx := ar.DB.MustBegin()
	defer func() error {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return nil
	}()

	q := `insert ignore into users (email, password) values (?, ?);`

	tx.MustExec(q, form.Email, form.Password)

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (ar *AuthRepository) Login(form model.Form) (*model.User, error) {
	var user model.User

	if err := ar.DB.Where("email = ?", form.Email).Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		return nil, errors.WithStack(err)
	}

	return &user, nil
}
