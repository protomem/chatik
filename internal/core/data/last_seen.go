package data

import (
	"context"

	"github.com/protomem/chatik/internal/core/entity"
)

//go:generate mockgen -destination mocks/last_seen.go -package mocks . LastSeen
type LastSeen interface {
	Read(context.Context, entity.ID) ([]entity.LastSeenLog, error)
	Write(context.Context, entity.LastSeenLog) error
}
