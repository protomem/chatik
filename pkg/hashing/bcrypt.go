package hashing

import (
	"errors"

	"github.com/protomem/chatik/pkg/werrors"
	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultBcryptCost = bcrypt.DefaultCost
	MaxBcryptCost     = bcrypt.MaxCost
	MinBcryptCost     = bcrypt.MinCost
)

var _ Hasher = BcrypHasher{}

type BcrypHasher struct {
	Cost int
}

func NewBcryptHasher(cost int) BcrypHasher {
	return BcrypHasher{
		Cost: cost,
	}
}

func (h BcrypHasher) Generate(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), h.Cost)
	return string(hash), werrors.Error(err, "hashing", "bcrypt")
}

func (BcrypHasher) Compare(plain, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if err != nil {
		werr := werrors.Wrap("hashing", "bcrypt")
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return werr(ErrInvalidHash)
		}

		return werr(err)
	}

	return nil
}
