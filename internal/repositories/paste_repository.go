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
	Find(criteria string) (*models.PasteModel, error)
	Create(dto *dtos.PasteDto) (*models.PasteModel, error)
	Update(criteria string, dto *dtos.PasteDto) (*models.PasteModel, error)
	Delete(criteria string) bool
}

const (
	FindByTitleStrict   = "SELECT id, title, paste, tags, created_at, updated_at FROM pastes WHERE title=$1"
	FindByIdStrict      = "SELECT id, title, paste, tags, created_at, updated_at FROM pastes WHERE id=$1"
	CreatePasteSql      = "INSERT INTO pastes (title, paste, tags) VALUES($1, $2, $3) RETURNING id, title, paste, created_at, updated_at"
	UpdateByIdSql       = "UPDATE pastes SET title=$1,paste=$2,tags=$3 WHERE id=$4 RETURNING id, title, paste, tags, created_at, updated_at"
	UpdateByTitleSql    = "UPDATE pastes SET title=$1,paste=$2,tags=$3 WHERE title=$4 RETURNING id, title, paste, tags, created_at, updated_at"
	DeleteByIdStrict    = "DELETE FROM pastes WHERE id=$1"
	DeleteByTitleStrict = "DELETE FROM pastes WHERE title=$1"
)

type pasteRepository struct {
	pool *pgxpool.Pool
}

func NewPasteRepository(p *pgxpool.Pool) PasteRepository {
	return &pasteRepository{pool: p}
}

func (p *pasteRepository) Find(value string) (*models.PasteModel, error) {
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

func (p *pasteRepository) Delete(criteria string) bool {
	var sql string = DeleteByTitleStrict

	if utils.IsNumber(criteria) {
		sql = DeleteByIdStrict
	}

	err := p.pool.QueryRow(context.Background(), sql, criteria)

	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (p *pasteRepository) Update(criteria string, dto *dtos.PasteDto) (*models.PasteModel, error) {
	var paste models.PasteModel
	var sql string = UpdateByTitleSql

	if utils.IsNumber(criteria) {
		sql = UpdateByIdSql
	}

	err := p.pool.QueryRow(
		context.Background(),
		sql,
		dto.Title,
		dto.Paste,
		dto.Tags,
		criteria,
	).Scan(
		&paste.Id,
		&paste.Title,
		&paste.Paste,
		&paste.Tags,
		&paste.CreatedAt,
		&paste.UpdatedAt,
	)

	if err != nil && !errors.Is(pgx.ErrNoRows, err) {
		log.Printf("Update paste failed: %v", err)
		return nil, err
	}

	return &paste, nil
}
