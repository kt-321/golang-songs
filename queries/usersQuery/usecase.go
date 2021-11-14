package usersQuery

import (
	"golang-songs/model"
)

type usecase struct {
	da dataAccessor
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

//func (ui *usecase) Update(userID int, p model.User) error {
//	return ui.da.Update(userID, p)
//}
