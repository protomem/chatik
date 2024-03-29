package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/protomem/chatik/internal/domain/model"
	"github.com/protomem/chatik/internal/jwt"
)

// Utilities types

type Token = string

// User Use Cases

type (
	FindUserUseCase interface {
		Invoke(ctx context.Context, id uuid.UUID) (model.User, error)
	}
)

type (
	CreateUserUCDTO struct {
		Nickname string
		Email    string
		Password string
	}

	CreateUserUseCase interface {
		Invoke(ctx context.Context, dto CreateUserUCDTO) (model.User, error)
	}
)

type (
	FindUserByEmailAndPasswordUCDTO struct {
		Email    string
		Password string
	}

	FindUserByEmailAndPasswordUseCase interface {
		Invoke(ctx context.Context, dto FindUserByEmailAndPasswordUCDTO) (model.User, error)
	}
)

// Auth Use Cases

type (
	RegisterUserUCDTO struct {
		Nickname string
		Email    string
		Password string
	}

	RegisterUserUseCase interface {
		Invoke(ctx context.Context, dto RegisterUserUCDTO) (model.User, Token, error)
	}
)

type (
	LoginUserUCDTO struct {
		Email    string
		Password string
	}

	LoginUserUseCase interface {
		Invoke(ctx context.Context, dto LoginUserUCDTO) (model.User, Token, error)
	}
)

type (
	VerifyTokenUseCase interface {
		Invoke(ctx context.Context, token string) (model.User, jwt.Payload, error)
	}
)

// Channel Use Cases

type (
	FindChannelUseCase interface {
		Invoke(ctx context.Context, id uuid.UUID) (model.Channel, error)
	}
)

type (
	FindAllChannelsUseCase interface {
		Invoke(ctx context.Context) ([]model.Channel, error)
	}
)

type (
	CreateChannelUCDTO struct {
		Title  string
		UserID uuid.UUID
	}

	CreateChannelUseCase interface {
		Invoke(ctx context.Context, dto CreateChannelUCDTO) (model.Channel, error)
	}
)

type (
	DeleteChannelUCDTO struct {
		ChannelID uuid.UUID
		UserID    uuid.UUID
	}

	DeleteChannelUseCase interface {
		Invoke(ctx context.Context, dto DeleteChannelUCDTO) error
	}
)

// Messages Use Cases

type (
	FindMessageUseCase interface {
		Invoke(ctx context.Context, id uuid.UUID) (model.Message, error)
	}
)

type (
	FindAllMessagesByChannelIDUseCase interface {
		Invoke(ctx context.Context, channelID uuid.UUID) ([]model.Message, error)
	}
)

type (
	CreateMessageUCDTO struct {
		Content   string
		ChannelID uuid.UUID
		UserID    uuid.UUID
	}

	CreateMessageUseCase interface {
		Invoke(ctx context.Context, dto CreateMessageUCDTO) (model.Message, error)
	}
)

type (
	DeleteMessageUCDTO struct {
		MessageID uuid.UUID
		UserID    uuid.UUID
	}

	DeleteMessageUseCase interface {
		Invoke(ctx context.Context, dto DeleteMessageUCDTO) error
	}
)

// Stream Use Cases

type (
	ObserveCreateChannelUseCase = CreateChannelUseCase

	ObserveDeleteChannelUseCase = DeleteChannelUseCase
)

type (
	ObserveCreateMessageUseCase = CreateMessageUseCase

	ObserveDeleteMessageUseCase = DeleteMessageUseCase
)
