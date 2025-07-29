package dtos

import (
	"api/internal/enums"
	"github.com/go-playground/validator/v10"
)

type RegisterUserDto struct {
	Username string         `json:"username" validate:"required,min=1,max=32"`
	SocialId string         `json:"socialId" validate:"required"`
	AuthType enums.AuthType `json:"authType" validate:"required,auth_type"`
}

const (
	AuthTypeValidation = "auth_type"
)

func RegisterCustomValidations(v *validator.Validate) {
	_ = v.RegisterValidation("auth_type", func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(enums.AuthType)
		if !ok {
			return false
		}

		for _, valid := range value.Values() {
			if value == valid {
				return true
			}
		}
		return false
	})
}
