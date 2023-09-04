package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidToken = fmt.Errorf("invalid token")

type (
	ctxKey struct{}

	Payload struct {
		UserID   uuid.UUID `json:"user_id"`
		Nickname string    `json:"nickname"`
		Email    string    `json:"email"`
		Verified bool      `json:"email_verified"`
	}

	Claims struct {
		Payload
		jwt.RegisteredClaims
	}
)

func Inject(ctx context.Context, payload Payload) context.Context {
	return context.WithValue(ctx, ctxKey{}, payload)
}

func Extract(ctx context.Context) (Payload, bool) {
	payload, ok := ctx.Value(ctxKey{}).(Payload)
	return payload, ok
}

type GenerateParams struct {
	SigningKey string
	TTL        time.Duration
}

func Generate(payload Payload, params GenerateParams) (string, error) {
	claims := Claims{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(params.TTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "chatik",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(params.SigningKey))
	if err != nil {
		return "", fmt.Errorf("jwt.Generate: %w", err)
	}

	return signedToken, nil
}

type ParseParams struct {
	SigningKey string
}

func Parse(tokenStr string, params ParseParams) (Payload, error) {
	const op = "jwt.Parse"

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(params.SigningKey), nil
	})
	if err != nil {
		return Payload{}, fmt.Errorf("%s: %w", op, err)
	}

	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		return claims.Payload, nil
	}

	return Payload{}, fmt.Errorf("%s: %w", op, ErrInvalidToken)
}
