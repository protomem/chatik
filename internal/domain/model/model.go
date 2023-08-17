package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrChannelNotFound      = errors.New("channel not found")
	ErrChannelAlreadyExists = errors.New("channel already exists")
)

type User struct {
	ID uuid.UUID `json:"id"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Nickname string `json:"nickname"`
	Password string `json:"-"`

	Email    string `json:"email"`
	Verified bool   `json:"isVerified"`
}

type Channel struct {
	ID uuid.UUID `json:"id"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Title string `json:"title"`

	User User `json:"user"`
}
