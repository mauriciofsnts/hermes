package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/storage"
	"github.com/mauriciofsnts/hermes/internal/types"
)

func SendEmail(c *fiber.Ctx) error {
	allowedOrigin := config.Hermes.AllowedOrigin

	var email types.Email

	if allowedOrigin != "" {
		origin := c.Get("Origin")

		if origin != allowedOrigin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Origin not allowed",
			})
		}
	}

	if err := c.BodyParser(&email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// !TODO: move this for another place grr
	var producer = storage.NewStorage()

	err := producer.Write(email)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send email: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Email sent successfully",
	})
}

func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hermes is up and running",
	})
}
