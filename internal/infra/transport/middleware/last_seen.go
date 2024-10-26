package middleware

import (
	"net/http"
	"time"

	"github.com/protomem/chatik/internal/core/data"
	"github.com/protomem/chatik/internal/core/entity"
	"github.com/protomem/chatik/internal/infra/ctxstore"
	"github.com/protomem/chatik/internal/infra/logging"
)

func LastSeen(log *logging.Logger, lastseen data.LastSeen) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			user, ok := ctxstore.From[entity.User](ctx, ctxstore.User)
			if !ok {
				h.ServeHTTP(w, r)
				return
			}

			clientIp := ctxstore.MustFrom[string](ctx, ctxstore.ClientIP)
			writeLog := entity.LastSeenLog{Timestamp: time.Now(), User: user.ID, Fingerprint: clientIp}

			if err := lastseen.Write(ctx, writeLog); err != nil {
				log.Error("failed to write timestamp last seen", logging.Error(err))
			}

			h.ServeHTTP(w, r)
		})
	}
}
