package model

import "time"

type User struct {
	ID               uint       `json:"id"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	DeletedAt        *time.Time `json:"deletedAt"`
	Name             string     `json:"name"`
	Email            string     `json:"email"`
	Age              int        `json:"age"`
	Gender           int        `json:"gender"`
	ImageUrl         string     `json:"imageUrl"`
	FavoriteMusicAge int        `json:"favoriteMusicAge"`
	FavoriteArtist   string     `json:"favoriteArtist"`
	Comment          string     `json:"comment"`
	Password         string     `json:"-"`
}

type Song struct {
	ID             uint       `json:"id"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	DeletedAt      *time.Time `json:"deletedAt"`
	Title          string     `json:"title"`
	Artist         string     `json:"artist"`
	MusicAge       int        `json:"musicAge"`
	Image          string     `json:"image"`
	Video          string     `json:"video"`
	Album          string     `json:"album"`
	Description    string     `json:"description"`
	SpotifyTrackId string     `json:"spotifyTrackId"`
	UserID         uint       `json:"userId"`
}

type JWT struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}

// Auth は署名前の認証トークン情報を表す。
type Auth struct {
	Email string
}
