package usecases

import (
	"golang-songs/model"
)

//"golang-songs/domain"

// A UserInteractor belong to the usecases layer.
type UserInteractor struct {
	UserRepository UserRepository
}

// Index is display a listing of the resource.
//func (ui *UserInteractor) Index() (users domain.Users, err error) {
func (ui *UserInteractor) Index() (*model.Users, error) {
	//	func (ui *UserInteractor) Index() *model.Users {
	//users := []model.User{}
	//users, err := ui.UserRepository.FindAll()
	users, err = ui.UserRepository.FindAll()
	return
	//if err != nil {
	//	var error model.Error
	//	error.Message = "JSONへの変換に失敗しました"
	//	errorInResponse(w, http.StatusInternalServerError, error)
	//	return
	//}
}

// Show is display the specified resource.
//func (ui *UserInteractor) Show(userID int) (user domain.User, err error) {
func (ui *UserInteractor) Show(userID int) (model.User, error) {
	user, err = ui.UserRepository.FindByID(userID)

	return
}
