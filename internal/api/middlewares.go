package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/protomem/chatik/internal/jwt"
	"github.com/protomem/chatik/internal/requestid"
)

func (srv *Server) requestLogging() fiber.Handler {
	conf := logger.ConfigDefault
	conf.Output = srv.logger
	conf.Format = "request [${time}] ${locals:requestid} ${status} - ${latency} ${method} ${path}"

	return logger.New(conf)
}

func (srv *Server) recovery() fiber.Handler {
	return recover.New()
}

func (srv *Server) CORS() fiber.Handler {
	return cors.New(cors.ConfigDefault)
}

func (srv *Server) authenticator() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			op     = "api.authenticator"
			ctx    = c.UserContext()
			logger = srv.logger.With(
				"operation", op,
				requestid.LogKey, requestid.Extract(ctx),
			)
		)

		var token string
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			token = c.Query("token")
			if token == "" {
				logger.Debug("auth token missing")

				return c.Next()
			}
		} else {
			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 {
				logger.Error("invalid auth header")

				return fiber.ErrUnauthorized
			}

			token = headerParts[1]
		}

		params := jwt.ParseParams{SigningKey: srv.conf.JWT.Secret}
		payload, err := jwt.Parse(token, params)
		if err != nil {
			logger.Error("invalid token", "error", err)

			return fiber.ErrUnauthorized
		}

		_, err = srv.db.UserRepo().FindByID(ctx, payload.UserID)
		if err != nil {
			logger.Error("failed to verify token", "error", err)

			return fiber.ErrUnauthorized
		}

		ctx = jwt.Inject(ctx, payload)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

func (srv *Server) authorizer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, ok := jwt.Extract(c.UserContext())
		if !ok {
			return fiber.NewError(fiber.StatusForbidden, "access denied")
		}

		return c.Next()
	}
}
