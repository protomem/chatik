package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/protomem/chatik/internal/domain/model"
	"github.com/protomem/chatik/internal/domain/port"
	"github.com/protomem/chatik/internal/requestid"
	"github.com/protomem/chatik/pkg/logging"
	"github.com/protomem/chatik/pkg/validation"
)

type AuthHandler struct {
	logger logging.Logger

	registerUserUC port.RegisterUserUseCase
	loginUserUC    port.LoginUserUseCase
}

func NewAuthHandler(
	logger logging.Logger,
	registerUserUC port.RegisterUserUseCase,
	loginUserUC port.LoginUserUseCase,
) *AuthHandler {
	return &AuthHandler{
		logger:         logger.With("handlerType", "http", "handler", "auth"),
		registerUserUC: registerUserUC,
		loginUserUC:    loginUserUC,
	}
}

func (h *AuthHandler) HandleRegisterUser() http.Handler {
	type Request struct {
		Nickname string `json:"nickname"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type Response struct {
		User  model.User `json:"user"`
		Token port.Token `json:"accessToken"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "http.AuthHandler.HandleRegisterUser"
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
			logger.Error("failed to decode request body", "error", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]any{
				"error": "invalid request body",
			})

			return
		}

		user, token, err := h.registerUserUC.Invoke(ctx, port.RegisterUserUCDTO{
			Nickname: req.Nickname,
			Password: req.Password,
			Email:    req.Email,
		})
		if err != nil {
			logger.Error("failed to register user", "error", err)

			code := http.StatusInternalServerError
			res := map[string]any{
				"error": "failed to register user",
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

			if errors.Is(err, model.ErrUserAlreadyExists) {
				code = http.StatusConflict
				res = map[string]any{
					"error": model.ErrUserAlreadyExists.Error(),
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
			User:  user,
			Token: token,
		})
	})
}

func (h *AuthHandler) HandleLoginUser() http.Handler {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type Response struct {
		User  model.User `json:"user"`
		Token port.Token `json:"accessToken"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "http.AuthHandler.HandleLoginUser"
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
			logger.Error("failed to decode request body", "error", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]any{
				"error": "invalid request body",
			})

			return
		}

		user, token, err := h.loginUserUC.Invoke(ctx, port.LoginUserUCDTO{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			logger.Error("failed to login user", "error", err)

			code := http.StatusInternalServerError
			res := map[string]any{
				"error": "failed to login user",
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

			if errors.Is(err, model.ErrUserNotFound) {
				code = http.StatusNotFound
				res = map[string]any{
					"error": model.ErrUserNotFound.Error(),
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(Response{
			User:  user,
			Token: token,
		})
	})
}
