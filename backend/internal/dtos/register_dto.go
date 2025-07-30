package dtos

import (
	"api/internal/enums"
)

type RegisterUserDto struct {
	Username string         `json:"username" validate:"required,min=1,max=32"`
	SocialId string         `json:"socialId" validate:"required"`
	AuthType enums.AuthType `json:"authType" validate:"required,oneof=discord telegram"`
}
