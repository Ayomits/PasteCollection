package dtos

type PasteDto struct {
	Title  string `json:"title" validate:"required,min=1,max=32"`
	Paste  string `json:"paste" validate:"required,min=1,max=2096"`
	UserId int    `json:"userId" validate:"required,min=1"`
}

type UpdatePasteDto struct {
	Title string `json:"title" validate:"required,min=1,max=32"`
	Paste string `json:"paste" validate:"required,min=1,max=2096"`
}
