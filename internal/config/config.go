package config

import (
	"fmt"

	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTP struct {
		Addr string `env:"ADDR,notEmpty"`
	} `envPrefix:"HTTP_"`

	JWT struct {
		Secret string `env:"SECRET,notEmpty"`
	} `envPrefix:"JWT_"`

	Log struct {
		Level string `env:"LEVEL,notEmpty"`
	} `envPrefix:"LOG_"`

	DB struct {
		URI string `env:"URI,notEmpty"`
	} `envPrefix:"DB_"`
}

func Load(filename string) error {
	err := godotenv.Load(filename)
	if err != nil {
		return fmt.Errorf("config.Load: env: %w", err)
	}

	return nil
}

func Parse() (Config, error) {
	var conf Config

	err := env.Parse(&conf)
	if err != nil {
		return Config{}, fmt.Errorf("config.Parse: env: %w", err)
	}

	return conf, nil
}
