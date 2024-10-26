package hashing

import "github.com/protomem/chatik/pkg/werrors"

var _ Hasher = NopHasher{}

type NopHasher struct{}

func NewNopHasher() NopHasher {
	return NopHasher{}
}

func (NopHasher) Generate(plain string) (string, error) {
	return plain, nil
}

func (NopHasher) Compare(plain, hash string) error {
	equal := plain == hash
	if !equal {
		return werrors.Error(ErrInvalidHash, "hashing", "nop")
	}
	return nil
}
