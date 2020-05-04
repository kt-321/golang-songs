package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name             string `json:"name"`
	Email            string `json:"email,omitempty"`
	Age              int    `json:"age,omitempty"`
	Gender           int    `json:"gender,omitempty"`
	ImageUrl         string `json:"imageUrl,omitempty"`
	FavoriteMusicAge int    `json:"favoriteMusicAge,omitempty"`
	FavoriteArtist   string `json:"favoriteArtist,omitempty"`
	Comment          string `json:"comment,omitempty"`
	Password         string `json:"-"`
}

type JWT struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}
