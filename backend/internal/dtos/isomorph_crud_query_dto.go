package dtos

type IsomorphQueryDto struct {
	AsUsername bool `json:"asUsername" validate:"omitempty,oneof=true false"`
}
