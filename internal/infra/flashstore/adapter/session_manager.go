package adapter

import (
	"context"
	"encoding/json"
	"time"

	"github.com/protomem/chatik/internal/core/data"
	"github.com/protomem/chatik/internal/core/entity"
	"github.com/protomem/chatik/internal/infra/flashstore"
	"github.com/protomem/chatik/pkg/werrors"
	"github.com/redis/go-redis/v9"
)

var _ data.SessionManager = (*SessionManager)(nil)

type SessionManagerOptions struct {
	flashstore.Options
}

func DefaultSessionManagerOptions() SessionManagerOptions {
	return SessionManagerOptions{
		Options: flashstore.Options{
			Addr: "localhost:6379",
			DB:   0,
			Ping: true,
		},
	}
}

type SessionManager struct {
	opts  SessionManagerOptions
	store *flashstore.Storage
}

func NewSessionManager(ctx context.Context, opts SessionManagerOptions) (*SessionManager, error) {
	werr := werrors.Wrap("sessionManager", "new")

	store, err := flashstore.New(ctx, opts.Options)
	if err != nil {
		return nil, werr(err)
	}

	return &SessionManager{
		opts:  opts,
		store: store,
	}, nil
}

func (s *SessionManager) Find(ctx context.Context, token string) (entity.Session, error) {
	werr := werrors.Wrap("sessionManager", "find")

	encodedSession, err := s.store.Get(ctx, s.fmtKey(token)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return entity.Session{}, werr(entity.ErrSessionNotFound)
		}

		return entity.Session{}, werr(err)
	}

	var session entity.Session
	if err := json.Unmarshal(encodedSession, &session); err != nil {
		return entity.Session{}, werr(err)
	}

	return session, nil
}

func (s *SessionManager) Save(ctx context.Context, session entity.Session) error {
	werr := werrors.Wrap("sessionManager", "save")

	encodedSession, err := json.Marshal(session)
	if err != nil {
		return werr(err)
	}

	if err := s.store.Set(
		ctx,
		s.fmtKey(session.Token),
		encodedSession,
		time.Until(session.ExpiredAt),
	).Err(); err != nil {
		return werr(err)
	}

	return nil
}

func (s *SessionManager) Remove(ctx context.Context, token string) error {
	err := s.store.Del(ctx, s.fmtKey(token)).Err()
	return werrors.Error(err, "sessionManager", "save")
}

func (s *SessionManager) Close(ctx context.Context) error {
	err := s.store.Close(ctx)
	return werrors.Error(err, "sessionManager", "close")
}

func (*SessionManager) fmtKey(token string) string {
	return "session:" + token
}
