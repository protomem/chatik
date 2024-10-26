package logging

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
	"github.com/protomem/chatik/pkg/werrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

var ErrInvalidLevel = errors.New("invalid log level")

type Extractor func(context.Context) []any

func NopExtractor(context.Context) []any {
	return nil
}

func MultiExtractor(extras ...Extractor) Extractor {
	return func(ctx context.Context) []any {
		var attrs []any
		for _, extra := range extras {
			attrs = append(attrs, extra(ctx)...)
		}
		return attrs
	}
}

type Options struct {
	Level  string `yaml:"level"`
	File   string `yaml:"file"`
	Pretty bool   `yaml:"pretty"`
}

func DefaultOptions() Options {
	return Options{
		Level:  "debug",
		File:   "",
		Pretty: false,
	}
}

type Logger struct {
	opts  Options
	extra Extractor
	*slog.Logger
}

func New(opts Options, extra ...Extractor) (*Logger, error) {
	if len(extra) == 0 {
		extra = []Extractor{NopExtractor}
	}

	lvl, err := ParseLevel(opts.Level)
	if err != nil {
		return nil, werrors.Error(err, "logger/new")
	}

	var out io.Writer = os.Stdout
	if opts.File != "" {
		fsOut := &lumberjack.Logger{
			Filename:   opts.File,
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     7,    // days
			Compress:   true, // disabled by default
		}
		out = io.MultiWriter(out, fsOut)
	}

	var h slog.Handler
	if opts.Pretty {
		o := &tint.Options{Level: lvl}
		h = tint.NewHandler(out, o)
	} else {
		o := &slog.HandlerOptions{Level: lvl}
		h = slog.NewJSONHandler(out, o)
	}

	return &Logger{
		opts:   opts,
		extra:  extra[len(extra)-1],
		Logger: slog.New(h),
	}, nil
}

func NewNop() *Logger {
	return &Logger{
		opts:   DefaultOptions(),
		extra:  NopExtractor,
		Logger: slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{})),
	}
}

func (l *Logger) With(attrs ...any) *Logger {
	return &Logger{
		opts:   l.opts,
		extra:  l.extra,
		Logger: l.Logger.With(attrs...),
	}
}

func (l *Logger) WithExtractor(extra Extractor) *Logger {
	if extra == nil {
		extra = l.extra
	}

	return &Logger{
		opts:   l.opts,
		extra:  extra,
		Logger: l.Logger,
	}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	return l.With(l.extra(ctx)...)
}

func ParseLevel(lvl string) (slog.Level, error) {
	switch strings.ToLower(lvl) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelDebug, ErrInvalidLevel
	}
}

func Group(key string, attrs ...any) slog.Attr {
	return slog.Group(key, attrs...)
}

func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}
