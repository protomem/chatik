package transport

import (
	"context"
	"net/http"
	"time"

	"github.com/go-pkgz/routegroup"
	"github.com/protomem/chatik/pkg/werrors"
)

type (
	ServerOptions struct {
		ListenAddr string `yaml:"listenAddr"`
		BaseAddr   string `yaml:"baseAddr"`

		Timeouts ServerTimeouts `yaml:"timeouts"`
	}

	ServerTimeouts struct {
		Read  time.Duration `yaml:"read"`
		Write time.Duration `yaml:"write"`
		Idle  time.Duration `yaml:"idle"`
	}
)

func DefaultServerOptions() ServerOptions {
	return ServerOptions{
		ListenAddr: ":8080",
		BaseAddr:   "localhost:8080",
		Timeouts: ServerTimeouts{
			Read:  15 * time.Second,
			Write: 15 * time.Second,
			Idle:  1 * time.Minute,
		},
	}
}

type Router = routegroup.Bundle

func NewRouter() *Router {
	return routegroup.New(http.NewServeMux())
}

type Server struct {
	opts ServerOptions
	srv  *http.Server
}

func NewServer(opts ServerOptions) *Server {
	srv := &http.Server{
		Addr:    opts.ListenAddr,
		Handler: http.NewServeMux(),

		ReadTimeout:  opts.Timeouts.Read,
		WriteTimeout: opts.Timeouts.Write,
		IdleTimeout:  opts.Timeouts.Idle,
	}

	return &Server{
		opts: opts,
		srv:  srv,
	}
}

func (s *Server) Addr() string {
	return s.srv.Addr
}

func (s *Server) SetHandler(h http.Handler) {
	if h == nil {
		return
	}
	s.srv.Handler = h
}

func (s *Server) Run() error {
	err := werrors.Filter(s.srv.ListenAndServe(), http.ErrServerClosed)
	return werrors.Error(err, "httpServer/run")
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.srv.Shutdown(ctx)
	return werrors.Error(err, "httpServer/shutdown")
}
