package repositories

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type PasteRepository interface{}

type pasteRepository struct {
	pool *pgxpool.Pool
}

func NewPasteRepository(p *pgxpool.Pool) PasteRepository {
	return &pasteRepository{pool: p}
}
