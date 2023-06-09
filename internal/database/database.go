package database

import (
	"context"
	"fmt"

	"github.com/protomem/chatik/internal/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	logger logging.Logger
	client *mongo.Client
}

func New(ctx context.Context, logger logging.Logger, uri string) (*DB, error) {
	var (
		err error

		op = "db.New"
	)

	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("%s: connect: %w", op, err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("%s: ping: %w", op, err)
	}

	return &DB{
		logger: logger.With("component", "db"),
		client: client,
	}, nil
}

func (db *DB) Close(ctx context.Context) error {
	err := db.client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("db.Close: %w", err)
	}

	return nil
}
