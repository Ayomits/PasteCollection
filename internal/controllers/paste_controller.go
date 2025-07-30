package controllers

import (
	"api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type PasteController interface {
	SearchPaste(c *fiber.Ctx) error
	FindPaste(c *fiber.Ctx) error
	CreatePaste(c *fiber.Ctx) error
	DeletePaste(c *fiber.Ctx) error
	UpdatePaste(c *fiber.Ctx) error
}

type pasteController struct {
	pasteService services.PasteService
}

func NewPasteController(pasteService services.PasteService) PasteController {
	return &pasteController{pasteService: pasteService}
}

func (p *pasteController) SearchPaste(c *fiber.Ctx) error {
	return p.pasteService.Search(c)
}

func (p *pasteController) FindPaste(c *fiber.Ctx) error {
	return p.pasteService.Find(c)
}

func (p *pasteController) DeletePaste(c *fiber.Ctx) error {
	return p.pasteService.Delete(c)
}

func (p *pasteController) UpdatePaste(c *fiber.Ctx) error {
	return p.pasteService.Update(c)
}

func (p *pasteController) CreatePaste(c *fiber.Ctx) error {
	return p.pasteService.Create(c)
}
