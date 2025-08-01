package repositories

import (
	"api/internal/dtos"
	"api/internal/models"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	CreatePasteSql = "INSERT INTO pastes (title, paste, user_id) VALUES ($1, $2, $3) RETURNING id, title, paste, user_id, created_at, updated_at"
	FindPasteSql   = "SELECT id, title, paste, user_id, created_at, updated_at FROM pastes %s"
	UpdatePasteSql = "UPDATE pastes SET title=$1, paste=$2, updated_at=now() %s RETURNING id, title, paste, user_id, created_at, updated_at"
	DeletePasteSql = "DELETE FROM pastes %s"
)

type PasteRepository interface {
	FindOne(filter *dtos.PastesFilterDto, pagination *dtos.PaginationDto) (*models.PasteModel, error)
	FindMany(filter *dtos.PastesFilterDto, pagination *dtos.PaginationDto) ([]*models.PasteModel, error)
	Create(dto *dtos.PasteDto) (*models.PasteModel, error)
	Update(filter *dtos.PastesFilterDto, dto *dtos.UpdatePasteDto) (*models.PasteModel, error)
	Delete(filter *dtos.PastesFilterDto) (bool, error)
}

type pasteRepository struct {
	pool *pgxpool.Pool
}

func NewPasteRepository(p *pgxpool.Pool) PasteRepository {
	return &pasteRepository{pool: p}
}

func (p *pasteRepository) FindOne(filter *dtos.PastesFilterDto, pagination *dtos.PaginationDto) (*models.PasteModel, error) {
	var paste models.PasteModel
	condition, args := p.buildFilters(filter, 0, pagination)
	err := p.pool.
		QueryRow(context.Background(), fmt.Sprintf(FindPasteSql, condition), args...).
		Scan(
			&paste.Id,
			&paste.Title,
			&paste.Paste,
			&paste.UserId,
			&paste.CreatedAt,
			&paste.UpdatedAt,
		)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &paste, nil
}

func (p *pasteRepository) FindMany(filter *dtos.PastesFilterDto, pagination *dtos.PaginationDto) ([]*models.PasteModel, error) {
	condition, args := p.buildFilters(filter, 0, pagination)
	rows, err := p.pool.Query(context.Background(), fmt.Sprintf(FindPasteSql, condition), args...)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var pastes []*models.PasteModel = []*models.PasteModel{}

	defer rows.Close()
	for rows.Next() {
		var paste models.PasteModel
		err := rows.
			Scan(
				&paste.Id,
				&paste.Title,
				&paste.Paste,
				&paste.UserId,
				&paste.CreatedAt,
				&paste.UpdatedAt,
			)
		if err != nil {
			return nil, err
		}
		pastes = append(pastes, &paste)
	}

	return pastes, nil
}

func (p *pasteRepository) Create(dto *dtos.PasteDto) (*models.PasteModel, error) {
	var paste models.PasteModel

	err := p.pool.QueryRow(context.Background(), CreatePasteSql, dto.Title, dto.Paste, dto.UserId).Scan(
		&paste.Id,
		&paste.Title,
		&paste.Paste,
		&paste.UserId,
		&paste.CreatedAt,
		&paste.UpdatedAt,
	)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &paste, nil
}

func (p *pasteRepository) Update(filter *dtos.PastesFilterDto, dto *dtos.UpdatePasteDto) (*models.PasteModel, error) {
	var paste models.PasteModel
	condition, args := p.buildFilters(filter, 2, nil)

	lastArgs := append([]any{dto.Title, dto.Paste}, args...)

	err := p.pool.QueryRow(context.Background(), fmt.Sprintf(UpdatePasteSql, condition), lastArgs...).Scan(
		&paste.Id,
		&paste.Title,
		&paste.Paste,
		&paste.UserId,
		&paste.CreatedAt,
		&paste.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &paste, nil
}

func (p *pasteRepository) Delete(filter *dtos.PastesFilterDto) (bool, error) {
	condition, args := p.buildFilters(filter, 0, nil)

	_, err := p.pool.Query(context.Background(), fmt.Sprintf(DeletePasteSql, condition), args...)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (p *pasteRepository) buildFilters(filter *dtos.PastesFilterDto, startFrom int, pagination *dtos.PaginationDto) (string, []interface{}) {
	var condition string = ""
	var conditions []string = []string{}
	var args []any = []any{}
	var position int = startFrom

	if filter != nil {
		if filter.Search != nil {
			position++
			if filter.Strict != nil && *filter.Strict {
				conditions = append(conditions, fmt.Sprintf("title = $%d", position))
				args = append(args, *filter.Search)
			} else {
				conditions = append(conditions, fmt.Sprintf("title ILIKE $%d", position))
				args = append(args, "%"+*filter.Search+"%")
			}
		}

		if filter.SocialId != nil {
			position++
			conditions = append(conditions, fmt.Sprintf("user_id = (SELECT id FROM users WHERE social_id = $%d)", position))
			args = append(args, filter.SocialId)
		}

		if filter.UserId != nil {
			position++
			conditions = append(conditions, fmt.Sprintf("user_id=$%d", position))
			args = append(args, filter.UserId)
		}

		if filter.PasteId != nil {
			position++
			conditions = append(conditions, fmt.Sprintf("id=$%d", position))
			args = append(args, filter.PasteId)
		}
	}

	if pagination != nil {
		if pagination.StartFrom != nil {
			position++
			if pagination.Order != nil && *pagination.Order == "next" {
				conditions = append(conditions, fmt.Sprintf("id > $%d", position))
			} else {
				conditions = append(conditions, fmt.Sprintf("id < $%d", position))
			}
			args = append(args, pagination.StartFrom)
		}

		if pagination.Sort != nil {
			conditions = append(conditions, fmt.Sprintf("ORDER BY id %s", *pagination.Order))
		}

		if pagination.Limit != nil {
			position++
			conditions = append(conditions, fmt.Sprintf("LIMIT $%d", position))
			args = append(args, *pagination.Limit+1)
		}
	}

	if len(conditions) > 0 {
		if len(conditions) > 1 {
			condition = fmt.Sprintf("WHERE %s", strings.Join(conditions, " OR "))
		} else {
			condition = fmt.Sprintf("WHERE %s", strings.Join(conditions, ""))
		}
	}

	return condition, args
}
