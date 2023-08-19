package http

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/protomem/chatik/internal/requestid"
	"github.com/protomem/chatik/pkg/logging"
)

func RequestID() mux.MiddlewareFunc {
	return requestid.HttpMiddleware()
}

func RequestLogger(logger logging.Logger) mux.MiddlewareFunc {
	logger = logger.With("middlewareType", "http", "middleware", "requestLogger")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlers.CombinedLoggingHandler(
				logger.With(requestid.LogKey, requestid.Extract(r.Context())),
				next,
			).ServeHTTP(w, r)
		})
	}
}

func Recovery(logger logging.Logger) mux.MiddlewareFunc {
	logger = logger.With("middlewareType", "http", "middleware", "recovery")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlers.RecoveryHandler(
				handlers.RecoveryLogger(
					logger.With(requestid.LogKey, requestid.Extract(r.Context())),
				),
				handlers.PrintRecoveryStack(true),
			)(next).ServeHTTP(w, r)
		})
	}
}

func CORS() mux.MiddlewareFunc {
	return handlers.CORS()
}
