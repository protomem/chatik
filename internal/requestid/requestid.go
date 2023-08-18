package requestid

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	Header = fiber.HeaderXRequestID
	LogKey = "requestId"
)

func Empty() string {
	return uuid.Nil.String()
}

func Generator() string {
	rid, err := uuid.NewRandom()
	if err != nil {
		return Empty()
	}

	return rid.String()
}

func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		rid, ok := c.GetReqHeaders()[Header]
		if !ok {
			rid = Generator()
		}

		ctx := Inject(c.UserContext(), rid)
		c.SetUserContext(ctx)

		c.Locals("requestid", rid)

		return c.Next()
	}
}

func HttpMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rid := r.Header.Get(Header)
			if rid == "" {
				rid = Generator()
			}

			ctx := Inject(r.Context(), rid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type ctxKey struct{}

func Inject(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, ctxKey{}, rid)
}

func Extract(ctx context.Context) string {
	rid, ok := ctx.Value(ctxKey{}).(string)
	if !ok {
		return Empty()
	}

	return rid
}
