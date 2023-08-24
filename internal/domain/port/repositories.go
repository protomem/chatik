package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/protomem/chatik/internal/domain/model"
)

type (
	CreateUserRepoDTO struct {
		UserID   uuid.UUID
		Nickname string
		Email    string
		Password string
	}

	UserRepository interface {
		FindUserByID(ctx context.Context, id uuid.UUID) (model.User, error)
		FindUserByEmail(ctx context.Context, email string) (model.User, error)
		CreateUser(ctx context.Context, dto CreateUserRepoDTO) (uuid.UUID, error)
	}
)

type (
	CreateChannelRepoDTO struct {
		ChannelID uuid.UUID
		Title     string
		UserID    uuid.UUID
	}

	ChannelRepository interface {
		FindChannelByID(ctx context.Context, id uuid.UUID) (model.Channel, error)
		FindAllChannels(ctx context.Context) ([]model.Channel, error)
		CreateChannel(ctx context.Context, dto CreateChannelRepoDTO) (uuid.UUID, error)
		DeleteChannelByID(ctx context.Context, id uuid.UUID) error
	}
)

type (
	CreateMessageRepoDTO struct {
		MessageID uuid.UUID
		Content   string
		ChannelID uuid.UUID
		UserID    uuid.UUID
	}

	MessageRepository interface {
		FindMessageByID(ctx context.Context, id uuid.UUID) (model.Message, error)
		FindAllMessagesByChannelID(ctx context.Context, channelID uuid.UUID) ([]model.Message, error)
		CreateMessage(ctx context.Context, dto CreateMessageRepoDTO) (uuid.UUID, error)
		DeleteMessageByID(ctx context.Context, id uuid.UUID) error
	}
)
