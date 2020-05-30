package domain

import "time"

// A Users belong to the domain layer.
type Users []User

// A User belong to the domain layer.
//type User struct {
//	ID   int    `json:"id"`
//	Name string `json:"name"`
//}

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
