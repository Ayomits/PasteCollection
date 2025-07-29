package controllers

import (
	"api/internal/services"
	"fmt"

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
	body := ctx.BodyParser(ctx.Body())
	fmt.Println(body)
	return ctx.JSON(map[string]int{"success": 10})
}

func (c *authController) Update(ctx *fiber.Ctx) error {
	return ctx.JSON(map[string]bool{"success": true})
}
