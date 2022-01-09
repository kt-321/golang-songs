package userQuery

import (
	"golang-songs/model"
	"time"
)

type usecase struct {
	da DataAccessor
}

type Usecase interface {
	GetAllUsers() (*[]model.User, error)
	FindUserByEmail(userEmail string) (*findUserByEmailRes, error)
	//GetUserByEail(string) (*model.User, error)
	FindUserByID(userID int) (*findUserByIDRes, error)
}

type findUserByEmailRes struct {
	ID               uint
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
	Name             string
	Email            string
	Age              int
	Gender           int
	ImageUrl         string
	FavoriteMusicAge int
	FavoriteArtist   string
	Comment          string
	//Password         string
	Bookmarkings []*model.Song
	Followings   []*model.User
}

type findUserByIDRes struct {
	ID               uint
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
	Name             string
	Email            string
	Age              int
	Gender           int
	ImageUrl         string
	FavoriteMusicAge int
	FavoriteArtist   string
	Comment          string
	//Password         string
	Bookmarkings []*model.Song
	Followings   []*model.User
}

type DataAccessor interface {
	GetAllUsers() (*[]model.User, error)
	GetUserInfoByEmail(userEmail string) (*getUserInfoByEmailRes, error)
	GetUserInfoByID(userID int) (*model.User, error)
	GetBookmarkings(userID int) (*getBookmarkingsRes, error)
	GetFollowees(userID int) (*getFolloweesRes, error)
}

type getUserInfoByEmailRes struct {
	ID               uint
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
	Name             string
	Email            string
	Age              int
	Gender           int
	ImageUrl         string
	FavoriteMusicAge int
	FavoriteArtist   string
	Comment          string
	Password         string
	//Bookmarkings     []*model.Song
	//Followings       []*model.User
}

type getBookmarkingsRes struct {
	Bookmarkings []*model.Song
}

type getFolloweesRes struct {
	Followees []*model.User
}

func (ui *usecase) GetAllUsers() (*[]model.User, error) {
	return ui.da.GetAllUsers()
}

func (ui *usecase) FindUserByEmail(userEmail string) (*findUserByEmailRes, error) {
	authUser, err := ui.da.GetUserInfoByEmail(userEmail)
	if err != nil {
		return nil, err
	}

	userID := int(authUser.ID)

	bookmarkingRes, err := ui.da.GetBookmarkings(int(userID))
	if err != nil {
		return nil, err
	}

	followeeRes, err := ui.da.GetFollowees(userID)
	if err != nil {
		return nil, err
	}

	return &findUserByEmailRes{
		ID:               authUser.ID,
		CreatedAt:        authUser.CreatedAt,
		UpdatedAt:        authUser.UpdatedAt,
		DeletedAt:        authUser.DeletedAt,
		Name:             authUser.Name,
		Email:            authUser.Email,
		Age:              authUser.Age,
		Gender:           authUser.Gender,
		ImageUrl:         authUser.ImageUrl,
		FavoriteMusicAge: authUser.FavoriteMusicAge,
		FavoriteArtist:   authUser.FavoriteArtist,
		Comment:          authUser.Comment,
		//Password         string
		Bookmarkings: bookmarkingRes.Bookmarkings,
		Followings:   followeeRes.Followees,
	}, nil
}

func (ui *usecase) FindUserByID(userID int) (*findUserByIDRes, error) {
	authUser, err := ui.da.GetUserInfoByID(userID)
	if err != nil {
		return nil, err
	}

	bookmarkingRes, err := ui.da.GetBookmarkings(int(userID))
	if err != nil {
		return nil, err
	}

	followeeRes, err := ui.da.GetFollowees(userID)
	if err != nil {
		return nil, err
	}

	return &findUserByIDRes{
		ID:               authUser.ID,
		CreatedAt:        authUser.CreatedAt,
		UpdatedAt:        authUser.UpdatedAt,
		DeletedAt:        authUser.DeletedAt,
		Name:             authUser.Name,
		Email:            authUser.Email,
		Age:              authUser.Age,
		Gender:           authUser.Gender,
		ImageUrl:         authUser.ImageUrl,
		FavoriteMusicAge: authUser.FavoriteMusicAge,
		FavoriteArtist:   authUser.FavoriteArtist,
		Comment:          authUser.Comment,
		//Password         string
		Bookmarkings: bookmarkingRes.Bookmarkings,
		Followings:   followeeRes.Followees,
	}, nil
}
