package domain

import "time"

// A Songs belong to the domain layer.
type Songs []Song

// A Post belong to the domain layer.
//type Song struct {
//	ID     int    `json:"id"`
//	UserID int    `json:"user_id"`
//	Body   string `json:"body"`
//}

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
