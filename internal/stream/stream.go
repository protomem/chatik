package stream

import (
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/protomem/chatik/internal/logging"
	"github.com/protomem/chatik/internal/requestid"
)

var _upgrader = &websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Stream struct {
	logger   logging.Logger
	sessions map[uuid.UUID]chan []byte
}

func New(logger logging.Logger) *Stream {
	return &Stream{
		logger:   logger.With("component", "stream"),
		sessions: make(map[uuid.UUID]chan []byte),
	}
}

func (s *Stream) Handle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := s.logger.With(
			"operation", "stream.Handle",
			requestid.LogKey, requestid.Extract(c.UserContext()),
		)

		session := make(chan []byte)
		sessionID := uuid.New()
		s.sessions[sessionID] = session

		_ = _upgrader.Upgrade(c.Context(), func(c *websocket.Conn) {
			for {
				logger.Debug("listening for new message")

				data, ok := <-session
				if !ok {
					delete(s.sessions, sessionID)
					return
				}

				logger.Debug("new message", "data", string(data))

				err := c.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					logger.Error("stream.WriteMessage", "error", err)
				}

				logger.Debug("message sent")
			}
		})

		return nil
	}
}

func (s *Stream) SendMessage(msg []byte) {
	for _, session := range s.sessions {
		session <- msg
	}
}

func (s *Stream) Close() error {
	for _, session := range s.sessions {
		close(session)
	}

	return nil
}
