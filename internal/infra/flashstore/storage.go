package flashstore

import (
	"context"

	"github.com/protomem/chatik/pkg/werrors"
	"github.com/redis/go-redis/v9"
)

type Options struct {
	Addr string `yaml:"addr"`
	DB   int    `yaml:"db"`
	Ping bool   `yaml:"ping"`
}

type Storage struct {
	opts Options
	*redis.Client
}

func New(ctx context.Context, opts Options) (*Storage, error) {
	client := redis.NewClient(&redis.Options{
		Addr: opts.Addr,
		DB:   opts.DB,
	})

	if opts.Ping {
		if err := client.Ping(ctx).Err(); err != nil {
			return nil, werrors.Error(err, "flashtore/new")
		}
	}

	return &Storage{
		opts:   opts,
		Client: client,
	}, nil
}

func (s *Storage) Close(_ context.Context) error {
	err := s.Client.Close()
	return werrors.Error(err, "flashtore/close")
}
