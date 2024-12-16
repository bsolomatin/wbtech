package config

import (
	"dockertest/pkg/db"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	postgres.Config
	ConsumerConfig
	AppPort string `env:"APP_PORT" env-default:"8080"`
}

type ConsumerConfig struct {
	Host string `env:"KAFKA_HOST" env-default:"localhost"`
	Port string `env:"KAFKA_PORT" env-default:"9092"`
	Topic string `env:"KAFKA_TOPIC" env-default:"test"`
}

func New() (*Config, error){
	config := Config{}
	err := cleanenv.ReadConfig("./configs/local.env", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}