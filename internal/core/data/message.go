package data

import (
	"context"

	"github.com/protomem/chatik/internal/core/entity"
)

type (
	SelectMessageByFromAndToOptions struct {
		SelectOptions
		From entity.ID
		To   entity.ID
	}

	InsertMessageDTO struct {
		From     entity.ID
		To       entity.ID
		Text     string
		Metadata map[string]string
	}
)

//go:generate mockgen -destination mocks/message.go -package mocks . MessageAccessor,MessageMutator
type (
	MessageAccessor interface {
		SelectByFromAndTo(context.Context, SelectMessageByFromAndToOptions) ([]entity.Message, error)

		Get(context.Context, entity.ID) (entity.Message, error)
	}

	MessageMutator interface {
		Insert(context.Context, InsertMessageDTO) (entity.ID, error)
	}
)
