package model

import "time"

type Songs []Song

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
	Bookmarkers    []*User    `json:"bookmarkers" gorm:"many2many:bookmarks;""`
}
