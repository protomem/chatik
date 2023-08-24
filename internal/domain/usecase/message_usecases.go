package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/protomem/chatik/internal/domain/model"
	"github.com/protomem/chatik/internal/domain/port"
	"github.com/protomem/chatik/internal/domain/vrule"
	"github.com/protomem/chatik/pkg/validation"
)

var (
	_ port.FindMessageUseCase                = (*FindMessage)(nil)
	_ port.FindAllMessagesByChannelIDUseCase = (*FindAllMessagesByID)(nil)
	_ port.CreateMessageUseCase              = (*CreateMessage)(nil)
	_ port.DeleteMessageUseCase              = (*DeleteMessage)(nil)
)

type FindMessage struct {
	messageRepo port.MessageRepository
}

func NewFindMessage(messageRepo port.MessageRepository) *FindMessage {
	return &FindMessage{
		messageRepo: messageRepo,
	}
}

func (uc *FindMessage) Invoke(ctx context.Context, id uuid.UUID) (model.Message, error) {
	const op = "usecase.FindMessage"
	var err error

	message, err := uc.messageRepo.FindMessageByID(ctx, id)
	if err != nil {
		return model.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	return message, nil
}

type FindAllMessagesByID struct {
	messageRepo port.MessageRepository
}

func NewFindAllMessagesByID(messageRepo port.MessageRepository) *FindAllMessagesByID {
	return &FindAllMessagesByID{
		messageRepo: messageRepo,
	}
}

func (uc *FindAllMessagesByID) Invoke(ctx context.Context, channelID uuid.UUID) ([]model.Message, error) {
	const op = "usecase.FindAllMessagesByID"
	var err error

	messages, err := uc.messageRepo.FindAllMessagesByChannelID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return messages, nil
}

type CreateMessage struct {
	messageRepo port.MessageRepository

	fincUserUC    port.FindUserUseCase
	findChannelUC port.FindChannelUseCase
}

func NewCreateMessage(
	messageRepo port.MessageRepository,
	fincUserUC port.FindUserUseCase,
	findChannelUC port.FindChannelUseCase,
) *CreateMessage {
	return &CreateMessage{
		messageRepo:   messageRepo,
		fincUserUC:    fincUserUC,
		findChannelUC: findChannelUC,
	}
}

func (uc *CreateMessage) Invoke(ctx context.Context, dto port.CreateMessageUCDTO) (model.Message, error) {
	const op = "usecase.CreateMessage"
	var err error

	err = validation.Validate(vrule.Content(dto.Content))
	if err != nil {
		return model.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = uc.fincUserUC.Invoke(ctx, dto.UserID)
	if err != nil {
		return model.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = uc.findChannelUC.Invoke(ctx, dto.ChannelID)
	if err != nil {
		return model.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	messageID, err := uuid.NewRandom()
	if err != nil {
		return model.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	messageID, err = uc.messageRepo.CreateMessage(ctx, port.CreateMessageRepoDTO{
		MessageID: messageID,
		Content:   dto.Content,
		UserID:    dto.UserID,
		ChannelID: dto.ChannelID,
	})
	if err != nil {
		return model.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	message, err := uc.messageRepo.FindMessageByID(ctx, messageID)
	if err != nil {
		return model.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	return message, nil
}

type DeleteMessage struct {
	messageRepo port.MessageRepository

	findMessageUC port.FindMessageUseCase
}

func NewDeleteMessage(messageRepo port.MessageRepository, findMessageUC port.FindMessageUseCase) *DeleteMessage {
	return &DeleteMessage{
		messageRepo:   messageRepo,
		findMessageUC: findMessageUC,
	}
}

func (uc *DeleteMessage) Invoke(ctx context.Context, dto port.DeleteMessageUCDTO) error {
	const op = "usecase.DeleteMessage"
	var err error

	message, err := uc.findMessageUC.Invoke(ctx, dto.MessageID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if message.User.ID != dto.UserID {
		return fmt.Errorf("%s: %w", op, model.ErrMessageNotFound)
	}

	err = uc.messageRepo.DeleteMessageByID(ctx, dto.MessageID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
