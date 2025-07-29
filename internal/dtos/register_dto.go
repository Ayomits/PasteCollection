package dtos

type UserDto struct {
	username string `validate"required,min=1,max=32"`
}
