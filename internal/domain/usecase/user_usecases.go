package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/protomem/chatik/internal/domain/model"
	"github.com/protomem/chatik/internal/domain/port"
	"github.com/protomem/chatik/internal/domain/vrule"
	"github.com/protomem/chatik/internal/passhash"
	"github.com/protomem/chatik/pkg/validation"
)

var (
	_ port.FindUserUseCase                   = (*FindUser)(nil)
	_ port.FindUserByEmailAndPasswordUseCase = (*FindUserByEmailAndPassword)(nil)
	_ port.CreateUserUseCase                 = (*CreateUser)(nil)
)

type FindUser struct {
	userRepo port.UserRepository
}

func NewFindUserByID(userRepo port.UserRepository) *FindUser {
	return &FindUser{
		userRepo: userRepo,
	}
}

func (uc *FindUser) Invoke(ctx context.Context, id uuid.UUID) (model.User, error) {
	const op = "usecase.FindUserByID"
	var err error

	user, err := uc.userRepo.FindUserByID(ctx, id)
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
	const op = "usecase.FindUserByEmailAndPassword"
	var err error

	err = validation.Validate(
		vrule.Email(dto.Email),
		vrule.Password(dto.Password),
	)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	user, err := uc.userRepo.FindUserByEmail(ctx, dto.Email)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	err = passhash.Compare(dto.Password, user.Password)
	if err != nil {
		if errors.Is(err, passhash.ErrWrongPassword) {
			return model.User{}, fmt.Errorf("%s: %w", op, model.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

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
		return model.User{}, fmt.Errorf("%s: %w", op, err)
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
