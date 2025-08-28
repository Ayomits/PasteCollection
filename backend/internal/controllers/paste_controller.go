package controllers

import (
	"api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type PasteController interface {
	FindSinglePaste(c *fiber.Ctx) error
	SearchPaste(c *fiber.Ctx) error
	Count(c *fiber.Ctx) error
	CreatePaste(c *fiber.Ctx) error
	DeletePaste(c *fiber.Ctx) error
	UpdatePaste(c *fiber.Ctx) error
	Increment(c *fiber.Ctx) error
}

type pasteController struct {
	pasteService services.PasteService
}

func NewPasteController(pasteService services.PasteService) PasteController {
	return &pasteController{pasteService: pasteService}
}

// CreatePaste godoc
//
//	@Summary		Создать новую пасту
//	@Description	Создает новую пасту с указанным содержимым
//	@Tags			Pastes
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dtos.CreatePasteDto	true	"Данные пасты"
//	@Success		201		{object}	models.PasteModel	"Созданная паста"
//	@Failure		401		{object}	responses.UnauthorizedError
//	@Failure		422		{object}	responses.ValidationError
//	@Failure		500		{object}	responses.InternalError
//	@Router			/api/pastes [post]
//	@Security		ApiKeyAuth
func (p *pasteController) CreatePaste(c *fiber.Ctx) error {
	return p.pasteService.Create(c)
}

// FindSinglePaste godoc
//
//	@Summary		Найти одну пасту по фильтру
//	@Description	Возвращает пасту по указанному фильтру
//	@Tags			Pastes
//	@Accept			json
//	@Produce		json
//	@Param			search		query		string				false	"Поиск по названию"
//	@Param			userId		query		int					false	"Айди автора"
//	@Param			socialId	query		string				false	"Айди в соц. сети"
//	@Param			pasteId		query		int					false	"Айди пасты"
//	@Param			strict		query		bool				false	"Строгое или нестрогое совпадение в названии"	default(false)
//	@Success		200			{object}	models.PasteModel	"Найденная паста"
//	@Failure		404			{object}	responses.NotFoundError
//	@Failure		401			{object}	responses.UnauthorizedError
//	@Failure		422			{object}	responses.ValidationError
//	@Failure		500			{object}	responses.InternalError
//	@Router			/api/pastes [get]
//	@Security		ApiKeyAuth
func (p *pasteController) FindSinglePaste(c *fiber.Ctx) error {
	return p.pasteService.Find(c)
}

// SearchPaste godoc
//
//	@Summary		Поиск паст
//	@Description	Поиск паст по различным критериям
//	@Tags			Pastes
//	@Accept			json
//	@Produce		json
//	@Param			pagination[sort]		query		string				false	"Сортировка результатов"	default("DESC")
//	@Param			pagination[startFrom]	query		int					false	"С какого айди начинать"	default(10)
//	@Param			pagination[limit]		query		int					true	"Лимит результатов"			default(10)
//	@Param			pagination[order]		query		int					false	"Сторона"					default(next)
//	@Param			filter[search]			query		string				false	"Поиск по названию"
//	@Param			filter[userId]			query		int					false	"Айди автора"
//	@Param			filter[socialId]		query		string				false	"Айди в соц. сети"
//	@Param			filter[pasteId]			query		int					false	"Айди пасты"
//	@Param			filter[strict]			query		bool				false	"Строгое или нестрогое совпадение в названии"	default(false)
//	@Success		200						{array}		models.PasteModel	"Список найденных паст"
//	@Failure		401						{object}	responses.UnauthorizedError
//	@Failure		422						{object}	responses.ValidationError
//	@Failure		500						{object}	responses.InternalError
//	@Router			/api/pastes/search [get]
//	@Security		ApiKeyAuth
func (p *pasteController) SearchPaste(c *fiber.Ctx) error {
	return p.pasteService.Search(c)
}

