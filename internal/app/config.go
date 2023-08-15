package app

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP struct {
		Addr string `yaml:"addr" env-required:"true"`
	} `yaml:"http"`

	Log struct {
		Level string `yaml:"level" env-required:"true"`
	} `yaml:"log"`

	Auth struct {
		Secret string `yaml:"secret" env-required:"true"`
	} `yaml:"auth"`

	Mongo struct {
		URI      string `yaml:"uri" env-required:"true"`
		User     string `yaml:"user" env-required:"true"`
		Password string `yaml:"password" env-required:"true"`
	} `yaml:"mongo"`
}

func NewConfig(file string) (Config, error) {
	var conf Config

	err := cleanenv.ReadConfig(file, &conf)
	if err != nil {
		return Config{}, fmt.Errorf("app.NewConfig: %w", err)
	}

	return conf, nil
}
