package service

import (
	"context"
	"errors"

	"github.com/protomem/chatik/internal/core/data"
	"github.com/protomem/chatik/internal/core/entity"
	"github.com/protomem/chatik/pkg/hashing"
	"github.com/protomem/chatik/pkg/validation"
	"github.com/protomem/chatik/pkg/werrors"
)

var _ User = (*UserImpl)(nil)

type (
	CreateUserDTO struct {
		Nickname string
		Password string
		Email    string
	}
)

//go:generate mockgen -destination mocks/user.go -package mocks . User
type (
	User interface {
		Find(context.Context) ([]entity.User, error)

		Get(context.Context, entity.ID) (entity.User, error)
		GetByNicknameAndPassword(context.Context, string, string) (entity.User, error)

		Create(context.Context, CreateUserDTO) (entity.User, error)
	}

	UserImpl struct {
		hasher hashing.Hasher
		acc    data.UserAccessor
		mut    data.UserMutator
	}
)

func NewUser(hasher hashing.Hasher, acc data.UserAccessor, mut data.UserMutator) *UserImpl {
	return &UserImpl{
		hasher: hasher,
		acc:    acc,
		mut:    mut,
	}
}

func (serv *UserImpl) Find(ctx context.Context) ([]entity.User, error) {
	users, err := serv.acc.Select(ctx)
	return users, werrors.Error(err, "service/user", "find")
}

func (serv *UserImpl) Get(ctx context.Context, id entity.ID) (entity.User, error) {
	user, err := serv.acc.Get(ctx, id)
	return user, werrors.Error(err, "service/user", "get")
}

func (serv *UserImpl) GetByNicknameAndPassword(ctx context.Context, nickname string, password string) (entity.User, error) {
	werr := werrors.Wrap("service/user", "getByNicknameAndPassword")

	user, err := serv.acc.GetByNickname(ctx, nickname)
	if err != nil {
		return entity.User{}, werr(err)
	}

	if err := serv.hasher.Compare(password, user.Password); err != nil {
		if errors.Is(err, hashing.ErrInvalidHash) {
			return entity.User{}, werr(entity.ErrUserNotFound)
		}

		return entity.User{}, werr(err)
	}

	return user, nil
}

func (serv *UserImpl) Create(ctx context.Context, dto CreateUserDTO) (entity.User, error) {
	werr := werrors.Wrap("service/user", "create")

	if err := dto.Validate(); err != nil {
		return entity.User{}, werr(err)
	}

	hashPassword, err := serv.hasher.Generate(dto.Password)
	if err != nil {
		return entity.User{}, werr(err)
	}

	insertDTO := data.InsertUserDTO(dto)
	insertDTO.Password = string(hashPassword)

	userID, err := serv.mut.Insert(ctx, insertDTO)
	if err != nil {
		var eerr entity.Error
		if errors.As(err, &eerr) && errors.Is(eerr.Err, entity.ErrExists) && eerr.Field != "" {
			verr := validation.New().AddFieldError(eerr.Field, "already taken")
			return entity.User{}, werr(verr)
		}

		return entity.User{}, werr(err)
	}

	user, err := serv.acc.Get(ctx, userID)
	if err != nil {
		return entity.User{}, werr(err)
	}

	return user, nil
}

func (dto CreateUserDTO) Validate() error {
	return validation.Validate(func(v *validation.Validator) {
		v.CheckField(validation.NotBlank(dto.Nickname), "nickname", "cannot be blank")
		v.CheckField(validation.MinRunes(dto.Nickname, 3), "nickname", "cannot be less than 3 characters")
		v.CheckField(validation.MaxRunes(dto.Nickname, 32), "nickname", "cannot be more than 32 characters")

		v.CheckField(validation.NotBlank(dto.Password), "password", "cannot be blank")
		v.CheckField(validation.MinRunes(dto.Password, 6), "password", "cannot be less than 6 characters")
		v.CheckField(validation.MaxRunes(dto.Password, 32), "password", "cannot be more than 32 characters")

		v.CheckField(validation.NotBlank(dto.Email), "email", "cannot be blank")
		v.CheckField(validation.IsEmail(dto.Email), "email", "must be a valid email")
	})
}
