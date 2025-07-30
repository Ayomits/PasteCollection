package models

import "time"

type PasteModel struct {
	Id         int       `json:"id"`
	Title      string    `json:"title"`
	Tags       []string  `json:"tags"`
	Paste      string    `json:"paste"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
