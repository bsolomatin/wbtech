package config

import (
	"dockertest/pkg/db"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	postgres.Config
}

type ConsumerConfig struct {
	Host string `env:"KAFKA_HOST" env-default:"localhost"`
	Port string `env:"KAFKA_PORT" env-default:"9092"`
	Topic string `env:"KAFKA_TOPIC" env-default:"test"`
}

func New() *Config {
	config := Config{}
	err := cleanenv.ReadConfig("./configs/local.env", &config)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &config
}