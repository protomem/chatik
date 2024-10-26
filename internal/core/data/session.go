package data

import (
	"context"

	"github.com/protomem/chatik/internal/core/entity"
)

//go:generate mockgen -destination mocks/session.go -package mocks . SessionManager
type SessionManager interface {
	Find(context.Context, string) (entity.Session, error)
	Save(context.Context, entity.Session) error
	Remove(context.Context, string) error
}
