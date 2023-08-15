package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/protomem/chatik/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ UserRepository = (*userRepository)(nil)

var (
	ErrUserNotFound = errors.New("user(s) not found")
	ErrUserExists   = errors.New("user already exists")
)

type (
	User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`

		Nickname string `json:"nickname"`
		Password string `json:"-"`

		Email    string `json:"email"`
		Verified bool   `json:"isVerified"`
	}

	userDocument struct {
		ID        primitive.ObjectID `bson:"_id"`
		CreatedAt time.Time          `bson:"created_at"`
		UpdatedAt time.Time          `bson:"updated_at"`

		UserID string `bson:"user_id"`

		Nickname string `bson:"nickname"`
		Password string `bson:"password"`

		Email    string `bson:"email"`
		Verified bool   `bson:"is_verified"`
	}
)

type (
	CreateUserDTO struct {
		Nickname string
		Password string
		Email    string
	}
)

type (
	UserRepository interface {
		CreateIndexes(ctx context.Context) error

		FindByID(ctx context.Context, userID uuid.UUID) (User, error)
		FindByEmail(ctx context.Context, email string) (User, error)

		FindAllByIDs(ctx context.Context, userIDs []uuid.UUID) ([]User, error)

		Create(ctx context.Context, dto CreateUserDTO) (uuid.UUID, error)
	}

	userRepository struct {
		logger logging.Logger
		coll   *mongo.Collection
	}
)

func (db *DB) UserRepo() *userRepository {
	return &userRepository{
		logger: db.logger.With("repo", "user"),
		coll:   db.client.Database("chatik").Collection("users"),
	}
}

func (repo *userRepository) CreateIndexes(ctx context.Context) error {
	_, err := repo.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
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
		return fmt.Errorf("userRepo.CreateIndexes: %w", err)
	}

	return nil
}

func (repo *userRepository) FindByID(ctx context.Context, userID uuid.UUID) (User, error) {
	var (
		err error

		op = "userRepo.FindByID"
	)

	filter := bson.D{{Key: "user_id", Value: userID.String()}}

	res := repo.coll.FindOne(ctx, filter)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return User{}, fmt.Errorf("%s: %w", op, res.Err())
	}

	doc := userDocument{}
	err = res.Decode(&doc)
	if err != nil {
		return User{}, fmt.Errorf("%s: decode: %w", op, err)
	}

	user := User{
		ID:        userID,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
		Nickname:  doc.Nickname,
		Password:  doc.Password,
		Email:     doc.Email,
		Verified:  doc.Verified,
	}

	return user, nil
}

func (repo *userRepository) FindByEmail(ctx context.Context, email string) (User, error) {
	var (
		err error

		op = "userRepo.FindByEmail"
	)

	filter := bson.D{{Key: "email", Value: email}}

	res := repo.coll.FindOne(ctx, filter)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return User{}, fmt.Errorf("%s: %w", op, res.Err())
	}

	doc := userDocument{}
	err = res.Decode(&doc)
	if err != nil {
		return User{}, fmt.Errorf("%s: decode: %w", op, err)
	}

	userID, err := uuid.Parse(doc.UserID)
	if err != nil {
		return User{}, fmt.Errorf("%s: parse user id: %w", op, err)
	}

	user := User{
		ID:        userID,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
		Nickname:  doc.Nickname,
		Password:  doc.Password,
		Email:     doc.Email,
		Verified:  doc.Verified,
	}

	return user, nil
}

func (repo *userRepository) FindAllByIDs(ctx context.Context, userIDs []uuid.UUID) ([]User, error) {
	var (
		err error

		op = "userRepo.FindAllByIDs"
	)

	userIDsStr := make([]string, len(userIDs))
	for i, userID := range userIDs {
		userIDsStr[i] = userID.String()
	}

	filter := bson.D{
		{Key: "user_id", Value: bson.D{{Key: "$in", Value: userIDsStr}}},
	}

	res, err := repo.coll.Find(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	docs := []userDocument{}
	err = res.All(ctx, &docs)
	if err != nil {
		return nil, fmt.Errorf("%s: decode: %w", op, err)
	}

	users := make([]User, 0, len(docs))
	for _, doc := range docs {
		userID, err := uuid.Parse(doc.UserID)
		if err != nil {
			return nil, fmt.Errorf("%s: parse user id: %w", op, err)
		}

		users = append(users, User{
			ID:        userID,
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
			Nickname:  doc.Nickname,
			Password:  doc.Password,
			Email:     doc.Email,
			Verified:  doc.Verified,
		})
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}

	return users, nil
}

func (repo *userRepository) Create(ctx context.Context, dto CreateUserDTO) (uuid.UUID, error) {
	var (
		err error

		op = "userRepo.Create"
	)

	now := time.Now()
	userID := uuid.New() // TODO: handle panic?
	doc := userDocument{
		ID: primitive.NewObjectID(),

		CreatedAt: now,
		UpdatedAt: now,

		UserID: userID.String(),

		Nickname: dto.Nickname,
		Password: dto.Password,

		Email:    dto.Email,
		Verified: false,
	}

	_, err = repo.coll.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}
