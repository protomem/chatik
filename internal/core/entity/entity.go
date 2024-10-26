package entity

import "time"

type Entity struct {
	ID        ID        `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type User struct {
	Entity

	Nickname string `json:"nickname"`
	Password string `json:"-"`

	Email    string `json:"email"`
	Verified bool   `json:"verified"`
}

type Session struct {
	CreatedAt time.Time `json:"createdAt"`
	ExpiredAt time.Time `json:"expiredAt"`

	User  ID     `json:"userId"`
	Token string `json:"token"`

	Fingerprint string `json:"fingerprint"`
}

type LastSeenLog struct {
	Timestamp   time.Time `json:"timestamp"`
	User        ID        `json:"userId"`
	Fingerprint string    `json:"fingerprint"`
}

type Message struct {
	Entity

	From User `json:"fromUser"`
	To   User `json:"toUser"`

	Text     string            `json:"text"`
	Metadata map[string]string `json:"metadata"`
}
