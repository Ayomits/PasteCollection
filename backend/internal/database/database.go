package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDatabase interface {
	Connect(dsn string) *pgxpool.Pool
}

type postgresDatabase struct{}

func (d *postgresDatabase) Connect(dsn string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		panic("Cannot connect to db")
	}
	return pool
}

func NewPostgresDatabase() PostgresDatabase {
	return &postgresDatabase{}
}
