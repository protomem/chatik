package port

import (
	"context"

	"github.com/protomem/chatik/internal/domain/model"
)

// Utilities types

type Token = string

// User Use Cases

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
