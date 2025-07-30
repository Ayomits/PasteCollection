package services

import (
	"api/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

type PasteService interface {
	Create(c *fiber.Ctx) error
}

type pasteService struct {
	pasteRepository repositories.PasteRepository
}

func NewPasteService(r repositories.PasteRepository) PasteService {
	return &pasteService{pasteRepository: r}
}

func (p *pasteService) Create(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"success": true})
}
