package dtos

type UserDto struct {
	Username    string `json:"username" validate:"required,min=1,max=32"`
	DisplayName string `json:"displayName" validate:"required,min=1,max=64"`
	SocialId    string `json:"socialId" validate:"required,min=1,max=255"`
}

type UpdateUserDto struct {
	Username    string `json:"username" validate:"required,min=1,max=32"`
	DisplayName string `json:"displayName" validate:"required,min=1,max=64"`
}
