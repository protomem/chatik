package api

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/protomem/chatik/internal/agregate"
	"github.com/protomem/chatik/internal/database"
	"github.com/protomem/chatik/internal/jwt"
	"github.com/protomem/chatik/internal/passhash"
	"github.com/protomem/chatik/internal/requestid"
	"github.com/protomem/chatik/internal/stream"
	"github.com/protomem/chatik/internal/validation"
	"github.com/protomem/chatik/internal/validation/vrule"
	"github.com/valyala/fasthttp"
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

		err = validation.Validate(
			vrule.Nickname(req.Nickname),
			vrule.Password(req.Password),
			vrule.Email(req.Email),
		)
		if err != nil {
			return err
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

		err = validation.Validate(
			vrule.Email(req.Email),
			vrule.Password(req.Password),
		)
		if err != nil {
			return err
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

func (srv *Server) handleListChannels() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			err error

			op     = "api.handleListChannels"
			ctx    = c.UserContext()
			logger = srv.logger.With(
				"operation", op,
				requestid.LogKey, requestid.Extract(ctx),
			)
		)
		defer func() {
			if err != nil {
				logger.Error("failed to handle list channels", "error", err)
			}
		}()

		// TODO: normalize data

		channels, err := srv.db.ChannelRepo().FindAll(ctx)
		if err != nil {
			if errors.Is(err, database.ErrChannelNotFound) {
				return fiber.NewError(fiber.StatusNotFound, database.ErrChannelNotFound.Error())
			}

			return fiber.ErrInternalServerError
		}

		userIDs := make([]uuid.UUID, 0, len(channels))
		for _, channel := range channels {
			userIDs = append(userIDs, channel.UserID)
		}

		users, err := srv.db.UserRepo().FindAllByIDs(ctx, userIDs)
		if err != nil {
			if errors.Is(err, database.ErrUserNotFound) {
				return fiber.NewError(fiber.StatusNotFound, database.ErrChannelNotFound.Error())
			}

			return fiber.ErrInternalServerError
		}

		channelAgrs := make([]agregate.Channel, 0, len(channels))
		for _, channel := range channels {
			var curUser database.User
			for _, user := range users {
				if user.ID == channel.UserID {
					curUser = user
					break
				}
			}

			if curUser.ID == uuid.Nil {
				continue
			}

			channelAgrs = append(channelAgrs, agregate.Channel{
				Channel: channel,
				User:    curUser,
			})
		}

		if len(channelAgrs) == 0 {
			return fiber.NewError(fiber.StatusNotFound, database.ErrChannelNotFound.Error())
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"channels": channelAgrs,
		})
	}
}

func (srv *Server) handleCreateChannel() fiber.Handler {
	type request struct {
		Title string `json:"title"`
	}

	return func(c *fiber.Ctx) error {
		var (
			err error

			op     = "api.handleCreateChannel"
			ctx    = c.UserContext()
			logger = srv.logger.With(
				"operation", op,
				requestid.LogKey, requestid.Extract(ctx),
			)
		)
		defer func() {
			if err != nil {
				logger.Error("failed to handle create channel", "error", err)
			}
		}()

		var req request
		err = c.BodyParser(&req)
		if err != nil {
			return fiber.ErrBadRequest
		}

		err = validation.Validate(
			vrule.Title(req.Title),
		)
		if err != nil {
			return err
		}

		authPayload, _ := jwt.Extract(ctx)

		channelID, err := srv.db.ChannelRepo().Create(ctx, database.CreateChannelDTO{
			Title:  req.Title,
			UserID: authPayload.UserID,
		})
		if err != nil {
			if errors.Is(err, database.ErrChannelExists) {
				return fiber.NewError(fiber.StatusConflict, database.ErrChannelExists.Error())
			}

			return fiber.ErrInternalServerError
		}

		channel, err := srv.db.ChannelRepo().FindByID(ctx, channelID)
		if err != nil {
			// TODO: handle error: database.ErrChannelNotFound?
			return fiber.ErrInternalServerError
		}

		user, err := srv.db.UserRepo().FindByID(ctx, channel.UserID)
		if err != nil {
			// TODO: handle error: database.ErrUserNotFound?
			return fiber.ErrInternalServerError
		}

		channelAgr := agregate.Channel{
			Channel: channel,
			User:    user,
		}

		// TODO: logging
		_ = stream.SendEvent(srv.stream, stream.NewChannelEvent(channelAgr))

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"channel": channelAgr,
		})
	}
}