// Count godoc
//
//	@Summary		Получить количество паст
//	@Description	Возвращает общее количество паст в системе
//	@Tags			Pastes
//	@Accept			json
//	@Produce		json
//	@Param			search		query		string					false	"Поиск по названию"
//	@Param			userId		query		int						false	"Айди автора"
//	@Param			socialId	query		string					false	"Айди в соц. сети"
//	@Param			pasteId		query		int						false	"Айди пасты"
//	@Param			strict		query		bool					false	"Строгое или нестрогое совпадение в названии"	default(false)
//	@Success		200			{object}	responses.CountResponse	"Количество паст"
//	@Failure		401			{object}	responses.UnauthorizedError
//	@Failure		422			{object}	responses.ValidationError
//	@Failure		500			{object}	responses.InternalError
//	@Router			/api/pastes/count [get]
//	@Security		ApiKeyAuth
func (p *pasteController) Count(c *fiber.Ctx) error {
	return p.pasteService.Count(c)
}

// UpdatePaste godoc
//
//	@Summary		Обновить пасту
//	@Description	Обновляет содержимое пасты
//	@Tags			Pastes
//	@Accept			json
//	@Produce		json
//	@Param			search		query		string				false	"Поиск по названию"
//	@Param			userId		query		int					false	"Айди автора"
//	@Param			socialId	query		string				false	"Айди в соц. сети"
//	@Param			pasteId		query		int					false	"Айди пасты"
//	@Param			strict		query		bool				false	"Строгое или нестрогое совпадение в названии"	default(false)
//	@Param			request		body		dtos.UpdatePasteDto	true	"Данные для обновления"
//	@Success		200			{object}	models.PasteModel	"Обновленная паста"
//	@Failure		404			{object}	responses.NotFoundError
//	@Failure		401			{object}	responses.UnauthorizedError
//	@Failure		422			{object}	responses.ValidationError
//	@Failure		500			{object}	responses.InternalError
//	@Router			/api/pastes [put]
//	@Security		ApiKeyAuth
func (p *pasteController) UpdatePaste(c *fiber.Ctx) error {
	return p.pasteService.Update(c)
}

// Increment godoc
//
//	@Summary					Увеличить счетчик просмотров
//	@Description				Увеличивает счетчик просмотров для конкретной пасты
//	@Tags						Pastes
//	@Accept						json
//	@Produce					json
//	@securitydefinitions.apikey	ApiKeyAuth
//	@Success					200	{object}	models.PasteModel	"Обновленная паста"
//	@Failure					401	{object}	responses.UnauthorizedError
//	@Failure					422	{object}	responses.ValidationError
//	@Failure					500	{object}	responses.InternalError
//	@Param						id	path		int	true	"ID пасты"
//	@Router						/api/pastes/{id}/increment [patch]
//	@Security					ApiKeyAuth
func (p *pasteController) Increment(c *fiber.Ctx) error {
	return p.pasteService.Increment(c)
}

// DeletePaste godoc
//
//	@Summary		Удалить пасту
//	@Description	Удаляет пасту по указанному фильтру
//	@Tags			Pastes
//	@Accept			json
//	@Produce		json
//	@Param			search		query	string	false	"Поиск по названию"
//	@Param			userId		query	int		false	"Айди автора"
//	@Param			socialId	query	string	false	"Айди в соц. сети"
//	@Param			pasteId		query	int		false	"Айди пасты"
//	@Param			strict		query	bool	false	"Строгое или нестрогое совпадение в названии"	default(false)
//	@Success		204
//	@Failure		401	{object}	responses.UnauthorizedError
//	@Failure		404	{object}	responses.NotFoundError
//	@Failure		422	{object}	responses.ValidationError
//	@Failure		500	{object}	responses.InternalError
//	@Router			/api/pastes [delete]
//	@Security		ApiKeyAuth
func (p *pasteController) DeletePaste(c *fiber.Ctx) error {
	return p.pasteService.Delete(c)
}
