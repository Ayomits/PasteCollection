package controllers

import (
	"api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	FindById(c *fiber.Ctx) error
	FindByUsername(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type userController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &userController{userService: userService}
}

func (u *userController) Create(c *fiber.Ctx) error {
	return u.userService.CreateUser(c)
}

func (u *userController) Delete(c *fiber.Ctx) error {
	return u.userService.DeleteUser(c)
}

func (u *userController) FindById(c *fiber.Ctx) error {
	return u.userService.FindUser(c)
}

func (u *userController) FindByUsername(c *fiber.Ctx) error {
	return u.userService.FindUser(c)
}

func (u *userController) Update(c *fiber.Ctx) error {
	return u.userService.UpdateUser(c)
}
