package services

import (
	"api/internal/dtos"
	"api/internal/formatter"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	Register(c *fiber.Ctx) error
}

type authService struct {
	configService       ConfigService
	validationFormatter formatter.ValidationFormatter
}

func NewAuthService(configService ConfigService) AuthService {
	v := validator.New()
	vf := formatter.NewValidationFormatter(v)
	dtos.RegisterCustomValidations(v)
	return &authService{configService: configService, validationFormatter: vf}
}

func (s *authService) Register(ctx *fiber.Ctx) error {
	var req dtos.RegisterUserDto

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Invalid body",
			"violations": err,
		})
	}

	err := s.validationFormatter.Validate(req)

	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(err)
	}

	return ctx.Status(201).JSON(fiber.Map{
		"success": true,
	})
}
