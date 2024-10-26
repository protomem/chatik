package app

import (
	"context"
	"net/http"
	"time"

	"github.com/protomem/chatik/cmd/api-server/config"
	"github.com/protomem/chatik/internal/core/service"
	"github.com/protomem/chatik/internal/infra/ctxstore"
	"github.com/protomem/chatik/internal/infra/database"
	dbAdapter "github.com/protomem/chatik/internal/infra/database/adapter"
	flashAdapter "github.com/protomem/chatik/internal/infra/flashstore/adapter"
	"github.com/protomem/chatik/internal/infra/logging"
	"github.com/protomem/chatik/internal/infra/transport"
	"github.com/protomem/chatik/internal/infra/transport/handler"
	"github.com/protomem/chatik/internal/infra/transport/middleware"
	"github.com/protomem/chatik/pkg/closing"
	"github.com/protomem/chatik/pkg/hashing"
	"github.com/protomem/chatik/pkg/syslog"
)

const _shutdownTimeout = 15 * time.Second

func Run(fs Features) {
	ctx := context.Background()
	syslog.Log(syslog.Info).Printf("Starting API server...")

	conf, err := config.New(fs.ConfigFile)
	if err != nil {
		syslog.Log(syslog.Error).Panicf("Failed to load config: %s", err)
	}

	log, err := logging.New(conf.Logger, LoggerExtractor)
	if err != nil {
		syslog.Log(syslog.Error).Panicf("Failed to initialize logger: %s", err)
	}

	db, err := database.Connect(ctx, log, conf.DB)
	if err != nil {
		syslog.Log(syslog.Error).Panicf("Failed to connect to database: %s", err)
	}

	sessmng, err := flashAdapter.NewSessionManager(ctx, conf.SessionManager)
	if err != nil {
		syslog.Log(syslog.Error).Panicf("Failed to initialize session manager: %s", err)
	}

	lseen, err := flashAdapter.NewLastSeen(ctx, conf.LastSeen)
	if err != nil {
		syslog.Log(syslog.Error).Panicf("Failed to initialize last seen: %s", err)
	}

	var (
		adapters    = NewAdapters(db, sessmng, lseen)
		services    = NewServices(conf, adapters)
		handlers    = NewHandlers(log, services)
		middlewares = NewMiddlewares(log, services)
		router      = SetupRoutes(log, adapters, handlers, middlewares)
	)

	srv := transport.NewServer(conf.HttpServer)
	srv.SetHandler(router)

	closer := closing.New()
	closer.Add(srv.Shutdown)
	closer.Add(db.Close)
	closer.Add(sessmng.Close)
	closer.Add(lseen.Close)

	closeErr := make(chan error, 1)
	go func() {
		<-closing.WaitQuit()

		syslog.Log(syslog.Info).Printf("Shutting down API server...")

		ctx, cancel := context.WithTimeout(context.Background(), _shutdownTimeout)
		defer cancel()

		closeErr <- closer.Close(ctx)
	}()

	syslog.Log(syslog.Info).Printf("API server started on %s", srv.Addr())
	defer syslog.Log(syslog.Info).Printf("API server stopped")

	if err := srv.Run(); err != nil {
		syslog.Log(syslog.Error).Panicf("Failed to start API server: %s", err)
	}

	if err := <-closeErr; err != nil {
		syslog.Log(syslog.Error).Panicf("Failed to shutdown API server: %s", err)
	}
}

func LoggerExtractor(ctx context.Context) []any {
	attrs := []any{}

	tid, ok := ctxstore.From[string](ctx, ctxstore.TraceID)
	if ok {
		attrs = append(attrs, ctxstore.TraceID.String(), tid)
	}

	return attrs
}

type Adapters struct {
	*dbAdapter.UserDAO
	*flashAdapter.SessionManager
	*flashAdapter.LastSeen
	*dbAdapter.MessageDAO
}

func NewAdapters(db *database.DB, sessmng *flashAdapter.SessionManager, lseen *flashAdapter.LastSeen) *Adapters {
	return &Adapters{
		UserDAO:        dbAdapter.NewUserDAO(db),
		SessionManager: sessmng,
		LastSeen:       lseen,
		MessageDAO:     dbAdapter.NewMessageDAO(db),
	}
}

type Services struct {
	service.User
	service.Token
	service.Auth
	service.Message
}

func NewServices(conf config.Config, adapters *Adapters) *Services {
	user := service.NewUser(hashing.NewBcryptHasher(hashing.DefaultBcryptCost), adapters.UserDAO, adapters.UserDAO)
	token := service.NewToken()
	auth := service.NewAuth(conf.Auth, user, token, adapters.SessionManager)
	message := service.NewMessage(adapters.MessageDAO, adapters.MessageDAO)

	return &Services{
		User:    user,
		Token:   token,
		Auth:    auth,
		Message: message,
	}
}

type Handlers struct {
	*handler.Auth
	*handler.User
	*handler.Message
}

func NewHandlers(log *logging.Logger, services *Services) *Handlers {
	return &Handlers{
		Auth:    handler.NewAuth(log, services.Auth),
		User:    handler.NewUser(log, services.User),
		Message: handler.NewMessage(log, services.Message),
	}
}

type Middlewares struct {
	*middleware.Auth
}

func NewMiddlewares(log *logging.Logger, services *Services) *Middlewares {
	return &Middlewares{
		Auth: middleware.NewAuth(log, services.Auth),
	}
}

func SetupRoutes(log *logging.Logger, adapters *Adapters, handlers *Handlers, middlewares *Middlewares) http.Handler {
	router := transport.NewRouter()

	router.Use(middleware.RealIP, middleware.TraceID)
	router.Use(middleware.LogAccess(log))
	router.Use(middlewares.Authorize)
	router.Use(middleware.LastSeen(log, adapters.LastSeen))

	router.HandleFunc("/health", handler.Healthz)

	authRoute := router.Mount("/auth")
	{
		authRoute.HandleFunc("POST /login", handlers.Login())
		authRoute.HandleFunc("POST /register", handlers.Register())
		authRoute.With(middlewares.Protect).HandleFunc("POST /logout", handlers.Logout())
		authRoute.With(middlewares.Protect).HandleFunc("GET /refresh", handlers.Refresh())
	}

	userRoute := router.Mount("/users")
	{
		userRoute.Use(middlewares.Protect)

		userRoute.HandleFunc("GET /", handlers.User.Find())
	}

	messageRoute := router.Mount("/messages")
	{
		messageRoute.Use(middlewares.Protect)

		messageRoute.HandleFunc("GET /users/{userId}", handlers.Message.Find())
		messageRoute.HandleFunc("POST /users/{userId}", handlers.Message.Create())
	}

	return router
}
