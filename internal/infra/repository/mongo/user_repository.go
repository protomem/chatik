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

var _ port.UserRepository = (*UserRepository)(nil)

type (
	UserDocument struct {
		ID primitive.ObjectID `bson:"_id"`

		UserID string `bson:"user_id"`

		CreatedAt time.Time `bson:"created_at"`
		UpdatedAt time.Time `bson:"updated_at"`

		Nickname string `bson:"nickname"`
		Password string `bson:"password"`

		Email    string `bson:"email"`
		Verified bool   `bson:"is_verified"`
	}

	UserRepository struct {
		logger     logging.Logger
		client     *mongo.Client
		database   string
		collection string
	}
)

func NewUserRepository(logger logging.Logger, client *mongo.Client) *UserRepository {
	return &UserRepository{
		logger:     logger.With("repositoryType", "mongo", "repository", "user"),
		client:     client,
		database:   "chatik",
		collection: "users",
	}
}

func (r *UserRepository) Migrate(ctx context.Context) error {
	const op = "mongo.UserRepository.Migrate"
	var err error

	_, err = r.client.
		Database(r.database).
		Collection(r.collection).
		Indexes().
		CreateMany(ctx, []mongo.IndexModel{
			{
				Keys: bson.M{
					"user_id": 1,
				},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys: bson.M{
					"nickname": 1,
				},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys: bson.M{
					"email": 1,
				},
				Options: options.Index().SetUnique(true),
			},
		})
	if err != nil {
		return fmt.Errorf("%s: create indexes: %w", op, err)
	}

	return nil
}

func (r *UserRepository) FindUserByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	const op = "mongo.UserRepository.FindUserByID"
	var err error

	filter := bson.D{
		{Key: "user_id", Value: id.String()},
	}

	var userDoc UserDocument
	err = r.client.
		Database(r.database).
		Collection(r.collection).
		FindOne(ctx, filter).
		Decode(&userDoc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	user, err := mapUserDocumentToUserModel(userDoc)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (model.User, error) {
	const op = "mongo.UserRepository.FindUserByEmail"
	var err error

	filter := bson.D{
		{Key: "email", Value: email},
	}

	var userDoc UserDocument
	err = r.client.
		Database(r.database).
		Collection(r.collection).
		FindOne(ctx, filter).
		Decode(&userDoc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	user, err := mapUserDocumentToUserModel(userDoc)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, dto port.CreateUserRepoDTO) (uuid.UUID, error) {
	const op = "mongo.UserRepository.CreateUser"
	var err error

	now := time.Now()
	doc := UserDocument{
		ID:        primitive.NewObjectID(),
		UserID:    dto.UserID.String(),
		CreatedAt: now,
		UpdatedAt: now,
		Nickname:  dto.Nickname,
		Password:  dto.Password,
		Email:     dto.Email,
		Verified:  false,
	}

	_, err = r.client.
		Database(r.database).
		Collection(r.collection).
		InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, model.ErrUserAlreadyExists)
		}

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return dto.UserID, nil
}

func mapUserDocumentToUserModel(doc UserDocument) (model.User, error) {
	userID, err := uuid.Parse(doc.UserID)
	if err != nil {
		return model.User{}, fmt.Errorf("parse user id: %w", err)
	}

	return model.User{
		ID:        userID,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
		Nickname:  doc.Nickname,
		Password:  doc.Password,
		Email:     doc.Email,
		Verified:  doc.Verified,
	}, nil
}
