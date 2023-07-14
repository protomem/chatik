package stream

import (
	"errors"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/protomem/chatik/internal/logging"
	"github.com/protomem/chatik/internal/requestid"
)

var ErrSessionClosed = errors.New("session closed")

var _upgrader = &websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Broadcast struct {
	logger   logging.Logger
	sessions map[uuid.UUID]*Session
}

func NewBroadcast(logger logging.Logger) *Broadcast {
	return &Broadcast{
		logger:   logger.With("component", "stream"),
		sessions: make(map[uuid.UUID]*Session),
	}
}

func (b *Broadcast) Handle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := b.logger.With(
			"operation", "stream.Handle",
			requestid.LogKey, requestid.Extract(c.UserContext()),
		)

		c.Request().Header.Del(fiber.HeaderOrigin)

		session := NewSession()
		b.sessions[session.id] = session

		_ = _upgrader.Upgrade(c.Context(), func(c *websocket.Conn) {
			for {
				data, ok := <-session.Receive()
				if !ok {
					delete(b.sessions, session.id)
					return
				}

				err := c.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					logger.Error("stream.WriteMessage", "error", err)
				}

			}
		})

		return nil
	}
}

func (b *Broadcast) SendMessage(msg []byte) {
	for _, session := range b.sessions {
		session.Send(msg)
	}
}

func (b *Broadcast) Close() {
	for _, session := range b.sessions {
		session.Close()
	}
}

type Session struct {
	id   uuid.UUID
	buff chan []byte
}

func NewSession() *Session {
	return &Session{
		id:   uuid.New(),
		buff: make(chan []byte),
	}
}

func (s *Session) Send(data []byte) {
	s.buff <- data
}

func (s *Session) Receive() <-chan []byte {
	return s.buff
}

func (s *Session) Close() {
	close(s.buff)
}
