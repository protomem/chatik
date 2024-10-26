package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/protomem/chatik/internal/core/entity"
	"github.com/protomem/chatik/pkg/werrors"
)

var _ Token = (*TokenImpl)(nil)

var ErrInvalidToken = errors.New("invalid token")

type (
	RandomTokenOptions struct {
		Expiration time.Duration `yaml:"expiration"`
	}

	JWTTokenOptions struct {
		SigningKey string        `yaml:"secretKey"`
		Expiration time.Duration `yaml:"expiration"`
	}

	JWTTokenPayload struct {
		User        entity.ID `json:"user_id"`
		Fingerprint string    `json:"fingerprint"`
		Email       string    `json:"email"`
		Verified    bool      `json:"email_verified"`
	}

	jwtTokenClaims struct {
		JWTTokenPayload
		jwt.RegisteredClaims
	}
)

//go:generate mockgen -destination mocks/token.go -package mocks . Token
type (
	Token interface {
		GenerateRandom(RandomTokenOptions) (string, error)
		GenerateJWT(JWTTokenPayload, JWTTokenOptions) (string, error)

		VerifyJWT(string, JWTTokenOptions) (JWTTokenPayload, error)
	}

	TokenImpl struct{}
)

func NewToken() *TokenImpl {
	return &TokenImpl{}
}

func (*TokenImpl) GenerateRandom(_ RandomTokenOptions) (string, error) {
	return entity.GenID().String(), nil
}

func (*TokenImpl) GenerateJWT(payload JWTTokenPayload, opts JWTTokenOptions) (string, error) {
	claims := jwtTokenClaims{
		JWTTokenPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(opts.Expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(opts.SigningKey))
	if err != nil {
		return "", werrors.Error(err, "service/token", "generateJWT")
	}

	return signedToken, nil
}

func (*TokenImpl) VerifyJWT(signedToken string, opts JWTTokenOptions) (JWTTokenPayload, error) {
	werr := werrors.Wrap("service/token", "verifyJWT")

	token, err := jwt.ParseWithClaims(signedToken, &jwtTokenClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(opts.SigningKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return JWTTokenPayload{}, werr(ErrInvalidToken)
		}

		return JWTTokenPayload{}, werr(err)
	}

	claims, ok := token.Claims.(*jwtTokenClaims)
	if ok && token.Valid {
		return claims.JWTTokenPayload, nil
	}

	return JWTTokenPayload{}, werr(ErrInvalidToken)
}
