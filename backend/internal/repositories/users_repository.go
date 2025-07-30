package repositories

import (
	"api/internal/dtos"
	"api/internal/enums"
	"api/internal/models"
	"api/internal/utils"
	"context"
	"errors"
	"fmt"

	sb "github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Find(criteria string) (*models.UserModel, error)
	FindById(id int) (*models.UserModel, error)
	FindByUsername(username string) (*models.UserModel, error)

	Create(dto *dtos.UserDto) (*models.UserModel, error)

	Update(criteria string, dto *dtos.UserDto) (*models.UserModel, error)
	UpdateById(id int, dto *dtos.UserDto) (*models.UserModel, error)
	UpdateByUsername(username string, dto *dtos.UserDto) (*models.UserModel, error)

	Delete(criteria string) bool
	DeleteById(id int) bool
	DeleteByUsername(username string) bool
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(p *pgxpool.Pool) UserRepository {
	return &userRepository{pool: p}
}

func (u *userRepository) Find(criteria string) (*models.UserModel, error) {
	if utils.IsNumber(criteria) {
		return u.FindById(int(utils.Numberize(criteria)))
	}
	return u.FindByUsername(criteria)
}

func (u *userRepository) FindById(id int) (*models.UserModel, error) {
	var usr models.UserModel
	sql := sb.Select("id", "username", "social_id").From("users")
	sql.Where(sql.Equal("id", id))

	query, args := sql.Build()

	err := u.pool.QueryRow(context.Background(), query, args...).Scan(
		&usr.Id,
		&usr.Username,
		&usr.SocialId,
	)

	if err != nil && !errors.Is(pgx.ErrNoRows, err) {
		return nil, err
	}

	return &usr, nil
}

func (u *userRepository) FindByUsername(username string) (*models.UserModel, error) {
	var usr models.UserModel
	sql := sb.Select("id", "username", "social_id").From("users")
	sql.Where(sql.Equal("username", username))

	query, args := sql.Build()

	err := u.pool.QueryRow(context.Background(), query, args...).Scan(
		&usr.Id,
		&usr.Username,
		&usr.SocialId,
	)

	if err != nil && !errors.Is(pgx.ErrNoRows, err) {
		return nil, err
	}

	return &usr, nil
}

func (u *userRepository) Create(dto *dtos.UserDto) (*models.UserModel, error) {
	var usr models.UserModel
	sql := sb.NewInsertBuilder()
	sql.InsertInto("users").Cols("username", "social_id").Values(dto.Username).Returning("id", "username", "social_id")

	query, args := sql.Build()

	err := u.pool.QueryRow(context.Background(), query, args...).Scan(
		&usr.Id,
		&usr.Username,
		&usr.SocialId,
	)

	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			switch pgError.Code {
			case enums.DbCodeDuplicateKey:
				return nil, errors.New("duplicate entry")
			default:
				return nil, err
			}
		}
		return nil, err
	}

	return &usr, nil
}

func (u *userRepository) Update(criteria string, dto *dtos.UserDto) (*models.UserModel, error) {
	if utils.IsNumber(criteria) {
		return u.UpdateById(int(utils.Numberize(criteria)), dto)
	}
	return u.UpdateByUsername(criteria, dto)
}

func (u *userRepository) UpdateById(id int, dto *dtos.UserDto) (*models.UserModel, error) {
	var usr models.UserModel
	sql := sb.NewUpdateBuilder()
	sql.Where("id", fmt.Sprint(id)).Set(
		sql.Assign("username", dto.Username),
	).Returning("id", "username", "social_id")

	query, args := sql.Build()

	err := u.pool.QueryRow(context.Background(), query, args...).Scan(
		&usr.Id,
		&usr.Username,
		&usr.SocialId,
	)

	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			switch pgError.Code {
			case enums.DbCodeDuplicateKey:
				return nil, errors.New("duplicate entry")
			default:
				return nil, err
			}
		}
		return nil, err
	}

	return &usr, nil
}

func (u *userRepository) UpdateByUsername(username string, dto *dtos.UserDto) (*models.UserModel, error) {
	var usr models.UserModel
	sql := sb.NewUpdateBuilder()
	sql.Where("id", fmt.Sprint(username)).Set(
		sql.Assign("username", dto.Username),
	).Returning("id", "username", "social_id")

	query, args := sql.Build()

	err := u.pool.QueryRow(context.Background(), query, args...).Scan(
		&usr.Id,
		&usr.Username,
		&usr.SocialId,
	)

	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			switch pgError.Code {
			case enums.DbCodeDuplicateKey:
				return nil, errors.New("duplicate entry")
			default:
				return nil, err
			}
		}
		return nil, err
	}

	return &usr, nil
}

func (u *userRepository) Delete(criteria string) bool {
	if utils.IsNumber(criteria) {
		return u.DeleteById(int(utils.Numberize(criteria)))
	}
	return u.DeleteByUsername(criteria)
}

func (u *userRepository) DeleteById(id int) bool {
	sb := sb.PostgreSQL.NewDeleteBuilder()
	sb.DeleteFrom("users")
	sb.Where(sb.Equal("id", id))

	query, args := sb.Build()

	_, err := u.pool.Exec(context.Background(), query, args...)
	if err != nil {
		return false
	}
	return true
}

func (u *userRepository) DeleteByUsername(username string) bool {
	sb := sb.PostgreSQL.NewDeleteBuilder()
	sb.DeleteFrom("users")
	sb.Where(sb.Equal("username", username))

	query, args := sb.Build()

	_, err := u.pool.Exec(context.Background(), query, args...)
	if err != nil {
		return false
	}
	return true
}
