package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/mauriciofsnts/hermes/internal/types"
	"github.com/pauloo27/logger"
)

var storageLocalName = "storage"

func CreateFiberInstance(storage types.Storage[types.Email]) *fiber.App {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals(storageLocalName, storage)
		return c.Next()
	})

	app.Use(limiter.New(limiter.Config{
		Max:        115,
		Expiration: 30 * time.Second,
	}))

	return app
}

func Listen(app *fiber.App) error {
	logger.Debug("Starting HTTP server...")

	app.Post("/api/send-email", SendEmail)
	app.Get("/api/health", HealthCheck)

	return app.Listen(":8082")
}
