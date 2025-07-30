package services

import (
	"api/internal/dtos"
	"api/internal/models"
	"api/internal/repositories"
	"api/internal/responses"
	"api/internal/services/querymap"
	"api/internal/services/validators"
	"api/internal/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	FindUser(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(r repositories.UserRepository) UserService {
	return &userService{userRepository: r}
}

func (u *userService) FindUser(c *fiber.Ctx) error {
	criteria := c.Params("criteria")
	url := c.BaseURL() + c.OriginalURL()
	queryObj, err := querymap.FromURLStringToStruct[dtos.IsomorphQueryDto](url)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Cannot parse query"))
	}

	var usr *models.UserModel

	if queryObj.AsUsername || !utils.IsNumber(criteria) {
		usr, err = u.userRepository.FindByUsername(criteria)
	} else {
		usr, err = u.userRepository.FindById(int(utils.Numberize(criteria)))
	}

	if usr == nil {
		return c.Status(fiber.StatusNotFound).JSON(responses.NewNotFoundError(fmt.Sprintf("User with criteria %s not found", criteria)))
	}

	return c.Status(fiber.StatusOK).JSON(usr)
}

func (u *userService) CreateUser(c *fiber.Ctx) error {
	var req dtos.UserDto

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Cannot parse body"))
	}

	violations := validators.AppValidatorInstance.Validate(req)

	if violations != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responses.NewValidationError(violations.Message, violations.Violations))
	}

	existed, err := u.userRepository.FindByUsername(req.Username)

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Cannot find existed user..."))
	}

	if existed != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responses.NewValidationError("User already exists", []responses.Violation{
			*responses.NewViolation("User already exists", "username"),
		}))
	}

	usr, err := u.userRepository.Create(&req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Cannot create user..."))
	}

	if usr == nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responses.NewValidationError("User already exists", []responses.Violation{
			*responses.NewViolation("User already exists", "username"),
		}))
	}

	return c.Status(fiber.StatusCreated).JSON(usr)
}

func (u *userService) UpdateUser(c *fiber.Ctx) error {
	var req dtos.UpdateUserDto
	criteria := c.Params("criteria")

	url := c.BaseURL() + c.OriginalURL()
	queryObj, err := querymap.FromURLStringToStruct[dtos.IsomorphQueryDto](url)

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Cannot parse body"))
	}

	violations := validators.AppValidatorInstance.Validate(req)

	if violations != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responses.NewValidationError(violations.Message, violations.Violations))
	}

	var oldExisted *models.UserModel

	if queryObj.AsUsername || !utils.IsNumber(criteria) {
		oldExisted, err = u.userRepository.Find(criteria)
	} else {
		oldExisted, err = u.userRepository.FindById(int(utils.Numberize(criteria)))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Cannot find existed user"))
	}

	if oldExisted == nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responses.NewValidationError("User with this criter does not exists", []responses.Violation{
			*responses.NewViolation("User with this criteria does not exists", "criteria"),
		}))
	}

	newExisted, err := u.userRepository.FindByUsername(req.Username)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Find new existed error..."))
	}

	if newExisted != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responses.NewValidationError("User with this username already exists", []responses.Violation{
			*responses.NewViolation("User already exists", "username"),
		}))
	}

	usr, err := u.userRepository.Update(criteria, &req)

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Cannot update user"))
	}

	if usr == nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responses.NewValidationError("User with this username already exists", []responses.Violation{
			*responses.NewViolation("User with this username already exists", "username"),
		}))
	}

	return c.Status(fiber.StatusCreated).JSON(usr)
}

func (u *userService) DeleteUser(c *fiber.Ctx) error {
	criteria := c.Params("criteria")

	url := c.BaseURL() + c.OriginalURL()
	queryObj, err := querymap.FromURLStringToStruct[dtos.IsomorphQueryDto](url)

	var usr *models.UserModel

	if queryObj.AsUsername || !utils.IsNumber(criteria) {
		usr, err = u.userRepository.FindByUsername(criteria)
	} else {
		usr, err = u.userRepository.FindById(int(utils.Numberize(criteria)))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError())
	}

	if usr == nil {
		return c.Status(fiber.StatusNotFound).JSON(responses.NewNotFoundError(fmt.Sprintf("User with criteria %s not found", criteria)))
	}

	is_deleted := u.userRepository.Delete(criteria)

	if !is_deleted {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Cannot delete user..."))
	}

	return c.SendStatus(fiber.StatusNoContent)
}
