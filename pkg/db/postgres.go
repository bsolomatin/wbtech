package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Username string `env:"POSTGRES_USER" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	Host     string `env:"POSTGRES_HOST" env-default:"db"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	DbName   string `env:"POSTGRES_DB" env-default:"postgres"`
}

type Database struct {
	Database *sqlx.DB
}

func New(config Config) (*Database, error) {
	dataSource := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", config.Username, config.Password, config.Host, config.Port, config.DbName)
	db, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		return nil, err
	}

	if _, err := db.Conn(context.Background()); err != nil {
		return nil, err
	}
	return &Database{
		Database: db,
	}, nil
}
