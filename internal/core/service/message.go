package service

import (
	"context"

	"github.com/protomem/chatik/internal/core/data"
	"github.com/protomem/chatik/internal/core/entity"
	"github.com/protomem/chatik/pkg/validation"
	"github.com/protomem/chatik/pkg/werrors"
)

var _ Message = (*MessageImpl)(nil)

type (
	FindMessageByFromAndTo struct {
		FindOptions
		From entity.ID
		To   entity.ID
	}

	CreateMessageDTO struct {
		From entity.ID
		To   entity.ID
		Text string
	}
)

//go:generate mockgen -destination mocks/message.go -package mocks . Message
type (
	Message interface {
		FindByFromAndTo(context.Context, FindMessageByFromAndTo) ([]entity.Message, error)

		Create(context.Context, CreateMessageDTO) (entity.Message, error)
	}

	MessageImpl struct {
		acc data.MessageAccessor
		mut data.MessageMutator
	}
)

func NewMessage(acc data.MessageAccessor, mut data.MessageMutator) *MessageImpl {
	return &MessageImpl{
		acc: acc,
		mut: mut,
	}
}

func (serv *MessageImpl) FindByFromAndTo(ctx context.Context, opts FindMessageByFromAndTo) ([]entity.Message, error) {
	werr := werrors.Wrap("service/message", "findByFromAndTo")

	selectOpts := data.SelectMessageByFromAndToOptions{
		SelectOptions: data.SelectOptions(opts.FindOptions),
		From:          opts.From,
		To:            opts.To,
	}

	// TODO: Check exists user

	messages, err := serv.acc.SelectByFromAndTo(ctx, selectOpts)
	if err != nil {
		return []entity.Message{}, werr(err)
	}

	return messages, nil
}

func (serv *MessageImpl) Create(ctx context.Context, dto CreateMessageDTO) (entity.Message, error) {
	werr := werrors.Wrap("service/message", "create")

	if err := dto.Validate(); err != nil {
		return entity.Message{}, werr(err)
	}

	insertDTO := data.InsertMessageDTO{
		From:     dto.From,
		To:       dto.To,
		Text:     dto.Text,
		Metadata: make(map[string]string),
	}

	id, err := serv.mut.Insert(ctx, insertDTO)
	if err != nil {
		return entity.Message{}, werr(err)
	}

	message, err := serv.acc.Get(ctx, id)
	if err != nil {
		return entity.Message{}, werr(err)
	}

	return message, nil
}

func (dto CreateMessageDTO) Validate() error {
	return validation.Validate(func(v *validation.Validator) {
		v.CheckField(validation.NotBlank(dto.Text), "text", "cannot be blank")
		v.CheckField(validation.MaxRunes(dto.Text, 1024), "text", "cannot be more than 1024 characters")
	})
}
