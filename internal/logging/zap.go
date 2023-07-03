package logging

import (
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _ Logger = (*Zap)(nil)

type Zap struct {
	logger *zap.SugaredLogger
}

func NewZap(lvlStr string, file string) (*Zap, error) {
	lvl, err := zapcore.ParseLevel(lvlStr)
	if err != nil {
		return nil, fmt.Errorf("zap.New: parse level: %w", err)
	}

	fsSync := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    50,
		MaxBackups: 3,
		MaxAge:     7,
		LocalTime:  true,
		Compress:   true,
	}

	enc := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	sync := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr), zapcore.AddSync(fsSync))
	core := zapcore.NewCore(enc, sync, lvl)

	logger := zap.New(core).Sugar()

	return &Zap{logger: logger}, nil
}

func (l *Zap) With(args ...any) Logger {
	return &Zap{logger: l.logger.With(args...)}
}

func (l *Zap) Debug(msg string, args ...any) {
	l.logger.Debugw(msg, args...)
}

func (l *Zap) Info(msg string, args ...any) {
	l.logger.Infow(msg, args...)
}

func (l *Zap) Error(msg string, args ...any) {
	l.logger.Errorw(msg, args...)
}

func (l *Zap) Write(p []byte) (int, error) {
	l.Info(string(p))

	return len(p), nil
}

func (l *Zap) Sync() error {
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
