package responses

type PaginationResponse[T any] struct {
	Items   *[]T `json:"items"`
	HasNext bool  `json:"hasNext"`
}

func NewPaginationResponse[T any](items *[]T, hasNext bool) *PaginationResponse[T] {
	return &PaginationResponse[T]{
		Items:   items,
		HasNext: hasNext,
	}
}
