package dtos

type FilterPastesDto struct {
	Search *string   `json:"search" validate:"omitempty"`
	Strict *bool     `json:"strict" validate:"omitempty,oneof=true false"`
	Tags   *[]string `json:"tags" validate:"omitempty"`
}
