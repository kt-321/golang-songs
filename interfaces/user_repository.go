package interfaces

import (
	"golang-songs/domain"
	"time"
)

// A UserRepository belong to the inteface layer
type UserRepository struct {
	SQLHandler SQLHandler
}

// FindAll is returns the number of entities.
func (ur *UserRepository) FindAll() (users domain.Users, err error) {
	const query = `
		SELECT
			id,
			createdAt,
			updatedAt,
			deletedAt
			name,
			email,
			age,
			gender,
			imageUrl,
			favoriteMusicAge,
			favoriteArtist,
			comment
		FROM
			users
	`
	rows, err := ur.SQLHandler.Query(query)

	defer rows.Close()

	if err != nil {
		return
	}

	for rows.Next() {
		var id uint
		var createdAt time.Time
		var updatedAt time.Time
		var deletedAt *time.Time
		var name string
		var email string
		var age int
		var gender int
		var imageUrl string
		var favoriteMusicAge int
		var favoriteArtist string
		var comment string

		if err = rows.Scan(&id, &createdAt, &updatedAt, &deletedAt, &name, &email, &age, &gender, &imageUrl, &favoriteMusicAge, &favoriteArtist, &comment); err != nil {
			return
		}
		user := domain.User{
			ID:               id,
			CreatedAt:        createdAt,
			UpdatedAt:        updatedAt,
			DeletedAt:        deletedAt,
			Name:             name,
			Email:            email,
			Age:              age,
			Gender:           gender,
			ImageUrl:         imageUrl,
			FavoriteMusicAge: favoriteMusicAge,
			FavoriteArtist:   favoriteArtist,
			Comment:          comment,
			//Password:         password,
			//Bookmarkings:     bookmarkings,
			//Followings:       followings,
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}

// FindByID is returns the entity identified by the given id.
func (ur *UserRepository) FindByID(userID int) (user domain.User, err error) {
	const query = `
		SELECT
			id,
			createdAt,
			updatedAt,
			deletedAt
			name,
			email,
			age,
			gender,
			imageUrl,
			favoriteMusicAge,
			favoriteArtist,
			comment
		FROM
			users
		WHERE
			id = ?
	`
	row, err := ur.SQLHandler.Query(query, userID)

	defer row.Close()

	if err != nil {
		return
	}

	var id uint
	var createdAt time.Time
	var updatedAt time.Time
	var deletedAt *time.Time
	var name string
	var email string
	var age int
	var gender int
	var imageUrl string
	var favoriteMusicAge int
	var favoriteArtist string
	var comment string

	row.Next()
	if err = row.Scan(&id, &createdAt, &updatedAt, &deletedAt, &name, &email, &age, &gender, &imageUrl, &favoriteMusicAge, &favoriteArtist, &comment); err != nil {
		return
	}
	user.ID = id
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt
	user.DeletedAt = deletedAt
	user.Name = name
	user.Email = email
	user.Age = age
	user.Gender = gender
	user.ImageUrl = imageUrl
	user.FavoriteMusicAge = favoriteMusicAge
	user.FavoriteArtist = favoriteArtist
	user.Comment = comment
	return
}
