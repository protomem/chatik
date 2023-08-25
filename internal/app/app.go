package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/protomem/chatik/internal/domain/port"
	"github.com/protomem/chatik/internal/domain/usecase"
	httphandl "github.com/protomem/chatik/internal/infra/handler/http"
	httpmdw "github.com/protomem/chatik/internal/infra/middleware/http"
	mongorepo "github.com/protomem/chatik/internal/infra/repository/mongo"
	"github.com/protomem/chatik/internal/stream"
	"github.com/protomem/chatik/pkg/closer"
	"github.com/protomem/chatik/pkg/logging"
	"github.com/protomem/chatik/pkg/logging/zap"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type (
	Repositories struct {
		port.UserRepository
		port.ChannelRepository
		port.MessageRepository
	}

	UseCases struct {
		port.FindUserUseCase
		port.FindUserByEmailAndPasswordUseCase
		port.CreateUserUseCase

		port.LoginUserUseCase
		port.RegisterUserUseCase
		port.VerifyTokenUseCase

		port.FindChannelUseCase
		port.FindAllChannelsUseCase
		port.CreateChannelUseCase
		port.DeleteChannelUseCase

		port.FindMessageUseCase
		port.FindAllMessagesByChannelIDUseCase
		port.CreateMessageUseCase
		port.DeleteMessageUseCase
	}

	Handlers struct {
		*httphandl.AuthHandler
		*httphandl.ChannelHandler
		*httphandl.MessageHandler
		*httphandl.StreamHandler
	}

	Middlewares struct {
		*httpmdw.AuthMiddleware
	}
)

func newRepositories(ctx context.Context, logger logging.Logger, mdb *mongo.Client) (*Repositories, error) {
	const opMigrate = "migrate"
	var err error

	userRepo := mongorepo.NewUserRepository(logger, mdb)
	err = userRepo.Migrate(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opMigrate, err)
	}

	channelRepo := mongorepo.NewChannelRepository(logger, mdb)
	err = channelRepo.Migrate(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opMigrate, err)
	}

	messageRepo := mongorepo.NewMessageRepository(logger, mdb)
	err = messageRepo.Migrate(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opMigrate, err)
	}

	return &Repositories{
		UserRepository:    userRepo,
		ChannelRepository: channelRepo,
		MessageRepository: messageRepo,
	}, nil
}

func newUseCases(authSecret string, stream *stream.Stream, repos *Repositories) *UseCases {
	findUserUC := usecase.NewFindUserByID(repos.UserRepository)
	findUserByEmailAndPasswordUC := usecase.NewFindUserByEmailAndPassword(repos.UserRepository)
	createUserUC := usecase.NewCreateUser(repos.UserRepository)

	registerUserUC := usecase.NewRegisterUser(authSecret, createUserUC)
	loginUserUC := usecase.NewLoginUser(authSecret, findUserByEmailAndPasswordUC)
	verifyTokenUC := usecase.NewVerifyToken(authSecret, findUserUC)

	findChannelUC := usecase.NewFindChannel(repos.ChannelRepository)
	findAllChannelsUC := usecase.NewFindAllChannels(repos.ChannelRepository)
	createChannelUC := usecase.NewCreateChannel(repos.ChannelRepository, findUserUC)
	deleteChannelUC := usecase.NewDeleteChannel(repos.ChannelRepository, findChannelUC)

	findMessageUC := usecase.NewFindMessage(repos.MessageRepository)
	findAllMessagesByChannelIDUC := usecase.NewFindAllMessagesByID(repos.MessageRepository)
	createMessageUC := usecase.NewCreateMessage(repos.MessageRepository, findUserUC, findChannelUC)
	deleteMessageUC := usecase.NewDeleteMessage(repos.MessageRepository, findMessageUC)

	observeCreateChannelUC := usecase.NewObserveCreateChannel(stream, createChannelUC)
	observeDeleteChannelUC := usecase.NewObserveDeleteChannel(stream, deleteChannelUC)
	observeCreateMessageUC := usecase.NewObserveCreateMessage(stream, createMessageUC)
	observeDeleteMessageUC := usecase.NewObserveDeleteMessage(stream, deleteMessageUC)

	return &UseCases{
		FindUserUseCase:                   findUserUC,
		FindUserByEmailAndPasswordUseCase: findUserByEmailAndPasswordUC,
		CreateUserUseCase:                 createUserUC,
		RegisterUserUseCase:               registerUserUC,
		LoginUserUseCase:                  loginUserUC,
		VerifyTokenUseCase:                verifyTokenUC,
		FindChannelUseCase:                findChannelUC,
		FindAllChannelsUseCase:            findAllChannelsUC,
		CreateChannelUseCase:              observeCreateChannelUC,
		DeleteChannelUseCase:              observeDeleteChannelUC,
		FindMessageUseCase:                findMessageUC,
		FindAllMessagesByChannelIDUseCase: findAllMessagesByChannelIDUC,
		CreateMessageUseCase:              observeCreateMessageUC,
		DeleteMessageUseCase:              observeDeleteMessageUC,
	}
}

func newHandlers(logger logging.Logger, stream *stream.Stream, ucs *UseCases) *Handlers {
	return &Handlers{
		AuthHandler: httphandl.NewAuthHandler(
			logger,
			ucs.RegisterUserUseCase,
			ucs.LoginUserUseCase,
		),

		ChannelHandler: httphandl.NewChannelHandler(
			logger,
			ucs.FindAllChannelsUseCase,
			ucs.CreateChannelUseCase,
			ucs.DeleteChannelUseCase,
		),

		MessageHandler: httphandl.NewMessageHandler(
			logger,
			ucs.FindAllMessagesByChannelIDUseCase,
			ucs.CreateMessageUseCase,
			ucs.DeleteMessageUseCase,
		),

		StreamHandler: httphandl.NewStreamHandler(logger, stream),
	}
}

