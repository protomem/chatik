package api

import "github.com/gofiber/fiber/v2"

func (srv *Server) handleHealth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	}
}
