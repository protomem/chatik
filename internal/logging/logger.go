package logging

import "io"

type Logger interface {
	With(args ...any) Logger

	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Error(msg string, args ...any)

	io.Writer

	Sync() error
}
