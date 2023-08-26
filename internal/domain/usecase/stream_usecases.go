package usecase

import (
	"context"
	"fmt"

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

	_ = stream.SendEvent(uc.stream, stream.NewChannelEvent(channel))

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

	_ = stream.SendEvent(uc.stream, stream.NewMessageEvent(message))

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
