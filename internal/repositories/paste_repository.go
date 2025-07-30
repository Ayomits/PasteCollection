package repositories

import (
	"api/internal/dtos"
	"api/internal/enums"
	"api/internal/models"
	"api/internal/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	sb "github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PasteRepository interface {
	Search(query dtos.PastesSearchQueryDto) (*[]models.PasteModel, error, bool)
	Find(criteria string, query *dtos.FindByTitleQueryDto) (*models.PasteModel, error)
	FindById(id int) (*models.PasteModel, error)
	FindByTitle(value string, strict *bool) (*models.PasteModel, error)
	Create(dto *dtos.PasteDto) (*models.PasteModel, error)
	Update(criteria string, dto *dtos.PasteDto) (*models.PasteModel, error)
	Delete(criteria string) bool
}

type pasteRepository struct {
	pool *pgxpool.Pool
}

func NewPasteRepository(p *pgxpool.Pool) PasteRepository {
	return &pasteRepository{pool: p}
}

func (p *pasteRepository) Search(query dtos.PastesSearchQueryDto) (*[]models.PasteModel, error, bool) {
	sql := sb.PostgreSQL.NewSelectBuilder()
	sql.Select("id", "title", "paste", "tags", "created_at", "updated_at")
	sql.From("pastes")

	limit := 10
	if query.Pagination != nil && query.Pagination.Limit != nil {
		limit = *query.Pagination.Limit
	}
	sql.Limit(limit + 1)

	if query.Pagination != nil && query.Pagination.StartFrom != nil {
		if query.Pagination.Order != nil && *query.Pagination.Order == enums.PaginationPrev {
			sql.Where(sql.GreaterThan("id", query.Pagination.StartFrom))
		} else {
			sql.Where(sql.GreaterThan("id", query.Pagination.StartFrom))
		}
	}

	if query.Filter != nil && query.Filter.Search != nil {
		searchTerm := *query.Filter.Search
		if query.Filter.Strict != nil && *query.Filter.Strict {
			sql.Where(
				sql.Or(
					sql.Equal("paste", searchTerm),
					sql.Equal("title", searchTerm),
				),
			)
		} else {
			likePattern := "%" + searchTerm + "%"
			sql.Where(
				sql.Or(
					sql.Like("paste", likePattern),
					sql.Like("title", likePattern),
				),
			)
		}
	}

	if query.Filter != nil && query.Filter.Tags != nil && len(*query.Filter.Tags) > 0 {
		tags := *query.Filter.Tags
		subQuery := sb.PostgreSQL.NewSelectBuilder()
		subQuery.Select("1")
		subQuery.From("unnest(tags) AS elem")

		var tagConditions []string
		for _, tag := range tags {
			tagConditions = append(tagConditions,
				subQuery.Like("LOWER(elem)", "LOWER("+subQuery.Var("%"+tag+"%")+")"))
		}

		subQuery.Where(strings.Join(tagConditions, " OR "))
		sql.Where("EXISTS (" + subQuery.String() + ")")
	}

	queryStr, args := sql.Build()
	rows, err := p.pool.Query(context.Background(), queryStr, args...)
	if err != nil {
		return nil, err, false
	}
	defer rows.Close()

	var pastes []models.PasteModel
	for rows.Next() {
		var paste models.PasteModel
		err := rows.Scan(
			&paste.Id,
			&paste.Title,
			&paste.Paste,
			&paste.Tags,
			&paste.CreatedAt,
			&paste.UpdatedAt,
		)
		if err != nil {
			return nil, err, false
		}
		pastes = append(pastes, paste)
	}

	hasNext := len(pastes) > limit

	return &pastes, nil, hasNext
}

func (p *pasteRepository) Find(value string, query *dtos.FindByTitleQueryDto) (*models.PasteModel, error) {
	if utils.IsNumber(value) {
		return p.FindById(int(utils.Numberize(value)))
	}
	return p.FindByTitle(value, &query.Strict)
}

func (p *pasteRepository) FindById(id int) (*models.PasteModel, error) {
	var paste models.PasteModel
	sb := sb.PostgreSQL.NewSelectBuilder()
	sb.Select("id", "title", "paste", "tags", "created_at", "updated_at")
	sb.From("pastes")
	sb.Where(sb.Equal("id", id))

	query, args := sb.Build()

	err := p.pool.QueryRow(
		context.Background(),
		query,
		args...,
	).Scan(
		&paste.Id,
		&paste.Title,
		&paste.Paste,
		&paste.Tags,
		&paste.CreatedAt,
		&paste.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &paste, nil
}

func (p *pasteRepository) FindByTitle(title string, strict *bool) (*models.PasteModel, error) {
	var paste models.PasteModel
	sb := sb.PostgreSQL.NewSelectBuilder()
	sb.Select("id", "title", "paste", "tags", "created_at", "updated_at")
	sb.From("pastes")

	if strict != nil && *strict {
		sb.Where(sb.Equal("title", title))
	} else {
		sb.Where(sb.Like("title", title))
	}

	query, args := sb.Build()

	err := p.pool.QueryRow(
		context.Background(),
		query,
		args...,
	).Scan(
		&paste.Id,
		&paste.Title,
		&paste.Paste,
		&paste.Tags,
		&paste.CreatedAt,
		&paste.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query error: %w", err)
	}

	return &paste, nil
}

func (p *pasteRepository) Create(dto *dtos.PasteDto) (*models.PasteModel, error) {
	var paste models.PasteModel
	sb := sb.PostgreSQL.NewInsertBuilder()
	sb.InsertInto("pastes")
	sb.Cols("title", "paste", "tags")
	sb.Values(dto.Title, dto.Paste, dto.Tags)
	sb.Returning("id", "title", "paste", "tags", "created_at", "updated_at")

	query, args := sb.Build()

	err := p.pool.QueryRow(
		context.Background(),
		query,
		args...,
	).Scan(
		&paste.Id,
		&paste.Title,
		&paste.Paste,
		&paste.Tags,
		&paste.CreatedAt,
		&paste.UpdatedAt,
	)

	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			switch pgError.Code {
			case "23505":
				return nil, errors.New("duplicate entry")
			default:
				return nil, err
			}
		}
		return nil, err
	}

	return &paste, nil
}

