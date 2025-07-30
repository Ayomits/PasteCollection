package dtos

type FindByTitleQueryDto struct {
	Strict bool `json:"strict"`
}

func NewFindByTitleQueryDto(strict bool) *FindByTitleQueryDto {
	return &FindByTitleQueryDto{
		Strict: strict,
	}
}
