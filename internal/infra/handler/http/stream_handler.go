package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/protomem/chatik/internal/requestid"
	"github.com/protomem/chatik/internal/stream"
	"github.com/protomem/chatik/pkg/logging"
)

type StreamHandler struct {
	logger logging.Logger

	stream *stream.Stream
}

func NewStreamHandler(logger logging.Logger, stream *stream.Stream) *StreamHandler {
	return &StreamHandler{
		logger: logger.With("handlerType", "http", "handler", "stream"),
		stream: stream,
	}
}

func (h *StreamHandler) HandleSSE() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "http.StreamHandler.HandleSSE"
		var err error

		ctx := r.Context()
		logger := h.logger.With(
			requestid.LogKey, requestid.Extract(ctx),
			"operation", op,
		)

		defer func() {
			if err != nil {
				logger.Error("failed to send response", "error", err)
			}
		}()

		subscriber, err := h.stream.Subscribe(ctx)
		if err != nil {
			logger.Error("failed to subscribe to stream", "error", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(map[string]any{
				"error": "failed to subscribe to stream",
			})

			return
		}
		msgs := subscriber()

		flusher, ok := w.(http.Flusher)
		if !ok {
			logger.Error("failed to set flusher")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(map[string]any{
				"error": "failed to set flusher",
			})

			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Transfer-Encoding", "chunked")

		for {
			select {
			case msg, ok := <-msgs:
				if !ok {
					logger.Debug("channel closed")
					return
				}

				_, err = fmt.Fprintf(w, "data: %s\n\n", msg.Payload())
				if err != nil {
					logger.Error("failed to write message", "error", err)
					continue
				}

				flusher.Flush()
			case <-ctx.Done():
				return
			}
		}
	})
}

func (h *StreamHandler) HandleWS() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "http.StreamHandler.HandleWS"
		var err error

		ctx := r.Context()
		logger := h.logger.With(
			requestid.LogKey, requestid.Extract(ctx),
			"operation", op,
		)

		defer func() {
			if err != nil {
				logger.Error("failed to send response", "error", err)
			}
		}()

		subscriber, err := h.stream.Subscribe(ctx)
		if err != nil {
			logger.Error("failed to subscribe to stream", "error", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(map[string]any{
				"error": "failed to subscribe to stream",
			})

			return
		}

		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error("failed to upgrade connection", "error", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(map[string]any{
				"error": "failed to upgrade connection",
			})

			return
		}

		for {
			select {
			case msg, ok := <-subscriber():
				if !ok {
					logger.Debug("channel closed")

					err = conn.Close()
					if err != nil {
						logger.Error("failed to close connection", "error", err)
					}

					return
				}

				err = conn.WriteMessage(websocket.TextMessage, msg.Payload())
				if err != nil {
					logger.Error("failed to write message", "error", err)
					continue
				}
			case <-ctx.Done():
				err = conn.Close()
				if err != nil {
					logger.Error("failed to close connection", "error", err)
				}

				return
			}
		}
	})
}
