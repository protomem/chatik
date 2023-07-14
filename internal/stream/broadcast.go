package stream

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/protomem/chatik/internal/logging"
	"github.com/protomem/chatik/internal/requestid"
)

var ErrSessionClosed = errors.New("session closed")

type Broadcast struct {
	logger    logging.Logger
	sessions  map[uuid.UUID]*Session
	transport Transport
}

func NewBroadcast(logger logging.Logger, transport Transport) *Broadcast {
	return &Broadcast{
		logger:    logger.With("component", "stream"),
		sessions:  make(map[uuid.UUID]*Session),
		transport: transport,
	}
}

func (b *Broadcast) Handle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			err error

			op     = "stream.Handle"
			ctx    = c.UserContext()
			logger = b.logger.With(
				"operation", op,
				requestid.LogKey, requestid.Extract(ctx),
			)
		)
		defer func() {
			if err != nil {
				logger.Error("failed to handle stream", "error", err)
			}
		}()

		c.Request().Header.Del(fiber.HeaderOrigin)

		session := NewSession()
		b.sessions[session.id] = session

		err = b.transport.Handle(ctx, c.Context(), session)
		if err != nil {
			return err
		}

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
