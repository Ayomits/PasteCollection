package services

import (
	"api/internal/dtos"
	"api/internal/repositories"
	"api/internal/responses"
	"api/internal/services/validators"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type PasteService interface {
	Create(c *fiber.Ctx) error
	FindByIdOrTitle(c *fiber.Ctx) error
}

type pasteService struct {
	pasteRepository repositories.PasteRepository
}

func NewPasteService(r repositories.PasteRepository) PasteService {
	return &pasteService{pasteRepository: r}
}

func (p *pasteService) FindByIdOrTitle(c *fiber.Ctx) error {
	criteria := c.Params("criteria")

	existed, err := p.pasteRepository.FindByTitleOrId(criteria)

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

	existed, err := p.pasteRepository.FindByTitleOrId(req.Title)

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
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			responses.NewInternalError("Failed to create new paste"),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": paste,
	})
}
