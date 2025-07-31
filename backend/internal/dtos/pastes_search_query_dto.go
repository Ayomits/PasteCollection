package dtos

type PastesSearchQueryDto struct {
	Filter     *PastesFilterDto `json:"filter"`
	Pagination *PaginationDto   `json:"pagination"`
}
