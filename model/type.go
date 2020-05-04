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
	Password         string `json:"password"`
}

//type UserInResponse struct {
//	//gorm.Model
//	Name             string `json:"name"`
//	Email            string `json:"email,omitempty"`
//	Age              int    `json:"age,omitempty"`
//	Gender           int    `json:"gender,omitempty"`
//	ImageUrl         string `json:"age,omitempty"`
//	FavoriteMusicAge int    `json:"favoriteMusicAge,omitempty"`
//	FavoriteArtist   string `json:"favoriteArtist,omitempty"`
//	Comment          string `json:"comment,omitempty"`
//	//Password         string `json:"password"`
//}

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
