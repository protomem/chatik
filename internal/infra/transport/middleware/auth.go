package middleware

import (
	"net/http"
	"strings"

	"github.com/protomem/chatik/internal/core/entity"
	"github.com/protomem/chatik/internal/core/service"
	"github.com/protomem/chatik/internal/infra/ctxstore"
	"github.com/protomem/chatik/internal/infra/logging"
	"github.com/protomem/chatik/internal/infra/transport"
)

type Auth struct {
	log  *logging.Logger
	serv service.Auth
}

func NewAuth(log *logging.Logger, serv service.Auth) *Auth {
	return &Auth{
		log:  log.With("middleware", "auth"),
		serv: serv,
	}
}

func (m *Auth) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := m.log.WithContext(ctx)

		authToken, authTokenOk := m.extractAuthToken(r)
		if !authTokenOk {
			next.ServeHTTP(w, r)
			return
		}

		if authToken == "" {
			log.Warn("invalid auth header")
			transport.WriteWithoutBody(w, http.StatusUnauthorized)
			return
		}

		user, err := m.serv.Verify(ctx, authToken)
		if err != nil {
			log.Warn("invalid auth token", logging.Error(err))
			transport.WriteWithoutBody(w, http.StatusUnauthorized)
			return
		}

		ctx = ctxstore.With(ctx, ctxstore.User, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Auth) Protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := m.log.WithContext(ctx)

		_, isAuth := ctxstore.From[entity.User](ctx, ctxstore.User)
		if !isAuth {
			log.Warn("user is not authenticated")
			transport.WriteWithoutBody(w, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (*Auth) extractAuthToken(r *http.Request) (string, bool) {
	header := r.Header.Get(transport.HeaderAuthorization)
	if header == "" {
		return "", false
	}

	parts := strings.Split(header, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1], true
	}

	return "", true
}
