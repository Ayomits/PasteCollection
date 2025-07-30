package controllers

import (
	"api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type PasteController interface {
	FindByIdOrTitle(c *fiber.Ctx) error
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

func (p *pasteController) FindByIdOrTitle(c *fiber.Ctx) error {
	return p.pasteService.FindByIdOrTitle(c)
}

func (p *pasteController) DeletePaste(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (p *pasteController) UpdatePaste(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (p *pasteController) CreatePaste(c *fiber.Ctx) error {
	return p.pasteService.Create(c)
}
