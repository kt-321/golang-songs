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

type JWT struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}
