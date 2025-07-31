package dtos

import (
	"api/internal/enums"
)

type PaginationDto struct {
	StartFrom *int                   `json:"startFrom" validate:"omitempty,min=1"`
	Order     *enums.PaginationOrder `json:"order" validate:"omitempty,oneof=prev next"`
	Limit     *int                   `json:"limit" validate:"omitempty,oneof=5 10 15 20 25 30 35 40 45 50"`
	Sort      *string                `json:"sort" validate:"omitempty,oneof=ASC DESC"`
}
