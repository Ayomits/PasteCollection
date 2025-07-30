package repositories

import (
	"api/internal/dtos"
	"api/internal/enums"
	"api/internal/models"
	"api/internal/utils"
	"context"
	"errors"
	"log"
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PasteRepository interface {
	Search(query dtos.PastesSearchQueryDto) (*[]models.PasteModel, error, bool)
	Find(criteria string) (*models.PasteModel, error)
	Create(dto *dtos.PasteDto) (*models.PasteModel, error)
	Update(criteria string, dto *dtos.PasteDto) (*models.PasteModel, error)
	Delete(criteria string) bool
}

const (
	SearchSql           = "SELECT id, title, paste, tags, created_at, updated_at FROM pastes WHERE"
	FindByTitleStrict   = "SELECT id, title, paste, tags, created_at, updated_at FROM pastes WHERE title=$1"
	FindByIdStrict      = "SELECT id, title, paste, tags, created_at, updated_at FROM pastes WHERE id=$1"
	CreatePasteSql      = "INSERT INTO pastes (title, paste, tags) VALUES($1, $2, $3) RETURNING id, title, paste, tags, created_at, updated_at"
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

func (p *pasteRepository) Search(query dtos.PastesSearchQueryDto) (*[]models.PasteModel, error, bool) {
	sql := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sql.Select("id", "title", "paste", "tags", "created_at", "updated_at")
	sql.From("pastes")

	limit := 10
	if query.Pagination != nil && query.Pagination.Limit != nil {
		limit = *query.Pagination.Limit
	}
	sql.Limit(limit)

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
					sql.Like("LOWER(paste)", "LOWER("+sql.Var(likePattern)+")"),
					sql.Like("LOWER(title)", "LOWER("+sql.Var(likePattern)+")"),
				),
			)
		}
	}

	if query.Filter != nil && query.Filter.Tags != nil && len(*query.Filter.Tags) > 0 {
		tags := *query.Filter.Tags
		subQuery := sqlbuilder.PostgreSQL.NewSelectBuilder()
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

	hasNext := len(pastes) == limit

	return &pastes, nil, hasNext
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
				return nil, nil
			default:
				return nil, nil
			}
		}
	}

	return &paste, nil
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

func (p *pasteRepository) Delete(criteria string) bool {
	var sql string = DeleteByTitleStrict

	if utils.IsNumber(criteria) {
		sql = DeleteByIdStrict
	}

	_, err := p.pool.Query(context.Background(), sql, criteria)

	if err != nil && !errors.Is(pgx.ErrNoRows, err) {
		log.Printf("Cannot delete paste with criteria %s: %v", criteria, err)
		return false
	}
	return true
}
