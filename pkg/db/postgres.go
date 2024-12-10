package postgres

import (
	"context"
	"fmt"
	"log"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Username string `env:"POSTGRES_USER" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	Host string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port string `env:"POSTGRES_PORT" env-default:"5432"`
	DbName string `env:"POSTGRES_DB" env-default:"postgres"`
}

type Database struct {
	Database *sqlx.DB
}

func New(config Config) (*Database, error) {
	dataSource := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s", config.Username, config.Password, config.DbName, config.Host, config.Port)
	db, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		log.Fatalln(err)
	}
	if _, err := db.Conn(context.Background()); err != nil {
		return nil, fmt.Errorf("fail %s", err)
	}

	return &Database{Database: db}, nil
}