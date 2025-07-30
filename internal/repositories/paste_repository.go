package repositories

import (
	"api/internal/dtos"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PasteRepository interface {
	Create(dto dtos.PasteDto) dtos.PasteDto
	UpdateById(id string)
	DeleteById(id int)
}

const ()

type pasteRepository struct {
	pool *pgxpool.Pool
}

func NewPasteRepository(p *pgxpool.Pool) PasteRepository {
	return &pasteRepository{pool: p}
}

func (p *pasteRepository) Create(dto dtos.PasteDto) dtos.PasteDto {
	panic("unimplemented")
}

func (p *pasteRepository) DeleteById(id int) {
	panic("unimplemented")
}

func (p *pasteRepository) UpdateById(id string) {
	panic("unimplemented")
}
