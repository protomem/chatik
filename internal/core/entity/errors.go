package entity

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotFound = errors.New("not found")
	ErrExists   = errors.New("already exists")
)

var (
	ErrUserNotFound = NewError("user", ErrNotFound)
	ErrUserExists   = NewError("user", ErrExists)

	ErrSessionNotFound = NewError("session", ErrNotFound)

	ErrMessageNotFound = NewError("message", ErrNotFound)
)

type Error struct {
	Entity string
	Field  string
	Err    error
}

func NewError(entity string, err error) Error {
	return Error{Entity: strings.ToLower(entity), Err: err}
}

func (err Error) WithErr(newErr error) Error {
	return Error{Entity: strings.ToLower(err.Entity), Field: strings.ToLower(err.Field), Err: newErr}
}

func (err Error) WithField(newField string) Error {
	return Error{Entity: strings.ToLower(err.Entity), Field: strings.ToLower(newField), Err: err.Err}
}

func (err Error) WithEntity(newEntity string) Error {
	return Error{Entity: strings.ToLower(newEntity), Field: strings.ToLower(err.Field), Err: err.Err}
}

func (err Error) Error() string {
	return fmt.Sprintf("%s %s", err.Entity, err.Err.Error())
}

func (err Error) Is(target error) bool {
	return errors.Is(err.Err, target)
}

func (err Error) As(target any) bool {
	_, ok := target.(Error)
	return ok
}
