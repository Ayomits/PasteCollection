package dtos

type PastesFilterDto struct {
	Search *string   `json:"search" validate:"omitempty"`
	Strict *bool     `json:"strict" validate:"omitempty,oneof=true false"`
	Tags   *[]string `json:"tags" validate:"omitempty,dive,min=1,max=32"`
}
