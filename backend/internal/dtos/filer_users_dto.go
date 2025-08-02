package dtos

type UserFiltersDto struct {
	Id          *int    `json:"userId" validate:"omitempty,min=1"`
	Username    *string `json:"username" validate:"omitempty,min=1,max=32"`
	DisplayName *string `json:"displayName" validate:"omitempty,min=1,max=64"`
	SocialId    *string `json:"socialId" validate:"omitempty,min=1,max=255"`
	MatchAll    *bool   `json:"matchAll" validate:"omitempty"`
	Strict      *bool   `json:"strict" validate:"omitempty"`
}
