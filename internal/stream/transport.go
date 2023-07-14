package stream

import (
	"context"

	"github.com/fasthttp/websocket"
	"github.com/protomem/chatik/internal/logging"
	"github.com/protomem/chatik/internal/requestid"
	"github.com/valyala/fasthttp"
)

var _ Transport = (*WebSocket)(nil)

type Transport interface {
	Handle(ctx context.Context, c *fasthttp.RequestCtx, session *Session) error
}

type WebSocket struct {
	logger   logging.Logger
	upgrader *websocket.FastHTTPUpgrader
}

func NewWebSocket(logger logging.Logger) *WebSocket {
	return &WebSocket{
		logger: logger,
		upgrader: &websocket.FastHTTPUpgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (ws *WebSocket) Handle(ctx context.Context, c *fasthttp.RequestCtx, session *Session) error {
	var (
		err error

		op     = "websocket.Handle"
		logger = ws.logger.With(
			"operation", op,
			"requestid", requestid.Extract(ctx),
		)
	)
	defer func() {
		if err != nil {
			logger.Error("failed to handle websocket", "error", err)
		}
	}()

	err = ws.upgrader.Upgrade(c, func(c *websocket.Conn) {
		for {
			data, ok := <-session.Receive()
			if !ok {
				return
			}

			err := c.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				logger.Error("stream.WriteMessage", "error", err)
			}

		}
	})
	if err != nil {
		return err
	}

	return nil
}
