package usecase

import (
	"context"
	"fmt"

	"github.com/protomem/chatik/internal/agregate"
	"github.com/protomem/chatik/internal/database"
	"github.com/protomem/chatik/internal/domain/model"
	"github.com/protomem/chatik/internal/domain/port"
	"github.com/protomem/chatik/internal/stream"
)

var (
	_ port.ObserveCreateChannelUseCase = (*ObserveCreateChannel)(nil)
	_ port.ObserveDeleteChannelUseCase = (*ObserveDeleteChannel)(nil)
	_ port.ObserveCreateMessageUseCase = (*ObserveCreateMessage)(nil)
	_ port.ObserveDeleteMessageUseCase = (*ObserveDeleteMessage)(nil)
)

type ObserveCreateChannel struct {
	stream *stream.Stream

	createChannelUC port.CreateChannelUseCase
}

func NewObserveCreateChannel(stream *stream.Stream, createChannelUC port.CreateChannelUseCase) *ObserveCreateChannel {
	return &ObserveCreateChannel{
		stream:          stream,
		createChannelUC: createChannelUC,
	}
}

func (uc *ObserveCreateChannel) Invoke(ctx context.Context, dto port.CreateChannelUCDTO) (model.Channel, error) {
	const op = "usecase.ObserveCreateChannel"
	var err error

	channel, err := uc.createChannelUC.Invoke(ctx, dto)
	if err != nil {
		return model.Channel{}, fmt.Errorf("%s: %w", op, err)
	}

	_ = stream.SendEvent(uc.stream, stream.NewChannelEvent(agregate.Channel{
		Channel: database.Channel{
			ID:        channel.ID,
			CreatedAt: channel.CreatedAt,
			UpdatedAt: channel.UpdatedAt,
			Title:     channel.Title,
			UserID:    channel.User.ID,
		},
		User: database.User{
			ID:        channel.User.ID,
			CreatedAt: channel.User.CreatedAt,
			UpdatedAt: channel.User.UpdatedAt,
			Nickname:  channel.User.Nickname,
			Email:     channel.User.Email,
			Password:  channel.User.Password,
			Verified:  channel.User.Verified,
		},
	}))

	return channel, nil
}

type ObserveDeleteChannel struct {
	stream *stream.Stream

	deleteChannelUC port.DeleteChannelUseCase
}

func NewObserveDeleteChannel(stream *stream.Stream, deleteChannelUC port.DeleteChannelUseCase) *ObserveDeleteChannel {
	return &ObserveDeleteChannel{
		stream:          stream,
		deleteChannelUC: deleteChannelUC,
	}
}

func (uc *ObserveDeleteChannel) Invoke(ctx context.Context, dto port.DeleteChannelUCDTO) error {
	const op = "usecase.ObserveDeleteChannel"
	var err error

	err = uc.deleteChannelUC.Invoke(ctx, dto)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_ = stream.SendEvent(uc.stream, stream.RemoveChannelEvent(dto.ChannelID))

	return nil
}

type ObserveCreateMessage struct {
	stream *stream.Stream

	createMessageUC port.CreateMessageUseCase
}

func NewObserveCreateMessage(stream *stream.Stream, createMessageUC port.CreateMessageUseCase) *ObserveCreateMessage {
	return &ObserveCreateMessage{
		stream:          stream,
		createMessageUC: createMessageUC,
	}
}

func (uc *ObserveCreateMessage) Invoke(ctx context.Context, dto port.CreateMessageUCDTO) (model.Message, error) {
	const op = "usecase.ObserveCreateMessage"
	var err error

	message, err := uc.createMessageUC.Invoke(ctx, dto)
	if err != nil {
		return model.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	_ = stream.SendEvent(uc.stream, stream.NewMessageEvent(agregate.Message{
		Message: database.Message{
			ID:        message.ID,
			CreatedAt: message.CreatedAt,
			UpdatedAt: message.UpdatedAt,
			Content:   message.Content,
			UserID:    message.User.ID,
			ChannelID: message.ChannelID,
		},
		User: database.User{
			ID:        message.User.ID,
			CreatedAt: message.User.CreatedAt,
			UpdatedAt: message.User.UpdatedAt,
			Nickname:  message.User.Nickname,
			Email:     message.User.Email,
			Password:  message.User.Password,
			Verified:  message.User.Verified,
		},
	}))

	return model.Message{}, nil
}

type ObserveDeleteMessage struct {
	stream *stream.Stream

	deleteMessageUC port.DeleteMessageUseCase
}

func NewObserveDeleteMessage(stream *stream.Stream, deleteMessageUC port.DeleteMessageUseCase) *ObserveDeleteMessage {
	return &ObserveDeleteMessage{
		stream:          stream,
		deleteMessageUC: deleteMessageUC,
	}
}

func (uc *ObserveDeleteMessage) Invoke(ctx context.Context, dto port.DeleteMessageUCDTO) error {
	const op = "usecase.ObserveDeleteMessage"
	var err error

	err = uc.deleteMessageUC.Invoke(ctx, dto)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_ = stream.SendEvent(uc.stream, stream.RemoveMessageEvent(dto.MessageID))

	return nil
}
