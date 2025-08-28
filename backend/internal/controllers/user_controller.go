package controllers

import (
	"api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	FindUser(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
}

type userController struct {
	userService services.UserService
}

func NewUserController(s services.UserService) UserController {
	return &userController{userService: s}
}

// FindUser godoc
//
//	@Summary		Найти пользователя по фильтру
//	@Description	Найдёт пользователя по фильтру или вернёт 404
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id			query		int		false	"ID пользователя"
//	@Param			username	query		string	false	"Уникальное имя"
//	@Param			displayName	query		string	false	"Отображаемое имя"
//	@Param			socialId	query		string	false	"Айди в соц. сети"
//	@Param			matchAll	query		bool	false	"Если в true, тогда будет условие AND"	default(false)
//	@Param			strict		query		bool	false	"Если в true, тогда строгое сравнение"	default(false)
//	@Success		200			{object}	models.UserModel
//	@Failure		401			{object}	responses.UnauthorizedError
//	@Failure		404			{object}	responses.NotFoundError
//	@Failure		422			{object}	responses.ValidationError
//	@Failure		429			{object}	responses.RateLimitError
//	@Failure		500			{object}	responses.InternalError
//	@Router			/api/users [get]
//	@Security		ApiKeyAuth
func (u *userController) FindUser(c *fiber.Ctx) error {
	return u.userService.Find(c)
}

// CreateUser godoc
//
//	@Summary		Создать пользователя
//	@Description	Создаст пользователя или вернёт 422
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dtos.CreateUserDto	true	"Данные юзера"
//	@Success		200		{object}	models.UserModel
//	@Failure		401		{object}	responses.UnauthorizedError
//	@Failure		422		{object}	responses.ValidationError
//	@Failure		429		{object}	responses.RateLimitError
//	@Failure		500		{object}	responses.InternalError
//	@Router			/api/users [post]
//	@Security		ApiKeyAuth
func (u *userController) CreateUser(c *fiber.Ctx) error {
	return u.userService.Create(c)
}

// UpdateUser godoc
//
//	@Summary		Обновит пользователя
//	@Description	Обновит пользователя или вернёт 404, 422
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id			query		int					false	"ID пользователя"
//	@Param			username	query		string				false	"Уникальное имя"
//	@Param			displayName	query		string				false	"Отображаемое имя"
//	@Param			socialId	query		string				false	"Айди в соц. сети"
//	@Param			matchAll	query		bool				false	"Если в true, тогда будет условие AND"	default(false)
//	@Param			strict		query		bool				false	"Если в true, тогда строгое сравнение"	default(false)
//	@Param			request		body		dtos.UpdateUserDto	true	"Данные юзера"
//	@Success		200			{object}	models.UserModel
//	@Failure		401			{object}	responses.UnauthorizedError
//	@Failure		404			{object}	responses.NotFoundError
//	@Failure		422			{object}	responses.ValidationError
//	@Failure		429			{object}	responses.RateLimitError
//	@Failure		500			{object}	responses.InternalError
//	@Router			/api/users [put]
//	@Security		ApiKeyAuth
func (u *userController) UpdateUser(c *fiber.Ctx) error {
	return u.userService.Update(c)
}

// UpdateUser godoc
//
//	@Summary		Удалить пользователя
//	@Description	Удалит пользователя или вернёт 404
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id			query		int		false	"ID пользователя"
//	@Param			username	query		string	false	"Уникальное имя"
//	@Param			displayName	query		string	false	"Отображаемое имя"
//	@Param			socialId	query		string	false	"Айди в соц. сети"
//	@Param			matchAll	query		bool	false	"Если в true, тогда будет условие AND"	default(false)
//	@Param			strict		query		bool	false	"Если в true, тогда строгое сравнение"	default(false)
//	@Success		200			{object}	models.UserModel
//	@Failure		401			{object}	responses.UnauthorizedError
//	@Failure		404			{object}	responses.NotFoundError
//	@Failure		422			{object}	responses.ValidationError
//	@Failure		429			{object}	responses.RateLimitError
//	@Failure		500			{object}	responses.InternalError
//	@Router			/api/users [delete]
//	@Security		ApiKeyAuth
func (u *userController) DeleteUser(c *fiber.Ctx) error {
	return u.userService.Delete(c)
}
