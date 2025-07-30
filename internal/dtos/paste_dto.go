package dtos

type PasteDto struct {
    Title string   `json:"title" validate:"required,min=1,max=32"`
    Tags  []string `json:"tags" validate:"max=3,dive,required,min=1,max=32"`
	Paste string `json:"paste" validate:"required,min=1,max=2096"`
}
