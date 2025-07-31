package dtos

type UserFiltersDto struct {
	Id          *int    `json:"userId" validate:"omitempty,min=1"`
	Username    *string `json:"username" validate:"omitempty,min=1,max=32"`
	DisplayName *string `json:"displayName" validate:"omitempty,min=1,max=64"`
	SocialId    *string `json:"socialId" validate:"omitempty,min=1,max=255"`
	MatchAll    *bool   `json:"matchAll" validate:"omitempty"`
	Strict      *bool   `json:"strict" validate:"omitempty"`
}

func NewUserFiltersDto(username *string, displayName *string, socialId *string, matchAll *bool, strict *bool, id *int) *UserFiltersDto {
	return &UserFiltersDto{
		Username:    username,
		DisplayName: displayName,
		SocialId:    socialId,
		MatchAll:    matchAll,
		Strict:      strict,
		Id:          id,
	}
}
