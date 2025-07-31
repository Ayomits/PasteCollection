package services

import (
	"api/internal/dtos"
	"api/internal/repositories"
	"api/internal/responses"
	"api/internal/services/validators"

	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	Register(c *fiber.Ctx) error
}

type authService struct {
	userRepository *repositories.UserRepository
}

func NewAuthService(r *repositories.UserRepository) AuthService {
	return &authService{userRepository: r}
}

func (s *authService) Register(ctx *fiber.Ctx) error {
	var req dtos.RegisterUserDto

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Invalid body"))
	}

	violations := validators.AppValidatorInstance.Validate(req)

	if violations != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(violations)
	}

	return ctx.Status(fiber.StatusCreated).JSON(req)
}
