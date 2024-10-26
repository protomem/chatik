package adapter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/protomem/chatik/internal/core/data"
	"github.com/protomem/chatik/internal/core/entity"
	"github.com/protomem/chatik/internal/infra/flashstore"
	"github.com/protomem/chatik/pkg/werrors"
)

var _ data.LastSeen = (*LastSeen)(nil)

type LastSeenOptions struct {
	flashstore.Options
	NumLogs uint
}

func DefaultLastSeenOptions() LastSeenOptions {
	return LastSeenOptions{
		Options: flashstore.Options{
			Addr: "localhost:6379",
			DB:   1,
			Ping: true,
		},
		NumLogs: 2,
	}
}

type LastSeen struct {
	opts  LastSeenOptions
	store *flashstore.Storage
}

func NewLastSeen(ctx context.Context, opts LastSeenOptions) (*LastSeen, error) {
	store, err := flashstore.New(ctx, opts.Options)
	if err != nil {
		return nil, werrors.Error(err, "lastSeen/new")
	}

	return &LastSeen{
		opts:  opts,
		store: store,
	}, nil
}

func (l *LastSeen) Read(ctx context.Context, user entity.ID) ([]entity.LastSeenLog, error) {
	werr := werrors.Wrap("lastSeen", "read")

	if l.opts.NumLogs == 0 {
		return make([]entity.LastSeenLog, 0), nil
	}

	res := l.store.LRange(ctx, l.fmtKey(user.String()), 0, int64(l.opts.NumLogs)-1)
	if err := res.Err(); err != nil {
		return nil, werr(err)
	}

	logs := make([]entity.LastSeenLog, 0, l.opts.NumLogs)
	for _, encodedLog := range res.Val() {
		var log entity.LastSeenLog
		if err := json.Unmarshal([]byte(encodedLog), &log); err != nil {
			return nil, werr(err)
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func (l *LastSeen) Write(ctx context.Context, log entity.LastSeenLog) error {
	werr := werrors.Wrap("lastSeen", "write")

	encodedLog, err := json.Marshal(log)
	if err != nil {
		return werr(err)
	}

	res := l.store.LPush(ctx, l.fmtKey(log.User.String()), encodedLog)
	if err := res.Err(); err != nil {
		return werr(err)
	}

	return nil
}

func (l *LastSeen) Close(ctx context.Context) error {
	err := l.store.Close(ctx)
	return werrors.Error(err, "lastSeen/close")
}

func (l *LastSeen) fmtKey(user string) string {
	return fmt.Sprintf("lastseen:%s", user)
}
