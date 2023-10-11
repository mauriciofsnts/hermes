package http

import (
	"time"

	"github.com/Pauloo27/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func CreateFiberInstance() *fiber.App {
	app := fiber.New()

	app.Use(limiter.New(limiter.Config{
		Max:        5,
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
