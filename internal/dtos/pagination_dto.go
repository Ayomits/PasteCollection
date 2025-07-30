package dtos

import "api/internal/enums"

type PaginationDto struct {
	StartFrom *int                   `json:"startFrom"`
	Order     *enums.PaginationOrder `json:"order" validate:"omitempty,oneof=prev next"`
	Limit     *int                   `json:"limit" validate:"omitempty,oneof=2"`
}