func (p *pasteRepository) Update(criteria string, dto *dtos.PasteDto) (*models.PasteModel, error) {
	if utils.IsNumber(criteria) {
		return p.UpdateById(int(utils.Numberize(criteria)), dto)
	}
	return p.UpdateByTitle(criteria, dto)
}

func (p *pasteRepository) UpdateById(id int, dto *dtos.PasteDto) (*models.PasteModel, error) {
	var paste models.PasteModel
	sb := sb.PostgreSQL.NewUpdateBuilder()
	sb.Update("pastes")
	sb.Set(
		sb.Assign("title", dto.Title),
		sb.Assign("paste", dto.Paste),
		sb.Assign("tags", dto.Tags),
		sb.Assign("updated_at", "NOW()"),
	)
	sb.Where(sb.Equal("id", id))
	sb.Returning("id", "title", "paste", "tags", "created_at", "updated_at")

	query, args := sb.Build()

	err := p.pool.QueryRow(
		context.Background(),
		query,
		args...,
	).Scan(
		&paste.Id,
		&paste.Title,
		&paste.Paste,
		&paste.Tags,
		&paste.CreatedAt,
		&paste.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &paste, nil
}

func (p *pasteRepository) UpdateByTitle(title string, dto *dtos.PasteDto) (*models.PasteModel, error) {
	var paste models.PasteModel
	sb := sb.PostgreSQL.NewUpdateBuilder()
	sb.Update("pastes")
	sb.Set(
		sb.Assign("title", dto.Title),
		sb.Assign("paste", dto.Paste),
		sb.Assign("tags", dto.Tags),
		sb.Assign("updated_at", "NOW()"),
	)
	sb.Where(sb.Equal("title", title))
	sb.Returning("id", "title", "paste", "tags", "created_at", "updated_at")

	query, args := sb.Build()

	err := p.pool.QueryRow(
		context.Background(),
		query,
		args...,
	).Scan(
		&paste.Id,
		&paste.Title,
		&paste.Paste,
		&paste.Tags,
		&paste.CreatedAt,
		&paste.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &paste, nil
}

func (p *pasteRepository) Delete(criteria string) bool {
	if utils.IsNumber(criteria) {
		return p.DeleteById(int(utils.Numberize(criteria)))
	}
	return p.DeleteByTitle(criteria)
}

func (p *pasteRepository) DeleteByTitle(criteria string) bool {
	sb := sb.PostgreSQL.NewDeleteBuilder()
	sb.DeleteFrom("pastes")
	sb.Where(sb.Equal("title", criteria))

	query, args := sb.Build()

	_, err := p.pool.Exec(context.Background(), query, args...)
	if err != nil {
		log.Printf("Cannot delete paste with criteria %s: %v", criteria, err)
		return false
	}
	return true
}

func (p *pasteRepository) DeleteById(id int) bool {
	sb := sb.PostgreSQL.NewDeleteBuilder()
	sb.DeleteFrom("pastes")
	sb.Where(sb.Equal("id", id))

	query, args := sb.Build()

	_, err := p.pool.Exec(context.Background(), query, args...)
	if err != nil {
		log.Printf("Cannot delete paste with id %d: %v", id, err)
		return false
	}
	return true
}
