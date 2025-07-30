package dtos

type PastesSearchQueryDto struct {
	Filter     *FilterPastesDto `json:"filter"`
	Pagination *PaginationDto   `json:"pagination"`
}
