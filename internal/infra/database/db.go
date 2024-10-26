package database

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/protomem/chatik/internal/infra/logging"
	"github.com/protomem/chatik/pkg/werrors"
)

type (
	Options struct {
		DSN        string `yaml:"dsn"`
		Ping       bool   `yaml:"ping"`
		TraceLevel string `yaml:"traceLevel"`
		MaxConns   int    `yaml:"maxConns"`

		Timeouts Timeouts `yaml:"timeouts"`
	}

	Timeouts struct {
		Idle time.Duration `yaml:"idle"`
		Life time.Duration `yaml:"life"`
	}
)

func DefaultOptions() Options {
	return Options{
		DSN:        "postgres://localhost:5432/postgres?sslmode=disable",
		Ping:       true,
		TraceLevel: "debug",
		MaxConns:   25,
	}
}

type DB struct {
	opts Options
	*pgxpool.Pool
}

func Connect(ctx context.Context, log *logging.Logger, opts Options) (*DB, error) {
	werr := werrors.Wrap("db/connect")

	conf, err := pgxpool.ParseConfig(opts.DSN)
	if err != nil {
		return nil, werr(err, "parse dsn")
	}

	tracer, err := newTracer(log, opts.TraceLevel)
	if err != nil {
		return nil, werr(err, "init tracer")
	}

	conf.ConnConfig.Tracer = tracer
	conf.MaxConns = int32(opts.MaxConns)

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, werr(err)
	}

	if opts.Ping {
		if err := pool.Ping(ctx); err != nil {
			return nil, werr(err, "ping")
		}
	}

	return &DB{
		opts: opts,
		Pool: pool,
	}, nil
}

func (db *DB) QueryBuilder() QueryBuilder {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

func (db *DB) Close(_ context.Context) error {
	db.Pool.Close()
	return nil
}
