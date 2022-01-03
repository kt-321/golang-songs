package userQuery

import (
	"golang-songs/model"
)

type usecase struct {
	da DataAccessor
}

type Usecase interface {
	Index() (*[]model.User, error)
	User(string) (*model.User, error)
	Show(int) (*model.User, error)
}

type DataAccessor interface {
	FindAll() (*[]model.User, error)
	GetUser(string) (*model.User, error)
	FindByID(int) (*model.User, error)
}

func (ui *usecase) Index() (*[]model.User, error) {
	return ui.da.FindAll()
}

func (ui *usecase) User(userEmail string) (*model.User, error) {
	return ui.da.GetUser(userEmail)
}

func (ui *usecase) Show(userID int) (*model.User, error) {
	return ui.da.FindByID(userID)
}
