package interfaces

import (
	"encoding/json"
	"golang-songs/model"
	"net/http"
)

// レスポンスにエラーを突っ込んで返却するメソッド.
func errorInResponse(w http.ResponseWriter, status int, error model.Error) {
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(error); err != nil {
		var error model.Error
		error.Message = "リクエストボディのデコードに失敗しました。"
		errorInResponse(w, http.StatusInternalServerError, error)
		return
	}
}
