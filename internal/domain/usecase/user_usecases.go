package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/protomem/chatik/internal/domain/model"
	"github.com/protomem/chatik/internal/domain/port"
	"github.com/protomem/chatik/internal/passhash"
	"github.com/protomem/chatik/internal/validation"
	"github.com/protomem/chatik/internal/validation/vrule"
)

var (
	_ port.CreateUserUseCase                 = (*CreateUser)(nil)
	_ port.FindUserByEmailAndPasswordUseCase = (*FindUserByEmailAndPassword)(nil)
)

type CreateUser struct {
	userRepo port.UserRepository
}

func NewCreateUser(userRepo port.UserRepository) *CreateUser {
	return &CreateUser{
		userRepo: userRepo,
	}
}

func (uc *CreateUser) Invoke(ctx context.Context, dto port.CreateUserUCDTO) (model.User, error) {
	const op = "usecase.CreateUser"
	var err error

	err = validation.Validate(
		vrule.Nickname(dto.Nickname),
		vrule.Email(dto.Email),
		vrule.Password(dto.Password),
	)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	hashPass, err := passhash.Generate(dto.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := uuid.NewRandom()
	if err != nil {
		return model.User{}, fmt.Errorf("%s: genete uuid: %w", op, err)
	}

	userID, err = uc.userRepo.CreateUser(ctx, port.CreateUserRepoDTO{
		UserID:   userID,
		Nickname: dto.Nickname,
		Email:    dto.Email,
		Password: hashPass,
	})
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	user, err := uc.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

type FindUserByEmailAndPassword struct {
	userRepo port.UserRepository
}

func NewFindUserByEmailAndPassword(userRepo port.UserRepository) *FindUserByEmailAndPassword {
	return &FindUserByEmailAndPassword{
		userRepo: userRepo,
	}
}

func (uc *FindUserByEmailAndPassword) Invoke(
	ctx context.Context,
	dto port.FindUserByEmailAndPasswordUCDTO,
) (model.User, error) {
	return model.User{}, nil
}
