package dtos

type UserDto struct {
	Username string `json:"username" validate:"required,min=1,max=32"`
}
