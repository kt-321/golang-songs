package domain

import "errors"

type (
	UserID           int
	UserName         string
	UserAge          int
	Gender           int
	UserImageUrl     string
	FavoriteMusicAge int
	FavoriteArtist   string
	//TODO modelから変更
	//BookmarkingSongIDs []int
	//FollowingUserIDs []int

	SongID         int
	Title          string
	Artist         string
	MusicAge       int
	SongImage      string
	Video          string
	Album          string
	SpotifyTrackId string
	//TODO modelから変更
	//BookmarkerIDs []int
)

const (
	Male   Gender = 1
	Female Gender = 2
)

func parseUserID(userID int) (UserID, error) {
	if userID < 1 {
		return UserID(nil), errors.New("InvalidID")
	}
	return UserID(userID), nil
}

func parseUserName(userName string) (UserName, error) {
	if userName == "" {
		return "", errors.New("EmptyNotAllowed")
	}
	return UserName(userName), nil
}

func parseUserAge(userAge int) (UserAge, error) {
	if err := validateUserAge(userAge); err != nil {
		return UserAge(nil), err
	}
	return UserAge(userAge), nil
}

func validateUserAge(userAge int) error {
	validAges := [7]int{10, 20, 30, 40, 50, 60, 70}
	for _, a := range validAges {
		if a == userAge {
			return nil
		}
	}
	return errors.New("InvalidUserAge")
}

func parseGender(gender int) (Gender, error) {
	if gender < 1 || 2 < gender {
		return Gender(nil), errors.New("InvalidArgument")
	}
	return Gender(gender), nil
}

func parseUserImageUrl(userImageUrl string) (UserImageUrl, error) {
	if userImageUrl == "" {
		return "", errors.New("EmptyNotAllowed")
	}
	return UserImageUrl(userImageUrl), nil
}

func parseFavoriteMusicAge(favoriteMusicAge int) (FavoriteMusicAge, error) {
	if err := validateFavoriteMusicAge(favoriteMusicAge); err != nil {
		return FavoriteMusicAge(nil), err
	}
	return FavoriteMusicAge(favoriteMusicAge), nil
}

func validateFavoriteMusicAge(favoriteMusicAge int) error {
	validAges := [6]int{0, 1970, 1980, 1990, 2000, 2010}
	for _, a := range validAges {
		if a == favoriteMusicAge {
			return nil
		}
	}
	return errors.New("InvalidMusicAge")
}

func parseFavoriteArtist(favoriteArtist string) (FavoriteArtist, error) {
	if favoriteArtist == "" {
		return "", errors.New("EmptyNotAllowed")
	}
	return FavoriteArtist(favoriteArtist), nil
}

func parseSongID(songID int) (SongID, error) {
	if songID < 1 {
		return SongID(nil), errors.New("InvalidID")
	}
	return SongID(songID), nil
}

func parseTitle(title string) (Title, error) {
	if title == "" {
		return "", errors.New("EmptyNotAllowed")
	}
	return Title(title), nil
}

func parseArtist(artist string) (Artist, error) {
	if artist == "" {
		return "", errors.New("EmptyNotAllowed")
	}
	return Artist(artist), nil
}

func parseMusicAge(musicAge int) (MusicAge, error) {
	if err := validateFavoriteMusicAge(musicAge); err != nil {
		return MusicAge(nil), err
	}
	return MusicAge(musicAge), nil
}

//
func validateMusicAge(musicAge int) error {
	validAges := [6]int{0, 1970, 1980, 1990, 2000, 2010}
	for _, a := range validAges {
		if a == musicAge {
			return nil
		}
	}
	return errors.New("InvalidMusicAge")
}

func parseSongImage(songImage string) (SongImage, error) {
	if songImage == "" {
		return "", errors.New("EmptyNotAllowed")
	}
	return SongImage(songImage), nil
}

func parseVideo(video string) (Video, error) {
	if video == "" {
		return "", errors.New("EmptyNotAllowed")
	}
	return Video(video), nil
}

func parseAlbum(album string) (Album, error) {
	if album == "" {
		return "", errors.New("EmptyNotAllowed")
	}
	return Album(album), nil
}

func parseSpotifyTrackId(spotifyTrackId string) (SpotifyTrackId, error) {
	if spotifyTrackId == "" {
		return "", errors.New("EmptyNotAllowed")
	}
	return SpotifyTrackId(spotifyTrackId), nil
}
