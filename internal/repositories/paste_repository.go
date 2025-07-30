package repositories

import (
	"api/internal/dtos"
	"api/internal/models"
	"api/internal/utils"
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PasteRepository interface {
	FindByTitleOrId(title string) (*models.PasteModel, error)
	Create(dto *dtos.PasteDto) (*models.PasteModel, error)
	UpdateById(id string)
	DeleteById(id int)
}

const (
	FindByTitleStrict = "SELECT id, title, paste, tags, created_at, updated_at FROM pastes WHERE title=$1"
	FindByIdStrict    = "SELECT id, title, paste, tags, created_at, updated_at FROM pastes WHERE id=$1"
	CreatePasteSql    = "INSERT INTO pastes (title, paste, tags) VALUES($1, $2, $3) RETURNING id, title, paste, created_at, updated_at"
)

type pasteRepository struct {
	pool *pgxpool.Pool
}

func NewPasteRepository(p *pgxpool.Pool) PasteRepository {
	return &pasteRepository{pool: p}
}

func (p *pasteRepository) FindByTitleOrId(value string) (*models.PasteModel, error) {
	var paste models.PasteModel
	var sql string = FindByTitleStrict

	if utils.IsNumber(value) {
		sql = FindByIdStrict
	}

	err := p.pool.QueryRow(
		context.Background(),
		sql,
		value,
	).Scan(
		&paste.Id,
		&paste.Title,
		&paste.Paste,
		&paste.Tags,
		&paste.CreatedAt,
		&paste.UpdatedAt,
	)

	if err != nil && errors.Is(pgx.ErrNoRows, err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &paste, nil
}

func (p *pasteRepository) Create(dto *dtos.PasteDto) (*models.PasteModel, error) {
	var paste models.PasteModel

	err := p.pool.QueryRow(
		context.Background(),
		CreatePasteSql,
		dto.Title,
		dto.Paste,
		dto.Tags,
	).Scan(
		&paste.Id,
		&paste.Title,
		&paste.Tags,
		&paste.Paste,
		&paste.CreatedAt,
		&paste.UpdatedAt,
	)

	if err != nil {
		log.Printf("Create new paste failed: %v", err)
		return nil, err
	}

	return &paste, nil
}

func (p *pasteRepository) DeleteById(id int) {
	panic("unimplemented")
}

func (p *pasteRepository) UpdateById(id string) {
	panic("unimplemented")
}
