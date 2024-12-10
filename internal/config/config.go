package config

import (
	"dockertest/pkg/db"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	postgres.Config
}

func New() *Config {
	config := Config{}
	pwd, err := os.Getwd()
	fmt.Println("TUT")
	fmt.Println(pwd)
	err = cleanenv.ReadConfig("./configs/local.env", &config)
	if err != nil {
		fmt.Println("WTF TUT")
		fmt.Println(err)
		return nil
	}

	return &config
}