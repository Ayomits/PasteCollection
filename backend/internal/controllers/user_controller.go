package controllers

import (
	"api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	Find(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type userController struct {
	userService services.UserService
}

func NewUserController(s services.UserService) UserController {
	return &userController{userService: s}
}

func (u *userController) Find(c *fiber.Ctx) error {
	return u.userService.Find(c)
}

func (u *userController) Create(c *fiber.Ctx) error {
	return u.userService.Create(c)
}

func (u *userController) Update(c *fiber.Ctx) error {
	return u.userService.Update(c)
}

func (u *userController) Delete(c *fiber.Ctx) error {
	return u.userService.Delete(c)
}
