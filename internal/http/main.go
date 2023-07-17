package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func CreateFiberInstance() *fiber.App {
	app := fiber.New()

	app.Use(limiter.New(limiter.Config{
		Max:        1,
		Expiration: 30 * time.Second,
	}))

	return app
}

func Listen(app *fiber.App) error {

	app.Post("/api/send-email", SendEmail)
	app.Get("/api/health", HealthCheck)

	app.Listen(":8080")

	return nil
}
