package api

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/protomem/chatik/internal/closer"
	"github.com/protomem/chatik/internal/config"
	"github.com/protomem/chatik/internal/database"
	"github.com/protomem/chatik/internal/logging"
	"github.com/protomem/chatik/internal/requestid"
)

type Server struct {
	conf   config.Config
	logger logging.Logger

	db *database.DB

	app *fiber.App

	closer *closer.Closer
}

func NewServer(conf config.Config) *Server {
	return &Server{conf: conf}
}

func (srv *Server) Run() error {
	var (
		err error

		op  = "api.Run"
		ctx = context.Background()
	)

	err = srv.configure(ctx)
	if err != nil {
		return fmt.Errorf("%s: configure: %w", op, err)
	}

	srv.logger.Debug("configure done", "conf", srv.conf)

	srv.setuoRoutes()
	srv.registerOnShutdown()

	errs := make(chan error)

	go func() { srv.startServer(ctx, errs) }()
	go func() { srv.gracefullShutdown(ctx, errs) }()

	err = <-errs
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (srv *Server) configure(ctx context.Context) error {
	var err error

	// init logger
	srv.logger, err = logging.NewZap(srv.conf.Log.Level, srv.conf.Log.File)
	if err != nil {
		return fmt.Errorf("new logger: %w", err)
	}

	// init database
	srv.db, err = database.New(ctx, srv.logger, srv.conf.DB.URI)
	if err != nil {
		return fmt.Errorf("new database: %w", err)
	}

	err = srv.db.UserRepo().CreateIndexes(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// init app
	srv.app = fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var apiErr *fiber.Error
			if errors.As(err, &apiErr) {
				code = apiErr.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// init closer
	srv.closer = closer.New()

	return nil
}

func (srv *Server) setuoRoutes() {
	srv.app.Use(requestid.Middleware())
	srv.app.Use(srv.requestLogging())
	srv.app.Use(srv.recovery())

	srv.app.Use(srv.authenticator())

	srv.app.Get("/health", srv.handleHealth())

	v1 := srv.app.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.Post("/register", srv.handleRegister())
			auth.Post("/login", srv.handleLogin())
		}
	}
}

func (srv *Server) registerOnShutdown() {
	srv.closer.Add(srv.app.ShutdownWithContext)
	srv.closer.Add(srv.db.Close)
	srv.closer.Add(func(ctx context.Context) error {
		return srv.logger.Sync()
	})
}

func (srv *Server) startServer(_ context.Context, errs chan<- error) {
	srv.logger.Info("start server", "addr", srv.conf.HTTP.Addr)

	err := srv.app.Listen(srv.conf.HTTP.Addr)
	if err != nil {
		errs <- fmt.Errorf("start server: %w", err)
	}
}

func (srv *Server) gracefullShutdown(ctx context.Context, errs chan<- error) {
	<-wait()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	srv.logger.Info("gracefull shutdown")

	err := srv.closer.Close(ctx)
	if err != nil {
		errs <- fmt.Errorf("gracefull shutdown: %w", err)
	}

	errs <- nil
}

func wait() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	return ch
}
