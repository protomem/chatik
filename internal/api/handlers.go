package api

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/protomem/chatik/internal/database"
	"github.com/protomem/chatik/internal/jwt"
	"github.com/protomem/chatik/internal/passhash"
	"github.com/protomem/chatik/internal/requestid"
)

func (srv *Server) handleHealth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		res := fiber.Map{
			"status": "ok",
		}

		authPayload, ok := jwt.Extract(c.UserContext())
		if ok {
			res["userId"] = authPayload.UserID
		}

		return c.Status(fiber.StatusOK).JSON(res)
	}
}

// TODO: Add validation
func (srv *Server) handleRegister() fiber.Handler {
	type request struct {
		Nickname string `json:"nickname"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	return func(c *fiber.Ctx) error {
		var (
			err error

			op     = "api.handleRegister"
			ctx    = c.UserContext()
			logger = srv.logger.With(
				"operation", op,
				requestid.LogKey, requestid.Extract(ctx),
			)
		)
		defer func() {
			if err != nil {
				logger.Error("failed to handle register", "error", err)
			}
		}()

		var req request
		err = c.BodyParser(&req)
		if err != nil {
			return fiber.ErrBadRequest
		}

		hashPass, err := passhash.Generate(req.Password)

		userID, err := srv.db.UserRepo().Create(ctx, database.CreateUserDTO{
			Nickname: req.Nickname,
			Password: hashPass,
			Email:    req.Email,
		})
		if err != nil {
			if errors.Is(err, database.ErrUserExists) {
				return fiber.NewError(fiber.StatusConflict, database.ErrUserExists.Error())
			}

			return fiber.ErrInternalServerError
		}

		user, err := srv.db.UserRepo().FindByID(ctx, userID)
		if err != nil {
			// TODO: handle error: database.ErrUserNotFound?
			return fiber.ErrInternalServerError
		}

		payload := jwt.Payload{
			UserID:   user.ID,
			Nickname: user.Nickname,
			Email:    user.Email,
			Verified: user.Verified,
		}
		params := jwt.GenerateParams{SigningKey: srv.conf.JWT.Secret, TTL: 3 * 24 * time.Hour}
		accessToken, err := jwt.Generate(payload, params)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"accessToken": accessToken,
			"user":        user,
		})
	}
}

func (srv *Server) handleLogin() fiber.Handler {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(c *fiber.Ctx) error {
		var (
			err error

			op     = "api.handleLogin"
			ctx    = c.UserContext()
			logger = srv.logger.With(
				"operation", op,
				requestid.LogKey, requestid.Extract(ctx),
			)
		)
		defer func() {
			if err != nil {
				logger.Error("failed to handle login", "error", err)
			}
		}()

		var req request
		err = c.BodyParser(&req)
		if err != nil {
			return fiber.ErrBadRequest
		}

		user, err := srv.db.UserRepo().FindByEmail(ctx, req.Email)
		if err != nil {
			if errors.Is(err, database.ErrUserNotFound) {
				return fiber.NewError(fiber.StatusNotFound, database.ErrUserNotFound.Error())
			}

			return fiber.ErrInternalServerError
		}

		err = passhash.Compare(req.Password, user.Password)
		if err != nil {
			if errors.Is(err, passhash.ErrWrongPassword) {
				return fiber.NewError(fiber.StatusNotFound, database.ErrUserNotFound.Error())
			}

			return fiber.ErrInternalServerError
		}

		payload := jwt.Payload{
			UserID:   user.ID,
			Nickname: user.Nickname,
			Email:    user.Email,
			Verified: user.Verified,
		}
		params := jwt.GenerateParams{SigningKey: srv.conf.JWT.Secret, TTL: 3 * 24 * time.Hour}
		accessToken, err := jwt.Generate(payload, params)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"accessToken": accessToken,
			"user":        user,
		})
	}
}
