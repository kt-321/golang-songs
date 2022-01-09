package userQuery

import (
	"golang-songs/model"

	"github.com/jmoiron/sqlx"
)

type dataAccessor struct {
	DB *sqlx.DB
}

func (ur *dataAccessor) GetAllUsers() (*[]model.User, error) {
	var users []model.User

	rows, err := ur.DB.Queryx("SELCT * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		err := rows.StructScan(&user)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return &users, nil
}

func (ur *dataAccessor) GetUserInfoByEmail(userEmail string) (*getUserInfoByEmailRes, error) {
	q := `
	select *
	from users
	where email = ?
`

	res := &getUserInfoByEmailRes{}

	rows, err := ur.DB.Queryx(q, userEmail)
	if err != nil {
		return nil, err
	}

	if err := rows.StructScan(res); err != nil {
		return nil, err
	}

	return res, nil
}

func (ur *dataAccessor) GetUserInfoByID(userID int) (*model.User, error) {
	var user model.User

	q := `
	select *
	from users
	where id = ?
`

	res := &getUserInfoByEmailRes{}

	rows, err := ur.DB.Queryx(q, userID)
	if err != nil {
		return nil, err
	}

	if err := rows.StructScan(res); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *dataAccessor) GetFollowees(userID int) (*getFolloweesRes, error) {
	q := `
	select *
	from users u1
	inner join user_follows uf on uf.user_id = u1.id
	inner join users u2 on u2.id = uf.follow_id
	where id = ?
`

	rows, err := ur.DB.Queryx(q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := &getFolloweesRes{}

	var followees []model.User
	for rows.Next() {
		var followee model.User
		err := rows.StructScan(&followee)
		if err != nil {
			return nil, err
		}

		followees = append(followees, followee)
	}

	return res, nil
}

func (ur *dataAccessor) GetBookmarkings(userID int) (*getBookmarkingsRes, error) {
	q := `
	select *
	from users u
	inner join bookmarks b on b.user_id = u.id
	inner join songs s on s.id = b.song_id
	where id = ?
`

	rows, err := ur.DB.Queryx(q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := &getBookmarkingsRes{}

	var boolmarkings []model.Song
	for rows.Next() {
		var bookmarking model.Song
		err := rows.StructScan(&bookmarking)
		if err != nil {
			return nil, err
		}

		boolmarkings = append(boolmarkings, bookmarking)
	}

	return res, nil
}
