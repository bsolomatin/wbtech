package postgres

import (
	"context"
	"fmt"
	"log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	dataSource := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", config.Username, config.Password, "db", config.Port, config.DbName)
	db, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		log.Fatalln(err)
	}
	
	if _, err := db.Conn(context.Background()); err != nil {
		return nil, fmt.Errorf("fail %s", err)
	}
	return &Database{Database: db}, nil
}