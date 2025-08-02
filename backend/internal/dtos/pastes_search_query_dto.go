package dtos

type PastesSearchQueryDto struct {
	Filter     *PastesFilterDto `json:"filter" validate:"omitempty"`
	Pagination *PaginationDto   `json:"pagination" validate:"omitempty"`
}
