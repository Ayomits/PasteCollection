package models

import "time"

type PasteModel struct {
	Id    int    `db:"id" json:"id" validate:"omitempty"`
	Title string `db:"title" json:"title" validate:"omitempty"`
	Paste string `db:"paste" json:"paste" validate:"omitempty"`

	UserId int `db:"user_id" json:"userId" validate:"omitempty"`

	CreatedAt time.Time `db:"created_at" json:"createdAt" validate:"omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt" validate:"omitempty"`
}
