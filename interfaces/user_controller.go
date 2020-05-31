package interfaces

import (
	"encoding/json"
	"golang-songs/usecases"
	"net/http"
	"strconv"
)

// A UserController belong to the interface layer.
type UserController struct {
	UserInteractor usecases.UserInteractor
	Logger         usecases.Logger
}

// NewUserController returns the resource of users.
func NewUserController(sqlHandler SQLHandler, logger usecases.Logger) *UserController {
	return &UserController{
		UserInteractor: usecases.UserInteractor{
			UserRepository: &UserRepository{
				SQLHandler: sqlHandler,
			},
		},
		Logger: logger,
	}
}

// Index return response which contain a listing of the resource of users.
func (uc *UserController) Index(w http.ResponseWriter, r *http.Request) {
	uc.Logger.LogAccess("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	//users, err := uc.UserInteractor.Index()

	allUsers, err := uc.UserInteractor.Index()

	//以下は元のコード
	v, err := json.Marshal(allUsers)
	if err != nil {
		//var error model.Error
		//error.Message = "ユーザー一覧の取得に失敗しました"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
	if _, err := w.Write(v); err != nil {
		//var error model.Error
		//error.Message = "ユーザー一覧の取得に失敗しました。"
		//errorInResponse(w, http.StatusInternalServerError, error)
		return
	}

	//↓以下は、とって来たコード
	//if err != nil {
	//	uc.Logger.LogError("%s", err)
	//
	//	w.Header().Set("Content-Type", "application/json")
	//	w.WriteHeader(500)
	//	json.NewEncoder(w).Encode(err)
	//}
	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(users)
}

// Show return response which contain the specified resource of a user.
func (uc *UserController) Show(w http.ResponseWriter, r *http.Request) {
	uc.Logger.LogAccess("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	userID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	user, err := uc.UserInteractor.Show(userID)
	if err != nil {
		uc.Logger.LogError("%s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
