package agregate

import "github.com/protomem/chatik/internal/database"

type Channel struct {
	database.Channel

	User database.User `json:"user"`
}

type Message struct {
	database.Message

	User database.User `json:"user"`
}
