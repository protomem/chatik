package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/protomem/chatik/internal/domain/model"
	"github.com/protomem/chatik/internal/domain/port"
	"github.com/protomem/chatik/internal/validation"
	"github.com/protomem/chatik/internal/validation/vrule"
)

var _ port.CreateChannelUseCase = (*CreateChannel)(nil)

type CreateChannel struct {
	channelRepo port.ChannelRepository

	findUserByIDUC port.FindUserByIDUseCase
}

func NewCreateChannel(channelRepo port.ChannelRepository, findUserByIDUC port.FindUserByIDUseCase) *CreateChannel {
	return &CreateChannel{
		channelRepo:    channelRepo,
		findUserByIDUC: findUserByIDUC,
	}
}

func (uc *CreateChannel) Invoke(ctx context.Context, dto port.CreateChannelUCDTO) (model.Channel, error) {
	const op = "usecase.CreateChannel"
	var err error

	err = validation.Validate(vrule.Title(dto.Title))
	if err != nil {
		return model.Channel{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = uc.findUserByIDUC.Invoke(ctx, dto.UserID)
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
