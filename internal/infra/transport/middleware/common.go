package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/protomem/chatik/internal/infra/ctxstore"
	"github.com/protomem/chatik/internal/infra/logging"
	"github.com/tomasen/realip"
)

func RealIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := realip.FromRequest(r)

		ctx := ctxstore.With(r.Context(), ctxstore.ClientIP, ip)
		wr := r.WithContext(ctx)

		w.Header().Set("X-Real-Ip", ip)
		next.ServeHTTP(w, wr)
	})
}

func TraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tid := r.Header.Get("X-Request-Id")
		if tid == "" {
			tid = genTraceID()
		}

		ctx := ctxstore.With(r.Context(), ctxstore.TraceID, tid)
		wr := r.WithContext(ctx)

		w.Header().Set("X-Request-Id", tid)
		next.ServeHTTP(w, wr)
	})
}

func genTraceID() string {
	id, _ := uuid.NewRandom()
	return id.String()
}

func LogAccess(log *logging.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := newResponseWrapper(w)
			next.ServeHTTP(ww, r)

			var (
				ip     = ctxstore.MustFrom[string](r.Context(), ctxstore.ClientIP)
				method = r.Method
				url    = r.URL.String()
				proto  = r.Proto
			)

			userAttrs := logging.Group("user", ctxstore.ClientIP.String(), ip)
			requestAttrs := logging.Group("request", "method", method, "url", url, "proto", proto)
			responseAttrs := logging.Group("response", "status", ww.StatusCode, "size", ww.BytesCount)

			log.WithContext(r.Context()).Info("access", userAttrs, requestAttrs, responseAttrs)
		})
	}
}

type responseWrapper struct {
	StatusCode    int
	BytesCount    int
	headerWritten bool
	wrapped       http.ResponseWriter
}

func newResponseWrapper(w http.ResponseWriter) *responseWrapper {
	return &responseWrapper{
		StatusCode: http.StatusOK,
		wrapped:    w,
	}
}

func (mw *responseWrapper) Header() http.Header {
	return mw.wrapped.Header()
}

func (mw *responseWrapper) WriteHeader(statusCode int) {
	mw.wrapped.WriteHeader(statusCode)

	if !mw.headerWritten {
		mw.StatusCode = statusCode
		mw.headerWritten = true
	}
}

func (mw *responseWrapper) Write(b []byte) (int, error) {
	mw.headerWritten = true

	n, err := mw.wrapped.Write(b)
	mw.BytesCount += n
	return n, err
}

func (mw *responseWrapper) Unwrap() http.ResponseWriter {
	return mw.wrapped
}
