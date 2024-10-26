package config

import (
	"os"

	"github.com/protomem/chatik/internal/core/service"
	"github.com/protomem/chatik/internal/infra/database"
	flashAdapter "github.com/protomem/chatik/internal/infra/flashstore/adapter"
	"github.com/protomem/chatik/internal/infra/logging"
	"github.com/protomem/chatik/internal/infra/transport"
	"github.com/protomem/chatik/pkg/werrors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger logging.Options `yaml:"log"`

	HttpServer transport.ServerOptions `yaml:"http"`
	DB         database.Options        `yaml:"db"`

	SessionManager flashAdapter.SessionManagerOptions `yaml:"sessionMng"`
	LastSeen       flashAdapter.LastSeenOptions       `yaml:"lastSeen"`

	Auth service.AuthOptions `yaml:"auth"`
}

func New(filename string) (Config, error) {
	werr := werrors.Wrap("config/new")

	file, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, werr(err, "read file")
	}

	var conf Config
	conf.Logger = logging.DefaultOptions()
	conf.HttpServer = transport.DefaultServerOptions()
	conf.DB = database.DefaultOptions()
	conf.SessionManager = flashAdapter.DefaultSessionManagerOptions()
	conf.LastSeen = flashAdapter.DefaultLastSeenOptions()
	conf.Auth = service.DefaultAuthOptions()

	if err := yaml.Unmarshal(file, &conf); err != nil {
		return Config{}, werr(err, "unmarshal yaml")
	}

	return conf, nil
}
