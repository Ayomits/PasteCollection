package dtos

type PastesFilterDto struct {
	Search   *string `json:"search" validate:"omitempty"`
	Strict   *bool   `json:"strict" validate:"omitempty,oneof=true false"`
	UserId   *int    `json:"userId" validate:"omitempty,min=1"`
	SocialId *string `json:"socialId" validate:"omitempty"`
	PasteId  *int    `json:"pasteId" validate:"omitempty"`
}
