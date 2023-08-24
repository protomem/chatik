package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/protomem/chatik/internal/domain/model"
	"github.com/protomem/chatik/internal/domain/port"
	"github.com/protomem/chatik/internal/jwt"
	"github.com/protomem/chatik/internal/requestid"
	"github.com/protomem/chatik/pkg/logging"
	"github.com/protomem/chatik/pkg/validation"
)

type MessageHandler struct {
	logger logging.Logger

	findAllMessagesByChannelIDUC port.FindAllMessagesByChannelIDUseCase
	createMessageUC              port.CreateMessageUseCase
}

func NewMessageHandler(
	logger logging.Logger,
	findAllMessagesByChannelIDUC port.FindAllMessagesByChannelIDUseCase,
	createMessageUC port.CreateMessageUseCase,
) *MessageHandler {
	return &MessageHandler{
		logger:                       logger.With("handlerType", "http", "handler", "message"),
		findAllMessagesByChannelIDUC: findAllMessagesByChannelIDUC,
		createMessageUC:              createMessageUC,
	}
}

func (h *MessageHandler) HandleFindAllMessages() http.Handler {
	type Response struct {
		Messages []model.Message `json:"messages"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "http.MessageHandler.HandleFindAllMessages"
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

		channelIDRaw, ok := mux.Vars(r)["channelID"]
		if !ok {
			logger.Error("failed to extract channel id from request")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]any{
				"error": "missing channel id",
			})

			return
		}

		channelID, err := uuid.Parse(channelIDRaw)
		if err != nil {
			logger.Error("failed to parse channel id", "error", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]any{
				"error": "invalid channel id",
			})

			return
		}

		messages, err := h.findAllMessagesByChannelIDUC.Invoke(ctx, channelID)
		if err != nil {
			logger.Error("failed to find messages", "error", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(map[string]any{
				"error": "failed to find messages",
			})

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(Response{
			Messages: messages,
		})
	})
}

func (h *MessageHandler) HandleCreateMessage() http.Handler {
	type Request struct {
		Content string `json:"content"`
	}

	type Response struct {
		Message model.Message `json:"message"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "http.MessageHandler.HandleCreateMessage"
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

		var req Request
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.Error("failed to decode request", "error", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid request",
			})

			return
		}

		channelIDRaw, ok := mux.Vars(r)["channelID"]
		if !ok {
			logger.Error("failed to extract channel id from request")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]any{
				"error": "missing channel id",
			})

			return
		}

		channelID, err := uuid.Parse(channelIDRaw)
		if err != nil {
			logger.Error("failed to parse channel id", "error", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]any{
				"error": "invalid channel id",
			})

			return
		}

		authPayload, ok := jwt.Extract(ctx)
		if !ok {
			logger.Error("failed to extract auth payload")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			err = json.NewEncoder(w).Encode(map[string]any{
				"error": "access denied",
			})

			return
		}

		message, err := h.createMessageUC.Invoke(ctx, port.CreateMessageUCDTO{
			Content:   req.Content,
			ChannelID: channelID,
			UserID:    authPayload.UserID,
		})
		if err != nil {
			logger.Error("failed to create message", "error", err)

			code := http.StatusInternalServerError
			res := map[string]any{
				"error": "failed to create message",
			}

			var v *validation.Validator
			if errors.As(err, &v) {
				var vErrs []string

				vErrs = append(vErrs, v.Errors...)

				for vFieldErrK, vFieldErrV := range v.FieldErrors {
					vErrs = append(vErrs, fmt.Sprintf("%s: %s", vFieldErrK, vFieldErrV))
				}

				code = http.StatusBadRequest
				res = map[string]any{
					"error":   v.Error(),
					"details": vErrs,
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(Response{
			Message: message,
		})
	})
}
