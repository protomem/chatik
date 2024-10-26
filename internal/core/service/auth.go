package service

import (
	"context"
	"errors"
	"time"

	"github.com/protomem/chatik/internal/core/data"
	"github.com/protomem/chatik/internal/core/entity"
	"github.com/protomem/chatik/pkg/crand"
	"github.com/protomem/chatik/pkg/werrors"
)

var _ Auth = (*AuthImpl)(nil)

type AuthOptions struct {
	AccessToken  JWTTokenOptions    `yaml:"accessToken"`
	RefreshToken RandomTokenOptions `yaml:"refreshToken"`
}

func DefaultAuthOptions() AuthOptions {
	return AuthOptions{
		AccessToken: JWTTokenOptions{
			SigningKey: crand.String(64),
			Expiration: 3 * time.Hour,
		},
		RefreshToken: RandomTokenOptions{
			Expiration: 3 * 24 * time.Hour,
		},
	}
}

type (
	RegisterDTO struct {
		CreateUserDTO
		Fingerprint string
	}

	LoginDTO struct {
		Login       string
		Password    string
		Fingerprint string
	}

	RefreshDTO struct {
		Token       string
		Fingerprint string
	}

	TokensDTO struct {
		Access  string `json:"accessToken"`
		Refresh string `json:"refreshToken"`
	}
)

//go:generate mockgen -destination mocks/auth.go -package mocks . Auth
type (
	Auth interface {
		Register(context.Context, RegisterDTO) (entity.User, TokensDTO, error)
		Login(context.Context, LoginDTO) (entity.User, TokensDTO, error)
		Logout(context.Context, string) error

		Refresh(context.Context, RefreshDTO) (TokensDTO, error)
		Verify(context.Context, string) (entity.User, error)
	}

	AuthImpl struct {
		opts    AuthOptions
		user    User
		token   Token
		sessmng data.SessionManager
	}
)

func NewAuth(opts AuthOptions, user User, token Token, sessmng data.SessionManager) *AuthImpl {
	return &AuthImpl{
		opts:    opts,
		user:    user,
		token:   token,
		sessmng: sessmng,
	}
}

func (serv *AuthImpl) Register(ctx context.Context, dto RegisterDTO) (entity.User, TokensDTO, error) {
	werr := werrors.Wrap("service/auth", "register")

	user, err := serv.user.Create(ctx, dto.CreateUserDTO)
	if err != nil {
		return entity.User{}, TokensDTO{}, err
	}

	tokens, err := serv.genTokens(user, dto.Fingerprint)
	if err != nil {
		return entity.User{}, TokensDTO{}, werr(err)
	}

	if err := serv.saveSession(ctx, user, tokens.Refresh, serv.opts.RefreshToken.Expiration, dto.Fingerprint); err != nil {
		return entity.User{}, TokensDTO{}, werr(err)
	}

	return user, tokens, nil
}

func (serv *AuthImpl) Login(ctx context.Context, dto LoginDTO) (entity.User, TokensDTO, error) {
	werr := werrors.Wrap("service/auth", "login")

	user, err := serv.user.GetByNicknameAndPassword(ctx, dto.Login, dto.Password)
	if err != nil {
		return entity.User{}, TokensDTO{}, err
	}

	tokens, err := serv.genTokens(user, dto.Fingerprint)
	if err != nil {
		return entity.User{}, TokensDTO{}, werr(err)
	}

	if err := serv.saveSession(ctx, user, tokens.Refresh, serv.opts.RefreshToken.Expiration, dto.Fingerprint); err != nil {
		return entity.User{}, TokensDTO{}, werr(err)
	}

	return user, tokens, nil
}

func (serv *AuthImpl) Logout(ctx context.Context, token string) error {
	err := serv.sessmng.Remove(ctx, token)
	return werrors.Error(err, "service/auth", "logout")
}

func (serv *AuthImpl) Refresh(ctx context.Context, dto RefreshDTO) (TokensDTO, error) {
	werr := werrors.Wrap("service/auth", "refresh")

	session, err := serv.takeSession(ctx, dto.Token)
	if err != nil {
		return TokensDTO{}, werr(err)
	}

	user, err := serv.user.Get(ctx, session.User)
	if err != nil {
		return TokensDTO{}, werr(err)
	}

	tokens, err := serv.genTokens(user, dto.Fingerprint)
	if err != nil {
		return TokensDTO{}, werr(err)
	}

	if err := serv.saveSession(ctx, user, tokens.Refresh, serv.opts.RefreshToken.Expiration, dto.Fingerprint); err != nil {
		return TokensDTO{}, werr(err)
	}

	return tokens, nil
}

func (serv *AuthImpl) Verify(ctx context.Context, token string) (entity.User, error) {
	werr := werrors.Wrap("service/auth", "verify")

	payload, err := serv.token.VerifyJWT(token, serv.opts.AccessToken)
	if err != nil {
		return entity.User{}, werr(err)
	}

	user, err := serv.user.Get(ctx, payload.User)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return entity.User{}, werr(ErrInvalidToken)
		}

		return entity.User{}, werr(err)
	}

	return user, nil
}

func (serv *AuthImpl) genTokens(user entity.User, fingerprint string) (TokensDTO, error) {
	var (
		err  error
		werr = werrors.Wrap("genTokens")

		tokens TokensDTO
	)

	payload := JWTTokenPayload{User: user.ID, Fingerprint: fingerprint, Email: user.Email, Verified: user.Verified}
	tokens.Access, err = serv.token.GenerateJWT(payload, serv.opts.AccessToken)
	if err != nil {
		return TokensDTO{}, werr(err)
	}

	tokens.Refresh, err = serv.token.GenerateRandom(serv.opts.RefreshToken)
	if err != nil {
		return TokensDTO{}, werr(err)
	}

	return tokens, nil
}

func (serv *AuthImpl) takeSession(ctx context.Context, token string) (entity.Session, error) {
	werr := werrors.Wrap("takeSession")

	session, err := serv.sessmng.Find(ctx, token)
	if err != nil {
		return entity.Session{}, werr(err)
	}

	if err := serv.sessmng.Remove(ctx, session.Token); err != nil {
		return entity.Session{}, werr(err)
	}

	if session.ExpiredAt.Before(time.Now()) {
		return entity.Session{}, werr(entity.ErrSessionNotFound)
	}

	return session, nil
}

func (serv *AuthImpl) saveSession(
	ctx context.Context,
	user entity.User,
	token string,
	expiration time.Duration,
	fingerprint string,
) error {
	now := time.Now()
	session := entity.Session{
		CreatedAt:   now,
		ExpiredAt:   now.Add(expiration),
		User:        user.ID,
		Token:       token,
		Fingerprint: fingerprint,
	}

	if err := serv.sessmng.Save(ctx, session); err != nil {
		return werrors.Error(err, "saveSession")
	}

	return nil
}
