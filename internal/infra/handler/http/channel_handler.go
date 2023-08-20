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

type ChannelHandler struct {
	logger logging.Logger

	findAllChannelsUC port.FindAllChannelsUseCase
	createChannelUC   port.CreateChannelUseCase
	deleteChannelUC   port.DeleteChannelUseCase
}

func NewChannelHandler(
	logger logging.Logger,
	findAllChannelsUC port.FindAllChannelsUseCase,
	createChannelUC port.CreateChannelUseCase,
	deleteChannelUC port.DeleteChannelUseCase,
) *ChannelHandler {
	return &ChannelHandler{
		logger:            logger.With("handlerType", "http", "handler", "channel"),
		findAllChannelsUC: findAllChannelsUC,
		createChannelUC:   createChannelUC,
		deleteChannelUC:   deleteChannelUC,
	}
}

func (h *ChannelHandler) HandleFindAllChannels() http.Handler {
	type Response struct {
		Channels []model.Channel `json:"channels"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "http.ChannelHandler.HandleFindAllChannels"
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

		channels, err := h.findAllChannelsUC.Invoke(ctx)
		if err != nil {
			logger.Error("failed to find all channels", "error", err)

			code := http.StatusInternalServerError
			res := map[string]any{
				"error": "failed to find all channels",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(Response{
			Channels: channels,
		})
	})
}

func (h *ChannelHandler) HandleCreateChannel() http.Handler {
	type Request struct {
		Title string `json:"title"`
	}

	type Response struct {
		Channel model.Channel `json:"channel"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "http.ChannelHandler.HandleCreateChannel"
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
			err = json.NewEncoder(w).Encode(map[string]any{
				"error": "invalid request body",
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

		channel, err := h.createChannelUC.Invoke(ctx, port.CreateChannelUCDTO{
			Title:  req.Title,
			UserID: authPayload.UserID,
		})
		if err != nil {
			logger.Error("failed to create channel", "error", err)

			code := http.StatusInternalServerError
			res := map[string]any{
				"error": "failed to create channel",
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

			if errors.Is(err, model.ErrChannelAlreadyExists) {
				code = http.StatusConflict
				res = map[string]any{
					"error": model.ErrChannelAlreadyExists.Error(),
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
			Channel: channel,
		})
	})
}

func (h *ChannelHandler) HandleDeleteChannel() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "http.ChannelHandler.HandleDeleteChannel"
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

		err = h.deleteChannelUC.Invoke(ctx, port.DeleteChannelUCDTO{
			ChannelID: channelID,
			UserID:    authPayload.UserID,
		})
		if err != nil {
			logger.Error("failed to delete channel", "error", err)

			code := http.StatusInternalServerError
			res := map[string]any{
				"error": "failed to delete channel",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
