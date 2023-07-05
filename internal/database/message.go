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

var _ MessageRepository = (*messageRepository)(nil)

var ErrMessageNotFound = errors.New("message(s) not found")

type (
	Message struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`

		Content string `json:"content"`

		UserID    uuid.UUID `json:"-"`
		ChannelID uuid.UUID `json:"channelId"`
	}

	messageDocument struct {
		ID        primitive.ObjectID `bson:"_id"`
		CreatedAt time.Time          `bson:"created_at"`
		UpdatedAt time.Time          `bson:"updated_at"`

		MessageID string `bson:"message_id"`

		Content string `bson:"content"`

		UserID    string `bson:"user_id"`
		ChannelID string `bson:"channel_id"`
	}
)

type (
	CreateMessageDTO struct {
		Content   string
		UserID    uuid.UUID
		ChannelID uuid.UUID
	}
)

type (
	MessageRepository interface {
		CreateIndexes(ctx context.Context) error

		FindByID(ctx context.Context, messageID uuid.UUID) (Message, error)

		FindAllByChannelID(ctx context.Context, channelID uuid.UUID) ([]Message, error)

		Create(ctx context.Context, dto CreateMessageDTO) (uuid.UUID, error)
	}

	messageRepository struct {
		logger logging.Logger
		coll   *mongo.Collection
	}
)

func (db *DB) MessageRepo() *messageRepository {
	return &messageRepository{
		logger: db.logger.With("repo", "message"),
		coll:   db.client.Database("chatik").Collection("messages"),
	}
}

func (repo *messageRepository) CreateIndexes(ctx context.Context) error {
	_, err := repo.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.M{
				"message_id": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return fmt.Errorf("messageRepo.CreateIndexes: %w", err)
	}

	return nil
}

func (repo *messageRepository) FindByID(ctx context.Context, messageID uuid.UUID) (Message, error) {
	var (
		err error

		op = "messageRepo.FindByID"
	)

	filter := bson.D{{Key: "message_id", Value: messageID.String()}}

	res := repo.coll.FindOne(ctx, filter)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return Message{}, fmt.Errorf("%s: %w", op, ErrMessageNotFound)
		}

		return Message{}, fmt.Errorf("%s: %w", op, res.Err())
	}

	doc := messageDocument{}
	err = res.Decode(&doc)
	if err != nil {
		return Message{}, fmt.Errorf("%s: decode: %w", op, err)
	}

	userID, err := uuid.Parse(doc.UserID)
	if err != nil {
		return Message{}, fmt.Errorf("%s: parse user id: %w", op, err)
	}

	channelID, err := uuid.Parse(doc.ChannelID)
	if err != nil {
		return Message{}, fmt.Errorf("%s: parse channel id: %w", op, err)
	}

	message := Message{
		ID:        messageID,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,

		Content: doc.Content,

		UserID:    userID,
		ChannelID: channelID,
	}

	return message, nil
}

func (repo *messageRepository) FindAllByChannelID(ctx context.Context, channelID uuid.UUID) ([]Message, error) {
	var (
		err error

		op = "messageRepo.FindAllByChannelID"
	)

	filter := bson.D{{Key: "channel_id", Value: channelID.String()}}

	res, err := repo.coll.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, ErrMessageNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	doc := []messageDocument{}
	err = res.All(ctx, &doc)
	if err != nil {
		return nil, fmt.Errorf("%s: decode: %w", op, err)
	}

	var messages []Message
	for _, doc := range doc {
		messageID, err := uuid.Parse(doc.MessageID)
		if err != nil {
			return nil, fmt.Errorf("%s: parse message id: %w", op, err)
		}

		userID, err := uuid.Parse(doc.UserID)
		if err != nil {
			return nil, fmt.Errorf("%s: parse user id: %w", op, err)
		}

		channelID, err := uuid.Parse(doc.ChannelID)
		if err != nil {
			return nil, fmt.Errorf("%s: parse channel id: %w", op, err)
		}

		messages = append(messages, Message{
			ID:        messageID,
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
			Content:   doc.Content,
			UserID:    userID,
			ChannelID: channelID,
		})
	}

	if len(messages) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrMessageNotFound)
	}

	return messages, nil
}

func (repo *messageRepository) Create(ctx context.Context, dto CreateMessageDTO) (uuid.UUID, error) {
	var (
		err error

		op = "messageRepo.CreateMessage"
	)

	now := time.Now()
	messageID := uuid.New() // TODO: handle panic?
	doc := messageDocument{
		ID:        primitive.NewObjectID(),
		CreatedAt: now,
		UpdatedAt: now,

		MessageID: messageID.String(),

		Content: dto.Content,

		UserID:    dto.UserID.String(),
		ChannelID: dto.ChannelID.String(),
	}

	// TODO: hadle duplicate key(message_id)
	_, err = repo.coll.InsertOne(ctx, doc)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return messageID, nil
}
