package data

import (
	"context"

	"github.com/protomem/chatik/internal/core/entity"
)

type (
	InsertUserDTO struct {
		Nickname string
		Password string
		Email    string
	}
)

//go:generate mockgen -destination mocks/user.go -package mocks . UserAccessor,UserMutator
type (
	UserAccessor interface {
		Select(context.Context) ([]entity.User, error)
		Get(context.Context, entity.ID) (entity.User, error)
		GetByNickname(context.Context, string) (entity.User, error)
	}

	UserMutator interface {
		Insert(context.Context, InsertUserDTO) (entity.ID, error)
	}
)
