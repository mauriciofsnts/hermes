package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/events"
	"github.com/mauriciofsnts/hermes/internal/types"
)

var producer = events.NewProducer[types.Email]()

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

	err := producer.Produce(uuid.New().String(), email, events.EmailTopic)

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
