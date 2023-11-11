package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mauriciofsnts/hermes/internal/config"
)

func Origin() fiber.Handler {
	return func(c *fiber.Ctx) error {

		allowedOrigin := config.Hermes.AllowedOrigin

		if allowedOrigin != "" {
			origin := c.Get("Origin")

			if origin != allowedOrigin {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Origin not allowed",
				})
			}

			return c.Next()
		}

		return c.Next()
	}
}
