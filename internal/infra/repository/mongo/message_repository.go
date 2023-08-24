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

var _ port.MessageRepository = (*MessageRepository)(nil)

type (
	MessageDocument struct {
		ID primitive.ObjectID `bson:"_id"`

		MessageID string `bson:"message_id"`

		CreatedAt time.Time `bson:"created_at"`
		UpdatedAt time.Time `bson:"updated_at"`

		Content string `bson:"content"`

		UserID    string `bson:"user_id"`
		ChannelID string `bson:"channel_id"`
	}

	MessageRepository struct {
		logger     logging.Logger
		client     *mongo.Client
		database   string
		collection string
	}
)

func NewMessageRepository(logger logging.Logger, client *mongo.Client) *MessageRepository {
	return &MessageRepository{
		logger:     logger.With("repositoryType", "mongo", "repository", "message"),
		client:     client,
		database:   "chatik",
		collection: "messages",
	}
}

func (r *MessageRepository) Migrate(ctx context.Context) error {
	const op = "mongo.MessageRepository.Migrate"
	var err error

	_, err = r.client.
		Database(r.database).
		Collection(r.collection).
		Indexes().
		CreateMany(ctx, []mongo.IndexModel{
			{
				Keys: bson.D{
					{Key: "message_id", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		})
	if err != nil {
		return fmt.Errorf("%s: create indexes: %w", op, err)
	}

	return nil
}

func (r *MessageRepository) FindMessageByID(ctx context.Context, id uuid.UUID) (model.Message, error) {
	const op = "mongo.MessageRepository.FindMessageByID"
	var err error

	messageFilter := bson.D{
		{Key: "message_id", Value: id.String()},
	}

	var messageDoc MessageDocument
	err = r.client.
		Database(r.database).
		Collection(r.collection).
		FindOne(ctx, messageFilter).
		Decode(&messageDoc)
	if err != nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return model.Message{}, fmt.Errorf("%s: find message: %w", op, model.ErrMessageNotFound)
		}

		return model.Message{}, fmt.Errorf("%s: find message: %w", op, err)
	}

	userFilter := bson.D{
		{Key: "user_id", Value: messageDoc.UserID},
	}

	var userDoc UserDocument
	err = r.client.
		Database(r.database).
		Collection("users").
		FindOne(ctx, userFilter).
		Decode(&userDoc)
	if err != nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return model.Message{}, fmt.Errorf("%s: find user: %w", op, model.ErrMessageNotFound)
		}

		return model.Message{}, fmt.Errorf("%s: find user: %w", op, err)
	}

	message, err := mapMessageDocumentAndUserDocumentToMessageModel(messageDoc, userDoc)
	if err != nil {
		return model.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	return message, nil
}

func (r *MessageRepository) FindAllMessagesByChannelID(ctx context.Context, channelID uuid.UUID) ([]model.Message, error) {
	const op = "mongo.MessageRepository.FindAllMessagesByChannelID"
	var err error

	messagesFilter := bson.D{
		{Key: "channel_id", Value: channelID.String()},
	}

	messagesCursor, err := r.client.
		Database(r.database).
		Collection(r.collection).
		Find(ctx, messagesFilter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []model.Message{}, nil
		}

		return []model.Message{}, fmt.Errorf("%s: find messages: %w", op, err)
	}

	messagesDocs := make([]MessageDocument, 0)
	err = messagesCursor.All(ctx, &messagesDocs)
	if err != nil {
		return []model.Message{}, fmt.Errorf("%s: decode: %w", op, err)
	}

	userIDs := make([]string, 0, len(messagesDocs))
	for _, doc := range messagesDocs {
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
			return []model.Message{}, nil
		}

		return []model.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	usersDocs := make([]UserDocument, 0)
	err = cursorUsers.All(ctx, &usersDocs)
	if err != nil {
		return []model.Message{}, fmt.Errorf("%s: decode: %w", op, err)
	}

	messages := make([]model.Message, 0, len(messagesDocs))
	for _, doc := range messagesDocs {
		for _, userDoc := range usersDocs {
			if userDoc.UserID == doc.UserID {
				message, err := mapMessageDocumentAndUserDocumentToMessageModel(doc, userDoc)
				if err != nil {
					return []model.Message{}, fmt.Errorf("%s: %w", op, err)
				}

				messages = append(messages, message)
			}
		}
	}

	return messages, nil
}

func (r *MessageRepository) CreateMessage(ctx context.Context, dto port.CreateMessageRepoDTO) (uuid.UUID, error) {
	const op = "mongo.MessageRepository.CreateMessage"
	var err error

	now := time.Now()
	doc := MessageDocument{
		ID:        primitive.NewObjectID(),
		MessageID: dto.MessageID.String(),
		CreatedAt: now,
		UpdatedAt: now,
		Content:   dto.Content,
		UserID:    dto.UserID.String(),
		ChannelID: dto.ChannelID.String(),
	}

	_, err = r.client.
		Database(r.database).
		Collection(r.collection).
		InsertOne(ctx, doc)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: create message: %w", op, err)
	}

	return dto.MessageID, nil
}

func (r *MessageRepository) DeleteMessageByID(ctx context.Context, id uuid.UUID) error {
	const op = "mongo.MessageRepository.DeleteMessage"
	var err error

	filter := bson.D{
		{Key: "message_id", Value: id.String()},
	}

	_, err = r.client.
		Database(r.database).
		Collection(r.collection).
		DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func mapMessageDocumentAndUserDocumentToMessageModel(
	messageDoc MessageDocument,
	userDoc UserDocument,
) (model.Message, error) {
	user, err := mapUserDocumentToUserModel(userDoc)
	if err != nil {
		return model.Message{}, err
	}

	messageID, err := uuid.Parse(messageDoc.MessageID)
	if err != nil {
		return model.Message{}, fmt.Errorf("parse message id: %w", err)
	}

	channelID, err := uuid.Parse(messageDoc.ChannelID)
	if err != nil {
		return model.Message{}, fmt.Errorf("parse channel id: %w", err)
	}

	return model.Message{
		ID:        messageID,
		CreatedAt: messageDoc.CreatedAt,
		UpdatedAt: messageDoc.UpdatedAt,
		Content:   messageDoc.Content,
		ChannelID: channelID,
		User:      user,
	}, nil
}
