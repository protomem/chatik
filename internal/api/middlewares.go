package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
