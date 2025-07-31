package services

import (
	"api/internal/repositories"

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

func (p *pasteService) Create(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (p *pasteService) Delete(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (p *pasteService) Find(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (p *pasteService) Search(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (p *pasteService) Update(c *fiber.Ctx) error {
	panic("unimplemented")
}

func NewPasteService(r repositories.PasteRepository) PasteService {
	return &pasteService{pasteRepository: r}
}
