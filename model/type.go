package model

type JWT struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}

// Auth は署名前の認証トークン情報を表す.
type Auth struct {
	Email string
}

type Form struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Code struct {
	Code string
}

type ErrorSet struct {
	StatusCode int
	//MessageNumber int
	Message string
	Err     *error
}
