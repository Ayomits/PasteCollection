package services

import (
	"api/internal/dtos"
	"api/internal/repositories"
	"api/internal/responses"
	"api/internal/services/querymap"
	"api/internal/services/validators"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type PasteService interface {
	Create(c *fiber.Ctx) error
	Search(c *fiber.Ctx) error
	Find(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type pasteService struct {
	pasteRepository repositories.PasteRepository
}

func NewPasteService(r repositories.PasteRepository) PasteService {
	return &pasteService{pasteRepository: r}
}

func (p *pasteService) Search(c *fiber.Ctx) error {
	origin := c.BaseURL()
	path := c.OriginalURL()
	url := origin + path

	queryObj, err := querymap.FromURLStringToStruct[dtos.PastesSearchQueryDto](url)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError(err.Error()))
	}

	violations := validators.AppValidatorInstance.Validate(queryObj)

	if violations != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(violations)
	}

	results, err, hasNext := p.pasteRepository.Search(*queryObj)

	if err != nil {
		log.Printf("Error while executing query: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError())
	}

	return c.Status(fiber.StatusOK).JSON(responses.NewPaginationResponse(results, hasNext))
}

func (p *pasteService) Find(c *fiber.Ctx) error {
	criteria := c.Params("criteria")

	origin := c.BaseURL()
	path := c.OriginalURL()
	url := origin + path

	queryObj, err := querymap.FromURLStringToStruct[dtos.FindByTitleQueryDto](url)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError())
	}

	existed, err := p.pasteRepository.Find(criteria, queryObj)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError())
	}

	if existed == nil {
		return c.Status(fiber.StatusNotFound).JSON(responses.NewNotFoundError(fmt.Sprintf("Paste with criteria %s not found", criteria)))
	}

	return c.Status(fiber.StatusOK).JSON(existed)
}

func (p *pasteService) Create(c *fiber.Ctx) error {
	var req dtos.PasteDto
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewInternalError("Invalid request body"))
	}

	if err := validators.AppValidatorInstance.Validate(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			responses.NewValidationError(err.Message, err.Violations),
		)
	}

	existed, err := p.pasteRepository.Find(req.Title, dtos.NewFindByTitleQueryDto(true))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewInternalError())
	}

	if existed != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			responses.NewValidationError("Paste with this title already exists", []responses.Violation{
				*responses.NewViolation("Already exists", "title"),
			}),
		)
	}

	paste, err := p.pasteRepository.Create(&req)

	if paste == nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			responses.NewValidationError("Paste with this title already exists", []responses.Violation{
				*responses.NewViolation("Already exists", "title"),
			}),
		)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			responses.NewInternalError("Failed to create new paste"),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(paste)
}

func (p *pasteService) Update(c *fiber.Ctx) error {
	var req dtos.PasteDto
	criteria := c.Params("criteria")

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewInternalError("Invalid request body"))
	}

	if err := validators.AppValidatorInstance.Validate(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			responses.NewValidationError(err.Message, err.Violations),
		)
	}

	existed, err := p.pasteRepository.Find(criteria, dtos.NewFindByTitleQueryDto(true))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewInternalError())
	}

	if existed == nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			responses.NewValidationError("Cannot update paste. Paste with this criteria does not exists", []responses.Violation{
				*responses.NewViolation("Does not exists", "title"),
			}),
		)
	}

	existed, err = p.pasteRepository.Find(req.Title, dtos.NewFindByTitleQueryDto(true))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewInternalError())
	}

	if existed != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			responses.NewValidationError("Cannot update paste. Paste with new title already exists", []responses.Violation{
				*responses.NewViolation("Already exists", "title"),
			}),
		)
	}

	updated, err := p.pasteRepository.Update(criteria, &req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			responses.NewInternalError("Failed to update paste..."),
		)
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

func (p *pasteService) Delete(c *fiber.Ctx) error {
	criteria := c.Params("criteria")

	existed, err := p.pasteRepository.Find(criteria, dtos.NewFindByTitleQueryDto(true))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError())
	}

	if existed == nil {
		return c.Status(fiber.StatusNotFound).JSON(responses.NewNotFoundError(fmt.Sprintf(`Paste with criteria "%s" not found`, criteria)))
	}

	isDeleted := p.pasteRepository.Delete(criteria)

	if !isDeleted {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Cannot delete paste..."))
	}

	return c.SendStatus(fiber.StatusNoContent)
}
