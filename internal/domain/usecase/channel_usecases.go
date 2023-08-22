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
	_ port.FindChannelUseCase     = (*FindChannel)(nil)
	_ port.FindAllChannelsUseCase = (*FindAllChannels)(nil)
	_ port.CreateChannelUseCase   = (*CreateChannel)(nil)
	_ port.DeleteChannelUseCase   = (*DeleteChannel)(nil)
)

type FindChannel struct {
	channelRepo port.ChannelRepository
}

func NewFindChannel(channelRepo port.ChannelRepository) *FindChannel {
	return &FindChannel{
		channelRepo: channelRepo,
	}
}

func (uc *FindChannel) Invoke(ctx context.Context, id uuid.UUID) (model.Channel, error) {
	const op = "usecase.FindChannel"
	var err error

	channel, err := uc.channelRepo.FindChannelByID(ctx, id)
	if err != nil {
		return model.Channel{}, fmt.Errorf("%s: %w", op, err)
	}

	return channel, nil
}

type FindAllChannels struct {
	channelRepo port.ChannelRepository
}

func NewFindAllChannels(channelRepo port.ChannelRepository) *FindAllChannels {
	return &FindAllChannels{
		channelRepo: channelRepo,
	}
}

func (uc *FindAllChannels) Invoke(ctx context.Context) ([]model.Channel, error) {
	const op = "usecase.FindAllChannel"
	var err error

	channels, err := uc.channelRepo.FindAllChannels(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return channels, nil
}

type CreateChannel struct {
	channelRepo port.ChannelRepository

	findUserUC port.FindUserUseCase
}

func NewCreateChannel(channelRepo port.ChannelRepository, findUserDUC port.FindUserUseCase) *CreateChannel {
	return &CreateChannel{
		channelRepo: channelRepo,
		findUserUC:  findUserDUC,
	}
}

func (uc *CreateChannel) Invoke(ctx context.Context, dto port.CreateChannelUCDTO) (model.Channel, error) {
	const op = "usecase.CreateChannel"
	var err error

	err = validation.Validate(vrule.Title(dto.Title))
	if err != nil {
		return model.Channel{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = uc.findUserUC.Invoke(ctx, dto.UserID)
	if err != nil {
		return model.Channel{}, fmt.Errorf("%s: %w", op, err)
	}

	channelID, err := uuid.NewRandom()
	if err != nil {
		return model.Channel{}, fmt.Errorf("%s: %w", op, err)
	}

	channelID, err = uc.channelRepo.CreateChannel(ctx, port.CreateChannelRepoDTO{
		ChannelID: channelID,
		Title:     dto.Title,
		UserID:    dto.UserID,
	})
	if err != nil {
		return model.Channel{}, fmt.Errorf("%s: %w", op, err)
	}

	channel, err := uc.channelRepo.FindChannelByID(ctx, channelID)
	if err != nil {
		return model.Channel{}, fmt.Errorf("%s: %w", op, err)
	}

	return channel, nil
}

type DeleteChannel struct {
	channelRepo port.ChannelRepository

	findChannelUC port.FindChannelUseCase
}

func NewDeleteChannel(channelRepo port.ChannelRepository, findChannelUC port.FindChannelUseCase) *DeleteChannel {
	return &DeleteChannel{
		channelRepo:   channelRepo,
		findChannelUC: findChannelUC,
	}
}

func (uc *DeleteChannel) Invoke(ctx context.Context, dto port.DeleteChannelUCDTO) error {
	const op = "usecase.DeleteChannel"
	var err error

	channel, err := uc.findChannelUC.Invoke(ctx, dto.ChannelID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if channel.User.ID != dto.UserID {
		return fmt.Errorf("%s: %w", op, model.ErrChannelNotFound)
	}

	err = uc.channelRepo.DeleteChannelByID(ctx, dto.ChannelID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
