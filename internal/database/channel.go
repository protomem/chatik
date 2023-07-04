package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/protomem/chatik/internal/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ ChannelRepository = (*channelRepository)(nil)

var (
	ErrChannelNotFound = errors.New("channel not found")
	ErrChannelExists   = errors.New("channel already exists")
)

type (
	Channel struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`

		Title string `json:"title"`

		UserID uuid.UUID `json:"-"`
	}

	channelDocument struct {
		ID        primitive.ObjectID `bson:"_id"`
		CreatedAt time.Time          `bson:"created_at"`
		UpdatedAt time.Time          `bson:"updated_at"`

		ChannelID string `bson:"channel_id"`

		Title string `bson:"title"`

		UserID string `bson:"user_id"`
	}
)

type (
	CreateChannelDTO struct {
		Title  string
		UserID uuid.UUID
	}
)

type (
	ChannelRepository interface {
		CreateIndexes(ctx context.Context) error

		Create(ctx context.Context, dto CreateChannelDTO) (uuid.UUID, error)
	}

	channelRepository struct {
		logger logging.Logger
		coll   *mongo.Collection
	}
)

func (db *DB) ChannelRepo() *channelRepository {
	return &channelRepository{
		logger: db.logger.With("repo", "channel"),
		coll:   db.client.Database("chatik").Collection("channels"),
	}
}

func (repo *channelRepository) CreateIndexes(ctx context.Context) error {
	_, err := repo.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.M{
				"channel_id": 1,
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{
				"title": 1,
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{
				"user_id": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return fmt.Errorf("channelRepo.CreateIndexes: %w", err)
	}

	return nil
}

func (repo *channelRepository) FindByID(ctx context.Context, channelID uuid.UUID) (Channel, error) {
	var (
		err error

		op = "channelRepo.FindByID"
	)

	filter := bson.D{{Key: "channel_id", Value: channelID.String()}}

	res := repo.coll.FindOne(ctx, filter)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return Channel{}, fmt.Errorf("%s: %w", op, ErrChannelNotFound)
		}

		return Channel{}, fmt.Errorf("%s: %w", op, res.Err())
	}

	doc := channelDocument{}
	err = res.Decode(&doc)
	if err != nil {
		return Channel{}, fmt.Errorf("%s: decode: %w", op, err)
	}

	userID, err := uuid.Parse(doc.UserID)
	if err != nil {
		return Channel{}, fmt.Errorf("%s: parse user id: %w", op, err)
	}

	channel := Channel{
		ID:        channelID,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
		Title:     doc.Title,
		UserID:    userID,
	}

	return channel, nil
}

func (repo *channelRepository) Create(ctx context.Context, dto CreateChannelDTO) (uuid.UUID, error) {
	var (
		err error

		op = "channelRepo.Create"
	)

	now := time.Now()
	channelID := uuid.New() // TODO: handle panic?
	doc := channelDocument{
		ID: primitive.NewObjectID(),

		CreatedAt: now,
		UpdatedAt: now,

		ChannelID: channelID.String(),

		Title: dto.Title,

		UserID: dto.UserID.String(),
	}

	_, err = repo.coll.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, ErrChannelExists)
		}

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return channelID, nil
}
