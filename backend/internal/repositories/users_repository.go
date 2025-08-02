package repositories

import (
	"api/internal/dtos"
	"api/internal/models"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	FindUserSql   = "SELECT id, username, display_name, social_id FROM users %s"
	CreateUserSql = "INSERT INTO users (username, display_name, social_id) VALUES ($1, $2, $3) RETURNING id, username, display_name, social_id"
	UpdateUserSql = "UPDATE users SET username=$1, display_name=$2 %s RETURNING id, username, display_name, social_id"
	DeleteUserSql = "DELETE FROM users %s"
)

type UserRepository interface {
	Create(dto *dtos.UserDto) (*models.UserModel, error)
	Find(filter *dtos.UserFiltersDto) (*models.UserModel, error)
	Update(filter *dtos.UserFiltersDto, dto *dtos.UpdateUserDto) (*models.UserModel, error)
	Delete(filter *dtos.UserFiltersDto) (bool, error)
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(p *pgxpool.Pool) UserRepository {
	return &userRepository{pool: p}
}

func (u *userRepository) Find(filter *dtos.UserFiltersDto) (*models.UserModel, error) {
	var usr models.UserModel
	condition, args := u.buildFilters(filter, 0)

	err := u.pool.QueryRow(context.Background(), fmt.Sprintf(FindUserSql, condition), args...).Scan(
		&usr.Id,
		&usr.Username,
		&usr.DisplayName,
		&usr.SocialId,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &usr, nil
}

func (u *userRepository) Create(dto *dtos.UserDto) (*models.UserModel, error) {
	var usr models.UserModel

	err := u.pool.QueryRow(context.Background(), CreateUserSql, &dto.Username, &dto.DisplayName, &dto.SocialId).Scan(
		&usr.Id,
		&usr.Username,
		&usr.DisplayName,
		&usr.SocialId,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &usr, nil
}

func (u *userRepository) Delete(filter *dtos.UserFiltersDto) (bool, error) {
	condition, args := u.buildFilters(filter, 0)

	_, err := u.pool.Query(context.Background(), fmt.Sprintf(DeleteUserSql, condition), args...)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (u *userRepository) Update(filter *dtos.UserFiltersDto, dto *dtos.UpdateUserDto) (*models.UserModel, error) {
	var usr models.UserModel
	condition, args := u.buildFilters(filter, 2)

	allArgs := make([]interface{}, 0, 2+len(args))
	allArgs = append(allArgs, &dto.Username, &dto.DisplayName)
	allArgs = append(allArgs, args...)

	fmt.Println(UpdateUserSql + " " + condition)

	err := u.pool.
		QueryRow(context.Background(), fmt.Sprintf(UpdateUserSql, condition), allArgs...).
		Scan(
			&usr.Id,
			&usr.Username,
			&usr.DisplayName,
			&usr.SocialId,
		)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &usr, nil
}

func (u *userRepository) buildFilters(filter *dtos.UserFiltersDto, startFrom int) (string, []interface{}) {
	var conditions []string = []string{}
	var args []any = []any{}
	var position int = startFrom

	likePattern := func(value string) string {
		return "%" + value + "%"
	}

	if filter.Id != nil {
		position++
		conditions = append(conditions, fmt.Sprintf("id=$%d", position))
		args = append(args, filter.Id)
	}

	if filter.SocialId != nil {
		position++
		conditions = append(conditions, fmt.Sprintf("social_id=$%d", position))
		args = append(args, filter.SocialId)
	}

	if filter.Username != nil {
		position++
		if filter.Strict != nil && *filter.Strict {
			conditions = append(conditions, fmt.Sprintf("username=$%d", position))
			args = append(args, filter.Username)
		} else {
			conditions = append(conditions, fmt.Sprintf("username LIKE $%d", position))
			args = append(args, likePattern(*filter.Username))
		}
	}

	if filter.DisplayName != nil {
		position++
		if filter.Strict != nil && *filter.Strict {
			conditions = append(conditions, fmt.Sprintf("display_name=$%d", position))
			args = append(args, filter.Username)
		} else {
			conditions = append(conditions, fmt.Sprintf("display_name LIKE $%d", position))
			args = append(args, likePattern(*filter.Username))
		}
	}

	var sign string

	if filter.MatchAll != nil && *filter.MatchAll {
		sign = "AND"
	} else {
		sign = "OR"
	}

	if len(conditions) > 0 {
		if len(conditions) == 1 {
			return "WHERE " + strings.Join(conditions, " "), args
		}
		return "WHERE " + strings.Join(conditions, " "+sign+" "), args
	}
	return "", args
}
