package hashing

import "errors"

var ErrInvalidHash = errors.New("invalid hash")

type Hasher interface {
	Generate(plain string) (string, error)
	Compare(plain, hash string) error
}
