package services

import (
	"api/internal/dtos"
	"api/internal/responses"
	"api/internal/services/validators"

	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	Register(c *fiber.Ctx) error
}

type authService struct {
	configService ConfigService
	v             validators.AppValidator
}

func NewAuthService(configService ConfigService) AuthService {
	return &authService{configService: configService, v: validators.AppValidatorInstance}
}

func (s *authService) Register(ctx *fiber.Ctx) error {
	var req dtos.RegisterUserDto

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Invalid body"))
	}

	err := s.v.Validate(req)

	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(req)
}
