package services

import (
	"api/internal/dtos"
	"api/internal/repositories"
	"api/internal/responses"
	"api/internal/services/querymap"
	"api/internal/services/validators"

	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	Find(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(r repositories.UserRepository) UserService {
	return &userService{userRepository: r}
}

func (u *userService) Find(c *fiber.Ctx) error {
	queryObj, err := querymap.FromURLStringToStruct[dtos.UserFiltersDto](c.BaseURL() + c.OriginalURL())
	if u.isEmptyQuery(queryObj) {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewBadRequestError("Query parametrs is empty"))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Cannot parse query parametrs.."))
	}
	filterViolations := validators.AppValidatorInstance.Validate(queryObj)
	if filterViolations != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(filterViolations)
	}

	result, err := u.userRepository.Find(queryObj)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Error while query to db"))
	}

	if result == nil {
		return c.Status(fiber.StatusNotFound).JSON(responses.NewNotFoundError("User not found"))
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (u *userService) Create(c *fiber.Ctx) error {
	var body *dtos.UserDto

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Cannot parse body..."))
	}

	bodyViolations := validators.AppValidatorInstance.Validate(body)

	if bodyViolations != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(bodyViolations)
	}

	strict := true
	existed, err := u.userRepository.Find(&dtos.UserFiltersDto{
		Username: &body.Username,
		SocialId: &body.SocialId,
		Strict:   &strict,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Error while query to db"))
	}

	if existed != nil {
		return c.
			Status(fiber.StatusUnprocessableEntity).
			JSON(responses.NewValidationError("User already exists", []responses.Violation{
				*responses.NewViolation("User already exists", "username"),
				*responses.NewViolation("User already exists", "social_id"),
			}))
	}

	newUsr, err := u.userRepository.Create(body)

	return c.Status(fiber.StatusOK).JSON(newUsr)
}

func (u *userService) Update(c *fiber.Ctx) error {
	var body dtos.UpdateUserDto

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Cannot parse body..."))
	}

	queryObj, err := querymap.FromURLStringToStruct[dtos.UserFiltersDto](c.BaseURL() + c.OriginalURL())
	if u.isEmptyQuery(queryObj) {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewBadRequestError("Query parametrs is empty"))
	}

	filterViolations := validators.AppValidatorInstance.Validate(queryObj)
	bodyViolations := validators.AppValidatorInstance.Validate(body)

	if filterViolations != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(filterViolations)
	}

	if bodyViolations != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(bodyViolations)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Cannot parse query parametrs.."))
	}

	result, err := u.userRepository.Find(queryObj)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Error while query to db"))
	}

	if result == nil {
		return c.Status(fiber.StatusNotFound).JSON(responses.NewNotFoundError("User not found"))
	}

	newUsr, err := u.userRepository.Update(queryObj, &body)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Error while query to db"))
	}

	if newUsr == nil {
		return c.Status(fiber.StatusNotFound).JSON(responses.NewNotFoundError("User not found"))
	}

	return c.Status(fiber.StatusOK).JSON(newUsr)
}

func (u *userService) Delete(c *fiber.Ctx) error {
	queryObj, err := querymap.FromURLStringToStruct[dtos.UserFiltersDto](c.BaseURL() + c.OriginalURL())
	if u.isEmptyQuery(queryObj) {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewBadRequestError("Query parametrs is empty"))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Cannot parse query parametrs.."))
	}

	filterViolations := validators.AppValidatorInstance.Validate(queryObj)

	if filterViolations != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(filterViolations)
	}

	existed, err := u.userRepository.Find(queryObj)

	if existed == nil {
		return c.Status(fiber.StatusNotFound).JSON(responses.NewNotFoundError("User not found"))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Error while quering existed user"))
	}

	_, err = u.userRepository.Delete(queryObj)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Error while deleting user"))
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (u *userService) isEmptyQuery(q *dtos.UserFiltersDto) bool {
	return q.DisplayName == nil && q.Id == nil && q.Username == nil && q.SocialId == nil
}
