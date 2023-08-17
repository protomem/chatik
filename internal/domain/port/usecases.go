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
	FindUserByIDUseCase interface {
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
	CreateChannelUCDTO struct {
		Title  string
		UserID uuid.UUID
	}

	CreateChannelUseCase interface {
		Invoke(ctx context.Context, dto CreateChannelUCDTO) (model.Channel, error)
	}
)
