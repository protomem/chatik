package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/protomem/chatik/internal/domain/port"
	"github.com/protomem/chatik/internal/jwt"
	"github.com/protomem/chatik/internal/requestid"
	"github.com/protomem/chatik/pkg/logging"
)

type AuthMiddleware struct {
	logger logging.Logger

	verifyTokenUC port.VerifyTokenUseCase
}

func NewAuthMiddleware(logger logging.Logger, verifyTokenUC port.VerifyTokenUseCase) *AuthMiddleware {
	return &AuthMiddleware{
		logger:        logger.With("middlewareType", "http", "middleware", "auth"),
		verifyTokenUC: verifyTokenUC,
	}
}

func (m *AuthMiddleware) Authenticator() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "http.AuthMiddleware.Authenticator"
			var err error

			ctx := r.Context()
			logger := m.logger.With(
				requestid.LogKey, requestid.Extract(ctx),
				"operation", op,
			)

			defer func() {
				if err != nil {
					logger.Error("failed to send response", "error", err)
				}
			}()

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Debug("no authorization header")

				next.ServeHTTP(w, r)
				return
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				logger.Error("invalid authorization header")

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				err = json.NewEncoder(w).Encode(map[string]string{
					"error": "invalid authorization header",
				})

				return
			}

			token := headerParts[1]
			_, payload, err := m.verifyTokenUC.Invoke(ctx, token)
			if err != nil {
				logger.Error("failed to verify token", "error", err)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				err = json.NewEncoder(w).Encode(map[string]string{
					"error": "invalid token",
				})

				return
			}

			ctxWithPayload := jwt.Inject(ctx, payload)
			next.ServeHTTP(w, r.WithContext(ctxWithPayload))
		})
	}
}