func (srv *Server) handleDeleteChannel() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			err error

			op     = "api.handleDeleteChannel"
			ctx    = c.UserContext()
			logger = srv.logger.With(
				"operation", op,
				requestid.LogKey, requestid.Extract(ctx),
			)
		)
		defer func() {
			if err != nil {
				logger.Error("failed to handle delete channel", "error", err)
			}
		}()

		channelID, _ := uuid.Parse(c.Params("channelID"))
		authPayload, _ := jwt.Extract(ctx)

		err = srv.db.ChannelRepo().DeleteByIDAndUserID(ctx, channelID, authPayload.UserID)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		// TODO: logging
		_ = stream.SendEvent(srv.stream, stream.RemoveChannelEvent(channelID))

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func (srv *Server) handleListMessages() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			err error

			op     = "api.handleListMessages"
			ctx    = c.UserContext()
			logger = srv.logger.With(
				"operation", op,
				requestid.LogKey, requestid.Extract(ctx),
			)
		)
		defer func() {
			if err != nil {
				logger.Error("failed to handle list messages", "error", err)
			}
		}()

		channelID, _ := uuid.Parse(c.Params("channelID"))

		messages, err := srv.db.MessageRepo().FindAllByChannelID(ctx, channelID)
		if err != nil {
			if errors.Is(err, database.ErrMessageNotFound) {
				return fiber.NewError(fiber.StatusNotFound, database.ErrMessageNotFound.Error())
			}

			return fiber.ErrInternalServerError
		}

		userIDs := make([]uuid.UUID, 0, len(messages))
		for _, message := range messages {
			userIDs = append(userIDs, message.UserID)
		}

		users, err := srv.db.UserRepo().FindAllByIDs(ctx, userIDs)
		if err != nil {
			if errors.Is(err, database.ErrUserNotFound) {
				return fiber.NewError(fiber.StatusNotFound, database.ErrMessageNotFound.Error())
			}

			return fiber.ErrInternalServerError
		}

		messageAgrs := make([]agregate.Message, 0, len(messages))
		for _, message := range messages {
			var curUser database.User
			for _, user := range users {
				if user.ID == message.UserID {
					curUser = user
					break
				}
			}

			if curUser.ID == uuid.Nil {
				continue
			}

			messageAgrs = append(messageAgrs, agregate.Message{
				Message: message,
				User:    curUser,
			})
		}

		if len(messageAgrs) == 0 {
			return fiber.NewError(fiber.StatusNotFound, database.ErrMessageNotFound.Error())
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"messages": messageAgrs,
		})
	}
}

func (srv *Server) handleCreateMessage() fiber.Handler {
	type request struct {
		Content string `json:"content"`
	}

	return func(c *fiber.Ctx) error {
		var (
			err error

			op     = "api.handleCreateMessage"
			ctx    = c.UserContext()
			logger = srv.logger.With(
				"operation", op,
				requestid.LogKey, requestid.Extract(ctx),
			)
		)
		defer func() {
			if err != nil {
				logger.Error("failed to handle create message", "error", err)
			}
		}()

		var req request
		err = c.BodyParser(&req)
		if err != nil {
			return fiber.ErrBadRequest
		}

		err = validation.Validate(
			vrule.Content(req.Content),
		)
		if err != nil {
			return err
		}

		channelID, _ := uuid.Parse(c.Params("channelID"))
		authPayload, _ := jwt.Extract(ctx)

		_, err = srv.db.ChannelRepo().FindByID(ctx, channelID)
		if err != nil {
			if errors.Is(err, database.ErrChannelNotFound) {
				return fiber.NewError(fiber.StatusBadRequest, database.ErrChannelNotFound.Error())
			}

			return fiber.ErrInternalServerError
		}

		messageID, err := srv.db.MessageRepo().Create(ctx, database.CreateMessageDTO{
			Content:   req.Content,
			UserID:    authPayload.UserID,
			ChannelID: channelID,
		})
		if err != nil {
			return fiber.ErrInternalServerError
		}

		// TODO: handle error: database.ErrChannelNotFound?
		message, err := srv.db.MessageRepo().FindByID(ctx, messageID)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		// TODO: handle error: database.ErrUserNotFound?
		user, err := srv.db.UserRepo().FindByID(ctx, message.UserID)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		messageAgr := agregate.Message{
			Message: message,
			User:    user,
		}

		// TODO: logging
		_ = stream.SendEvent(srv.stream, stream.NewMessageEvent(messageAgr))

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": messageAgr,
		})
	}
}

func (srv *Server) handleDeleteMessage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			err error

			op     = "api.handleDeleteMessage"
			ctx    = c.UserContext()
			logger = srv.logger.With(
				"operation", op,
				requestid.LogKey, requestid.Extract(ctx),
			)
		)
		defer func() {
			if err != nil {
				logger.Error("failed to handle delete message", "error", err)
			}
		}()

		messageID, _ := uuid.Parse(c.Params("messageID"))
		authPayload, _ := jwt.Extract(ctx)

		err = srv.db.MessageRepo().DeleteByIDAndUserID(ctx, messageID, authPayload.UserID)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		// TODO: logging
		_ = stream.SendEvent(srv.stream, stream.RemoveMessageEvent(messageID))

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func (srv *Server) handleStreamWS() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		var (
			op     = "api.handleStreamWS"
			ctx    = context.Background()
			logger = srv.logger.With(
				"operation", op,
				requestid.LogKey, c.Locals("requestid"),
			)
		)

		_ = c.SetCompressionLevel(7)

		subscriber, err := srv.stream.Subscribe(ctx)
		if err != nil {
			_ = c.Close()
			return
		}
		msgs := subscriber()

		for {
			msg, ok := <-msgs
			if !ok {
				_ = c.Close()
				return
			}

			err := c.WriteMessage(websocket.TextMessage, msg.Payload())
			if err != nil {
				logger.Error("failed to write message", "error", err)
			}
		}
	})
}

func (s *Server) handleStreamSSE() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			op     = "api.handleStreamSSE"
			ctx    = context.Background()
			logger = s.logger.With(
				"operation", op,
				requestid.LogKey, c.Locals("requestid"),
			)
		)

		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")
		c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			subscriber, err := s.stream.Subscribe(ctx)
			if err != nil {
				return
			}

			msgs := subscriber()

			for {
				msg, ok := <-msgs
				if !ok {
					logger.Debug("channel closed")

					break
				}

				_, err = fmt.Fprintf(w, "data: %s\n\n", msg.Payload())
				if err != nil {
					logger.Error("failed to write sse message", "error", err)

					continue
				}

				err = w.Flush()
				if err != nil {
					logger.Error("failed to flush sse message", "error", err)

					break
				}
			}
		}))

		return nil
	}
}
