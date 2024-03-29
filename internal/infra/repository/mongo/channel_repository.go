package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/protomem/chatik/internal/domain/model"
	"github.com/protomem/chatik/internal/domain/port"
	"github.com/protomem/chatik/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ port.ChannelRepository = (*ChannelRepository)(nil)

type (
	ChannelDocument struct {
		ID primitive.ObjectID `bson:"_id"`

		ChannelID string `bson:"channel_id"`

		CreatedAt time.Time `bson:"created_at"`
		UpdatedAt time.Time `bson:"updated_at"`

		Title string `bson:"title"`

		UserID string `bson:"user_id"`
	}

	ChannelRepository struct {
		logger     logging.Logger
		client     *mongo.Client
		database   string
		collection string
	}
)

func NewChannelRepository(logger logging.Logger, client *mongo.Client) *ChannelRepository {
	return &ChannelRepository{
		logger:     logger.With("repositoryType", "mongo", "repository", "channel"),
		client:     client,
		database:   "chatik",
		collection: "channels",
	}
}

func (r *ChannelRepository) Migrate(ctx context.Context) error {
	const op = "mongo.ChannelRepository.Migrate"
	var err error

	_, err = r.client.
		Database(r.database).
		Collection(r.collection).
		Indexes().
		CreateMany(ctx, []mongo.IndexModel{
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
		})
	if err != nil {
		return fmt.Errorf("%s: create indexes: %w", op, err)
	}

	return nil
}

func (r *ChannelRepository) FindChannelByID(ctx context.Context, channelID uuid.UUID) (model.Channel, error) {
	const op = "mongo.ChannelRepository.FindChannelByID"
	var err error

	channelFilter := bson.D{
		{Key: "channel_id", Value: channelID.String()},
	}

	var channelDoc ChannelDocument
	err = r.client.
		Database(r.database).
		Collection(r.collection).
		FindOne(ctx, channelFilter).
		Decode(&channelDoc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Channel{}, fmt.Errorf("%s: find channel: %w", op, model.ErrChannelNotFound)
		}

		return model.Channel{}, fmt.Errorf("%s: find channel: %w", op, err)
	}

	userFilter := bson.D{
		{Key: "user_id", Value: channelDoc.UserID},
	}

	var userDoc UserDocument
	err = r.client.
		Database(r.database).
		Collection("users").
		FindOne(ctx, userFilter).
		Decode(&userDoc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Channel{}, fmt.Errorf("%s: find user: %w", op, model.ErrChannelNotFound)
		}

		return model.Channel{}, fmt.Errorf("%s: find user: %w", op, err)
	}

	channel, err := mapChannelDocumentAndUserDocumentToChannelModel(channelDoc, userDoc)
	if err != nil {
		return model.Channel{}, fmt.Errorf("%s: %w", op, err)
	}

	return channel, nil
}

func (r *ChannelRepository) FindAllChannels(ctx context.Context) ([]model.Channel, error) {
	const op = "mongo.ChannelRepository.FindAllChannels"
	var err error

	cursorChannels, err := r.client.
		Database(r.database).
		Collection(r.collection).
		Find(ctx, bson.D{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []model.Channel{}, nil
		}

		return []model.Channel{}, fmt.Errorf("%s: %w", op, err)
	}

	channelsDocs := make([]ChannelDocument, 0)
	err = cursorChannels.All(ctx, &channelsDocs)
	if err != nil {
		return []model.Channel{}, fmt.Errorf("%s: decode: %w", op, err)
	}

	userIDs := make([]string, 0, len(channelsDocs))
	for _, doc := range channelsDocs {
		userIDs = append(userIDs, doc.UserID)
	}

	usersFilter := bson.D{
		{Key: "user_id", Value: bson.D{{Key: "$in", Value: userIDs}}},
	}

	cursorUsers, err := r.client.
		Database(r.database).
		Collection("users").
		Find(ctx, usersFilter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []model.Channel{}, nil
		}

		return []model.Channel{}, fmt.Errorf("%s: %w", op, err)
	}

	usersDocs := make([]UserDocument, 0)
	err = cursorUsers.All(ctx, &usersDocs)
	if err != nil {
		return []model.Channel{}, fmt.Errorf("%s: decode: %w", op, err)
	}

	// TODO: debug
	channels := make([]model.Channel, 0, len(channelsDocs))
	for _, doc := range channelsDocs {
		for _, userDoc := range usersDocs {
			if userDoc.UserID == doc.UserID {
				channel, err := mapChannelDocumentAndUserDocumentToChannelModel(doc, userDoc)
				if err != nil {
					return []model.Channel{}, fmt.Errorf("%s: %w", op, err)
				}

				channels = append(channels, channel)
			}
		}
	}

	return channels, nil
}

func (r *ChannelRepository) CreateChannel(ctx context.Context, dto port.CreateChannelRepoDTO) (uuid.UUID, error) {
	const op = "mongo.ChannelRepository.CreateChannel"
	var err error

	now := time.Now()
	doc := ChannelDocument{
		ID:        primitive.NewObjectID(),
		ChannelID: dto.ChannelID.String(),
		CreatedAt: now,
		UpdatedAt: now,
		Title:     dto.Title,
		UserID:    dto.UserID.String(),
	}

	_, err = r.client.
		Database(r.database).
		Collection(r.collection).
		InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return uuid.Nil, fmt.Errorf("%s: create channel: %w", op, model.ErrChannelAlreadyExists)
		}

		return uuid.Nil, fmt.Errorf("%s: create channel: %w", op, err)
	}

	return dto.ChannelID, nil
}

func (r *ChannelRepository) DeleteChannelByID(ctx context.Context, id uuid.UUID) error {
	const op = "mongo.ChannelRepository.DeleteChannelByID"
	var err error

	filter := bson.D{
		{Key: "channel_id", Value: id.String()},
	}

	_, err = r.client.
		Database(r.database).
		Collection(r.collection).
		DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("%s: delete channel: %w", op, err)
	}

	return nil
}

func mapChannelDocumentAndUserDocumentToChannelModel(
	channelDoc ChannelDocument,
	userDoc UserDocument,
) (model.Channel, error) {
	channelID, err := uuid.Parse(channelDoc.ChannelID)
	if err != nil {
		return model.Channel{}, fmt.Errorf("parse channel id: %w", err)
	}

	user, err := mapUserDocumentToUserModel(userDoc)
	if err != nil {
		return model.Channel{}, err
	}

	return model.Channel{
		ID:        channelID,
		CreatedAt: channelDoc.CreatedAt,
		UpdatedAt: channelDoc.UpdatedAt,
		Title:     channelDoc.Title,
		User:      user,
	}, nil
}
