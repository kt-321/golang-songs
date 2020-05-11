package model

import (
	"time"
)

//type Model struct {
//	ID        uint       `json:"id"`
//	CreatedAt time.Time  `json:"createdAt"`
//	UpdatedAt time.Time  `json:"updatedAt"`
//	DeletedAt *time.Time `json:"deletedAt"`
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
	Followings       []*User    `json:"followings" gorm:"many2many:user_follows;association_jointable_foreignkey:follow_id"`
	Bookmarkings     []*Song    `json:"bookmarkings" gorm:"many2many:bookmarks;"`
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
	Bookmarkers    []*User    `json:"bookmarkers" gorm:"many2many:bookmarks;"`
}

type UserFollow struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	UserID    uint       `json:"userId"`
	FollowID  uint       `json:"followId"`
}

type Bookmark struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	UserID    uint       `json:"userId"`
	SongID    uint       `json:"songId"`
}

type JWT struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}

type Response struct {
	Tracks Tracks `json:"tracks"`
}

type Tracks struct {
	Tracks struct {
		Href  string `json:"href"`
		Items []struct {
			Album struct {
				AlbumType string `json:"album_type"`
				Artists   []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href string `json:"href"`
					ID   string `json:"id"`
					Name string `json:"name"`
					Type string `json:"type"`
					URI  string `json:"uri"`
				} `json:"artists"`
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href   string `json:"href"`
				ID     string `json:"id"`
				Images []struct {
					Height int    `json:"height"`
					URL    string `json:"url"`
					Width  int    `json:"width"`
				} `json:"images"`
				Name                 string `json:"name"`
				ReleaseDate          string `json:"release_date"`
				ReleaseDatePrecision string `json:"release_date_precision"`
				TotalTracks          int    `json:"total_tracks"`
				Type                 string `json:"type"`
				URI                  string `json:"uri"`
			} `json:"album"`
			Artists []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
			DiscNumber  int  `json:"disc_number"`
			DurationMs  int  `json:"duration_ms"`
			Explicit    bool `json:"explicit"`
			ExternalIds struct {
				Isrc string `json:"isrc"`
			} `json:"external_ids"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href        string `json:"href"`
			ID          string `json:"id"`
			IsLocal     bool   `json:"is_local"`
			IsPlayable  bool   `json:"is_playable"`
			Name        string `json:"name"`
			Popularity  int    `json:"popularity"`
			PreviewURL  string `json:"preview_url"`
			TrackNumber int    `json:"track_number"`
			Type        string `json:"type"`
			URI         string `json:"uri"`
		} `json:"items"`
		Limit    int         `json:"limit"`
		Next     string      `json:"next"`
		Offset   int         `json:"offset"`
		Previous interface{} `json:"previous"`
		Total    int         `json:"total"`
	} `json:"tracks"`
}
