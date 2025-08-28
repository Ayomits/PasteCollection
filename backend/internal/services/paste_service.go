package services

import (
	"api/internal/dtos"
	"api/internal/repositories"
	"api/internal/responses"
	"api/internal/services/querymap"
	"api/internal/services/validators"

	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type PasteService interface {
	Create(c *fiber.Ctx) error
	Find(c *fiber.Ctx) error
	Count(c *fiber.Ctx) error
	Increment(c *fiber.Ctx) error
	Search(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type pasteService struct {
	pasteRepository repositories.PasteRepository
}

func NewPasteService(r repositories.PasteRepository) PasteService {
	return &pasteService{pasteRepository: r}
}

func (p *pasteService) Count(c *fiber.Ctx) error {
	queryObj, err := querymap.FromURLStringToStruct[dtos.PastesFilterDto](c.BaseURL() + c.OriginalURL())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError())
	}

	violations := validators.AppValidatorInstance.Validate(queryObj)

	if violations != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responses.NewValidationError(violations.Message, violations.Violations))
	}

	count, err := p.pasteRepository.Count(queryObj)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError())
	}

	return c.Status(fiber.StatusOK).JSON(responses.NewCountResponse(*count))
}

func (p *pasteService) Increment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responses.NewValidationError("Is not integer id", []responses.Violation{
			*responses.NewViolation("Not integer", "id"),
		}))
	}

	existed, err := p.pasteRepository.FindOne(&dtos.PastesFilterDto{
		PasteId: &id,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError())
	}

	if existed == nil {
		return c.Status(fiber.StatusNotFound).JSON(responses.NewNotFoundError("Paste not existed"))
	}

	updated, err := p.pasteRepository.IncrementViews(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError())
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

func (p *pasteService) Create(c *fiber.Ctx) error {
	var body dtos.CreatePasteDto

	if err := c.BodyParser(&body); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError())
	}

	violations := validators.AppValidatorInstance.Validate(body)

	if violations != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(violations)
	}

	strict := true
	existed, err := p.pasteRepository.FindOne(&dtos.PastesFilterDto{
		Search: &body.Title,
		Strict: &strict,
	})

	if err != nil {
		log.Errorf("While quering db %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError())
	}

	if existed != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responses.NewValidationError("Paste already exists", []responses.Violation{
			*responses.NewViolation("Paste already exists", "title"),
		}))
	}

	newPaste, err := p.pasteRepository.Create(&body)

	if err != nil {
		log.Errorf("While quering db %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError())
	}

	return c.Status(fiber.StatusCreated).JSON(newPaste)
}

func (p *pasteService) Delete(c *fiber.Ctx) error {
	queryObj, err := querymap.FromURLStringToStruct[dtos.PastesFilterDto](c.BaseURL() + c.OriginalURL())

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Failed to parse query..."))
	}

	if p.isEmptyFilter(queryObj) {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewBadRequestError("Empty query parametrs"))
	}

	existed, err := p.pasteRepository.FindOne(queryObj)

	if existed == nil {
		return c.Status(fiber.StatusNotFound).JSON(responses.NewBadRequestError("Paste not found"))
	}

	_, err = p.pasteRepository.Delete(queryObj)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Error while quering db..."))
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (p *pasteService) Find(c *fiber.Ctx) error {
	url := c.BaseURL() + c.OriginalURL()
	queryObj, err := querymap.FromURLStringToStruct[dtos.PastesFilterDto](url)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Failed to parse query..."))
	}
	existed, err := p.pasteRepository.FindOne(queryObj)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Error while quering db..."))
	}

	if existed == nil {
		return c.Status(fiber.StatusNotFound).JSON(responses.NewBadRequestError("Paste not found"))
	}

	return c.Status(fiber.StatusOK).JSON(existed)
}

func (p *pasteService) Search(c *fiber.Ctx) error {
	url := c.BaseURL() + c.OriginalURL()
	queryObj, err := querymap.FromURLStringToStruct[dtos.PastesSearchQueryDto](url)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Failed to parse query..."))
	}

	entries, err := p.pasteRepository.FindMany(queryObj.Filter, queryObj.Pagination)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Error while quering db..."))
	}

	if entries == nil {
		return c.Status(fiber.StatusNotFound).JSON(responses.NewBadRequestError("Paste not found"))
	}

	var limit int = 10

	if queryObj.Pagination != nil && queryObj.Pagination.Limit != nil {
		limit = *queryObj.Pagination.Limit
	}

	items := entries[0:limit]

	return c.Status(fiber.StatusOK).JSON(responses.NewPaginationResponse(&items, len(entries) > limit))
}

func (p *pasteService) Update(c *fiber.Ctx) error {
	var body dtos.UpdatePasteDto

	queryObj, err := querymap.FromURLStringToStruct[dtos.PastesFilterDto](c.BaseURL() + c.OriginalURL())
	filterViolations := validators.AppValidatorInstance.Validate(queryObj)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError("Failed to parse query..."))
	}

	if filterViolations != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(filterViolations)
	}

	if p.isEmptyFilter(queryObj) {
		return c.Status(fiber.StatusBadRequest).JSON(responses.NewBadRequestError("Empty query parametrs"))
	}

	if err := c.BodyParser(&body); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError())
	}

	existed, err := p.pasteRepository.FindOne(queryObj)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError())
	}

	if existed == nil {
		return c.Status(fiber.StatusNotFound).JSON(responses.NewNotFoundError("paste not found"))
	}

	bodyViolations := validators.AppValidatorInstance.Validate(body)

	if bodyViolations != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(bodyViolations)
	}

	if existed.Title != body.Title {
		strict := true
		existed, err = p.pasteRepository.FindOne(&dtos.PastesFilterDto{
			Search: &body.Title,
			Strict: &strict,
		})

		if err != nil {
			log.Errorf("While quering db %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError())
		}

		if existed != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(responses.NewValidationError("Paste already exists", []responses.Violation{
				*responses.NewViolation("Paste already exists", "title"),
			}))
		}
	}

	newPaste, err := p.pasteRepository.Update(queryObj, &body)

	if err != nil {
		log.Errorf("While quering db %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(responses.NewInternalError())
	}

	return c.Status(fiber.StatusCreated).JSON(newPaste)
}

func (p *pasteService) isEmptyFilter(filter *dtos.PastesFilterDto) bool {
	return filter.Search == nil && filter.UserId == nil && filter.PasteId == nil
}
