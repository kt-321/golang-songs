package interfaces

//
//import (
//	"golang-songs/usecases"
//
//	"github.com/jinzhu/gorm"
//)
//
//type UserController struct {
//	UserInteractor usecases.UserInteractor
//}
//
//func NewUserController(DB *gorm.DB) *UserController {
//	return &UserController{
//		UserInteractor: usecases.UserInteractor{
//			UserRepository: &UserRepository{
//				DB: DB,
//			},
//		},
//	}
//}

// 全てのユーザーを返す.
//func (uc *UserController) AllUsersHandler(w http.ResponseWriter, r *http.Request) {
//	allUsers, err := uc.UserInteractor.Index()
//
//	if err != nil {
//		ErrorInResponse(w, http.StatusInternalServerError, GetUserError)
//
//		return
//	}
//
//	v, err := json.Marshal(allUsers)
//
//	if err != nil {
//		ErrorInResponse(w, http.StatusInternalServerError, JsonEncodeError)
//
//		return
//	}
//
//	if _, err := w.Write(v); err != nil {
//		ErrorInResponse(w, http.StatusInternalServerError, GetUsersListError)
//
//		return
//	}
//}
//
//// リクエストユーザーの情報を返す.
//func (uc *UserController) UserHandler(w http.ResponseWriter, r *http.Request) {
//	// リクエストユーザーのメアドを取得.
//	userEmail, errorSet := GetEmail(r)
//
//	if errorSet != nil {
//		ErrorInResponse(w, errorSet.StatusCode, errorSet.Message)
//
//		return
//	}
//
//	var user *model.User
//
//	user, err := uc.UserInteractor.User(userEmail)
//
//	if err != nil {
//		ErrorInResponse(w, http.StatusInternalServerError, GetAccountError)
//
//		return
//	}
//
//	v, err := json.Marshal(user)
//
//	if err != nil {
//		ErrorInResponse(w, http.StatusInternalServerError, JsonEncodeError)
//
//		return
//	}
//
//	if _, err := w.Write(v); err != nil {
//		ErrorInResponse(w, http.StatusInternalServerError, GetUserDetailError)
//
//		return
//	}
//}
//
//// idで指定したユーザーの情報を返す.
//func (uc *UserController) GetUserHandler(w http.ResponseWriter, r *http.Request) {
//	// 対象のユーザーidを取得.
//	userID, errorSet := GetId(r)
//
//	if errorSet != nil {
//		ErrorInResponse(w, errorSet.StatusCode, errorSet.Message)
//
//		return
//	}
//
//	user, err := uc.UserInteractor.Show(userID)
//
//	if err != nil {
//		ErrorInResponse(w, http.StatusInternalServerError, GetAccountError)
//
//		return
//	}
//
//	v, err := json.Marshal(user)
//
//	if err != nil {
//		ErrorInResponse(w, http.StatusInternalServerError, JsonEncodeError)
//
//		return
//	}
//
//	if _, err := w.Write(v); err != nil {
//		ErrorInResponse(w, http.StatusInternalServerError, GetUserDetailError)
//
//		return
//	}
//}

// idで指定したユーザーの情報を更新する.
//func (uc *UserController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
//	// 対象のユーザーidを取得.
//	userID, errorSet := GetId(r)
//
//	if errorSet != nil {
//		ErrorInResponse(w, errorSet.StatusCode, errorSet.Message)
//
//		return
//	}
//
//	var d model.User
//
//	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
//		ErrorInResponse(w, http.StatusInternalServerError, DecodeError)
//
//		return
//	}
//
//	if err := uc.UserInteractor.Update(userID, d); err != nil {
//		ErrorInResponse(w, http.StatusInternalServerError, UpdateUserError)
//
//		return
//	}
//
//	// 204 No Content.
//	w.WriteHeader(204)
//}
