package model

import "time"

type UserFollow struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	UserID    uint       `json:"userId"`
	FollowID  uint       `json:"followId"`
}
