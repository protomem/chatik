package zap

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/protomem/chatik/pkg/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ logging.Logger = (*Logger)(nil)

type Logger struct {
	logger *zap.SugaredLogger
}

func New(lvlStr string) (*Logger, error) {
	lvl, err := zapcore.ParseLevel(lvlStr)
	if err != nil {
		return nil, fmt.Errorf("zap.New: parse level: %w", err)
	}

	enc := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	sync := zapcore.AddSync(os.Stderr)
	core := zapcore.NewCore(enc, sync, lvl)

	logger := zap.New(core).Sugar()

	return &Logger{logger: logger}, nil
}

func (l *Logger) With(args ...any) logging.Logger {
	return &Logger{logger: l.logger.With(args...)}
}

func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debugw(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.logger.Infow(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.logger.Errorw(msg, args...)
}

func (l *Logger) Write(p []byte) (int, error) {
	logStr := string(p)
	logStr = strings.TrimSpace(logStr)
	logStr = strings.Trim(logStr, "\n")

	l.Info(logStr)

	return len(p), nil
}

func (l *Logger) Println(args ...any) {
	l.Debug(fmt.Sprint(args...))
}

func (l *Logger) Sync(_ context.Context) error {
	err := l.logger.Sync()
	if err != nil {
		var pErr *os.PathError
		if errors.As(err, &pErr) && (pErr.Path == "/dev/stderr" || pErr.Path == "/dev/stdout") {
			return nil
		}

		return fmt.Errorf("zap.Sync: %w", err)
	}

	return nil
}