func newMiddlewares(logger logging.Logger, ucs *UseCases) *Middlewares {
	return &Middlewares{
		AuthMiddleware: httpmdw.NewAuthMiddleware(logger, ucs.VerifyTokenUseCase),
	}
}

type App struct {
	conf   Config
	logger logging.Logger

	mdb    *mongo.Client
	stream *stream.Stream

	repositories *Repositories
	useCases     *UseCases
	handlers     *Handlers
	middlewares  *Middlewares

	router *mux.Router
	server *http.Server

	closer *closer.Closer
}

func New(conf Config) (*App, error) {
	const op = "app.New"
	var err error
	ctx := context.Background()

	logger, err := zap.New(conf.Log.Level)
	if err != nil {
		return nil, fmt.Errorf("%s: init logger: %w", op, err)
	}

	logger.Debug("app configured", "config", conf)

	mdb, err := newMongo(ctx, conf.Mongo.URI, conf.Mongo.User, conf.Mongo.Password)
	if err != nil {
		return nil, fmt.Errorf("%s: init mongo: %w", op, err)
	}

	stream := stream.New(logger)

	repositories, err := newRepositories(ctx, logger, mdb)
	if err != nil {
		return nil, fmt.Errorf("%s: init repositories: %w", op, err)
	}

	useCases := newUseCases(conf.Auth.Secret, stream, repositories)
	handlers := newHandlers(logger, stream, useCases)
	middlewares := newMiddlewares(logger, useCases)

	router := mux.NewRouter()
	server := newServer(router, conf.HTTP.Addr)

	return &App{
		conf:         conf,
		logger:       logger,
		mdb:          mdb,
		stream:       stream,
		repositories: repositories,
		useCases:     useCases,
		handlers:     handlers,
		middlewares:  middlewares,
		router:       router,
		server:       server,
		closer:       closer.New(),
	}, nil
}

func (app *App) Run() error {
	const op = "app.Run"
	var err error
	ctx := context.Background()

	app.registerOnShutdown()
	app.setupRoutes()

	errs := make(chan error)

	go app.startServer(ctx, errs)
	go app.gracefullShutdown(ctx, errs)

	app.logger.Info("app start ...", "addr", app.conf.HTTP.Addr)
	defer app.logger.Info("app stoped.")

	err = <-errs
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (app *App) registerOnShutdown() {
	app.closer.Add(app.server.Shutdown)
	app.closer.Add(func(_ context.Context) error {
		return app.stream.Close()
	})
	app.closer.Add(app.mdb.Disconnect)
	app.closer.Add(app.logger.Sync)
}

func (app *App) setupRoutes() {
	app.router.Use(httpmdw.RequestID())
	app.router.Use(httpmdw.RequestLogger(app.logger))
	app.router.Use(httpmdw.Recovery(app.logger))

	app.router.Use(app.middlewares.AuthMiddleware.Authenticator())

	app.router.Handle("/health", httphandl.HealthCheck()).Methods(http.MethodGet)

	app.router.Handle("/api/v1/auth/register", app.handlers.AuthHandler.HandleRegisterUser()).Methods(http.MethodPost)
	app.router.Handle("/api/v1/auth/login", app.handlers.AuthHandler.HandleLoginUser()).Methods(http.MethodPost)

	app.router.Handle("/api/v1/channels", app.handlers.ChannelHandler.HandleFindAllChannels()).Methods(http.MethodGet)
	app.router.Handle("/api/v1/channels", app.handlers.ChannelHandler.HandleCreateChannel()).Methods(http.MethodPost)
	app.router.Handle("/api/v1/channels/{channelID}", app.handlers.
		ChannelHandler.HandleDeleteChannel()).Methods(http.MethodDelete)

	app.router.Handle("/api/v1/channels/{channelID}/messages", app.handlers.
		MessageHandler.HandleFindAllMessages()).Methods(http.MethodGet)
	app.router.Handle("/api/v1/channels/{channelID}/messages", app.handlers.
		MessageHandler.HandleCreateMessage()).Methods(http.MethodPost)
	app.router.Handle("/api/v1/channels/{channelID}/messages/{messageID}", app.handlers.
		MessageHandler.HandleDeleteMessage()).Methods(http.MethodDelete)

	app.router.Handle("/api/v1/stream/sse", app.handlers.StreamHandler.HandleSSE()).Methods(http.MethodGet)
	app.router.Handle("/api/v1/stream/ws", app.handlers.StreamHandler.HandleWS()).Methods(http.MethodGet)

	app.server.Handler = httpmdw.CORS()(app.router)
}

func (app *App) startServer(_ context.Context, errs chan<- error) {
	err := app.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		errs <- fmt.Errorf("start server: %w", err)
	}
}

func (app *App) gracefullShutdown(ctx context.Context, errs chan<- error) {
	<-wait()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := app.closer.Close(ctx)
	if err != nil {
		errs <- fmt.Errorf("gracefull shutdown: %w", err)
	}

	errs <- nil
}

func newMongo(ctx context.Context, uri, user, password string) (*mongo.Client, error) {
	var err error

	opts := options.Client().
		ApplyURI(uri).
		SetAuth(options.Credential{
			Username:      user,
			Password:      password,
			AuthSource:    "admin",
			AuthMechanism: "SCRAM-SHA-256",
		})

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return client, nil
}

func newServer(h http.Handler, addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: h,
	}
}

func wait() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	return ch
}
