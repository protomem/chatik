package ctxstore

import (
	"context"
	"fmt"
)

func With[T any](ctx context.Context, key Key, value T) context.Context {
	return context.WithValue(ctx, key, value)
}

func From[T any](ctx context.Context, key Key) (T, bool) {
	v, ok := ctx.Value(key).(T)
	return v, ok
}

func MustFrom[T any](ctx context.Context, key Key) T {
	v, ok := ctx.Value(key).(T)
	if !ok {
		failKeyNotFound(key)
	}
	return v
}

func failKeyNotFound(key Key) {
	err := fmt.Errorf("ctxstore: key(%s) not found", key)
	panic(err)
}
