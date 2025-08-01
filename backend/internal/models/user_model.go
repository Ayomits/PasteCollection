package models

import "time"

type UserModel struct {
	Id          int    `json:"id" validate:"omitempty"`
	Username    string `json:"username" validate:"omitempty"`
	DisplayName string `json:"displayName" validate:"omitempty"`
	SocialId    string `json:"socialId" validate:"omitempty"`

	CreatedAt time.Time `json:"created_at" validate:"omitempty"`
	UpdatedAt time.Time `json:"updated_at" validate:"omitempty"`
}
