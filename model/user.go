package model

import "time"

type Users []User

//TODO 調整
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
	Bookmarkings     []*Song    `json:"bookmarkings" gorm:"many2many:bookmarks;"`
	Followings       []*User    `json:"followings" gorm:"many2many:user_follows;association_jointable_foreignkey:follow_id"`
}
