package controllers

import (
	"api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	Register(c *fiber.Ctx) error
}

type authController struct {
	authService services.AuthService
}

func NewAuthController(s services.AuthService) AuthController {
	return &authController{authService: s}
}

func (c *authController) Register(ctx *fiber.Ctx) error {
	return c.authService.Register(ctx)
}
