package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDatabase interface {
	Connect(dsn string) *gorm.DB
}

type postgresDatabase struct{}

func (d *postgresDatabase) Connect(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Cannot connect to db")
	}
	return db
}

func NewPostgresDatabase() PostgresDatabase {
	return &postgresDatabase{}
}
