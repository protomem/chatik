package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/protomem/chatik/internal/domain/model"
	"github.com/protomem/chatik/internal/domain/port"
	"github.com/protomem/chatik/internal/jwt"
)

const _tokenTTL = 3 * 24 * time.Hour

var (
	_ port.RegisterUserUseCase = (*RegisterUser)(nil)
	_ port.LoginUserUseCase    = (*LoginUser)(nil)
)

type RegisterUser struct {
	authSecret   string
	createUserUC port.CreateUserUseCase
}

func NewRegisterUser(authSecret string, createUserUC port.CreateUserUseCase) *RegisterUser {
	return &RegisterUser{
		authSecret:   authSecret,
		createUserUC: createUserUC,
	}
}

func (uc *RegisterUser) Invoke(ctx context.Context, dto port.RegisterUserUCDTO) (model.User, port.Token, error) {
	const op = "usecase.RegisterUser"
	var err error

	user, err := uc.createUserUC.Invoke(ctx, port.CreateUserUCDTO(dto))
	if err != nil {
		return model.User{}, "", fmt.Errorf("%s: %w", op, err)
	}

	payload := jwt.Payload{UserID: user.ID, Nickname: user.Nickname, Email: user.Email, Verified: user.Verified}
	params := jwt.GenerateParams{SigningKey: uc.authSecret, TTL: _tokenTTL}
	token, err := jwt.Generate(payload, params)
	if err != nil {
		return model.User{}, "", fmt.Errorf("%s: %w", op, err)
	}

	return user, token, nil
}

type LoginUser struct {
	authSecret                   string
	findUserByEmailAndPasswordUC port.FindUserByEmailAndPasswordUseCase
}

func NewLoginUser(authSecret string, findUserByEmailAndPasswordUC port.FindUserByEmailAndPasswordUseCase) *LoginUser {
	return &LoginUser{
		authSecret:                   authSecret,
		findUserByEmailAndPasswordUC: findUserByEmailAndPasswordUC,
	}
}

func (uc *LoginUser) Invoke(ctx context.Context, dto port.LoginUserUCDTO) (model.User, port.Token, error) {
	const op = "usecase.LoginUser"
	var err error

	user, err := uc.findUserByEmailAndPasswordUC.Invoke(ctx, port.FindUserByEmailAndPasswordUCDTO(dto))
	if err != nil {
		return model.User{}, "", fmt.Errorf("%s: %w", op, err)
	}

	payload := jwt.Payload{UserID: user.ID, Nickname: user.Nickname, Email: user.Email, Verified: user.Verified}
	params := jwt.GenerateParams{SigningKey: uc.authSecret, TTL: _tokenTTL}
	token, err := jwt.Generate(payload, params)
	if err != nil {
		return model.User{}, "", fmt.Errorf("%s: %w", op, err)
	}

	return user, token, nil
}
