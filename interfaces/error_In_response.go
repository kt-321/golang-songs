package interfaces

import (
	"encoding/json"
	"golang-songs/model"
	"log"
	"net/http"
)

// レスポンスにエラーメッセージを突っ込んで返却するメソッド.
func ErrorInResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)

	log.Printf("%v: %v", status, message)

	error2 := model.Error{Message: message}

	if err := json.NewEncoder(w).Encode(error2); err != nil {
		log.Println("構造体errorのJSONエンコードに失敗しました")

		return
	}
}

const (
	DecodeError                     = "リクエストボディのデコードに失敗しました"
	RequiredEmailError              = "Emailは必須です"
	RequiredPasswordError           = "パスワードは必須です"
	InvalidPasswordError            = "パスワードの値が不正です"
	CreateAccountError              = "アカウントの作成に失敗しました"
	JsonEncodeError                 = "JSONへの変換に失敗しました"
	GetUserDetailError              = "ユーザー情報の取得に失敗しました"
	CreateTokenError                = "トークンの作成に失敗しました"
	GetJwtTokenError                = "JWTトークンの取得に失敗しました"
	GetSongError                    = "曲が見つかりません"
	GetSongsListError               = "曲一覧の取得に失敗しました"
	GetIdError                      = "パスパラメータからidを取得できません"
	ConvertIdToIntError             = "idのint型への型変換に失敗しました"
	GetSongDetailError              = "曲情報の取得に失敗しました"
	ParseAuthenticationCodeError    = "認証コードのパースに失敗しました"
	PostSongError                   = "曲の投稿に失敗しました"
	UpdateSongError                 = "曲の更新に失敗しました"
	DeleteSongError                 = "曲の削除に失敗しました"
	GetUserError                    = "ユーザーが見つかりません"
	GetUsersListError               = "ユーザー一覧の取得に失敗しました"
	GetAccountError                 = "該当するアカウントが見つかりません"
	UpdateUserError                 = "ユーザー情報の更新に失敗しました"
	GetAuthenticationTokenError     = "認証トークンの取得に失敗しました"
	GetBearerTokenError             = "bearerトークンの取得に失敗しました"
	FollowUserError                 = "ユーザーのフォローに失敗しました"
	UnfollowUserError               = "ユーザーのフォロー解除に失敗しました"
	BookmarkSongError               = "曲のお気に入り登録に失敗しました"
	RemoveBookmarkError             = "曲のお気に入り解除に失敗しました"
	GetSpotifyTokenError            = "Spotifyトークンの取得に失敗しました"
	GetUrlError                     = "URLの取得に失敗しました"
	GetSpotifyTokenFromReqBodyError = "リクエストボディからSpotifyトークンの取得に失敗しました"
	GetTraksError                   = "トラックの取得に失敗しました"
)
