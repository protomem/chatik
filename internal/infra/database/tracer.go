package database

import (
	"context"

	"github.com/jackc/pgx/v5/tracelog"
	"github.com/protomem/chatik/internal/infra/logging"
)

func newTracer(log *logging.Logger, level string) (*tracelog.TraceLog, error) {
	parsedLevel, err := tracelog.LogLevelFromString(level)
	if err != nil {
		return nil, err
	}

	return &tracelog.TraceLog{
		LogLevel: parsedLevel,
		Logger: tracelog.LoggerFunc(func(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
			log.WithContext(ctx).Debug("QUERY TRACE", "queryLevel", level.String(), "queryMsg", msg, "queryData", data)
		}),
	}, nil
}
