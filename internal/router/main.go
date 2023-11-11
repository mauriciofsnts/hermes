package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/mauriciofsnts/hermes/internal/controller"
	"github.com/mauriciofsnts/hermes/internal/types"
	"github.com/pauloo27/logger"
)

const storageLocalName = "storage"

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

	app.Use(Origin())

	return app
}

func Listen(app *fiber.App) error {
	logger.Debug("Starting HTTP server...")

	healthController := controller.NewHealthController()
	emailController := controller.NewEmailController()
	templateController := controller.NewTemplateController()

	app.Post("/api/send-email", emailController.SendEmail)
	app.Get("/api/health", healthController.Health)

	app.Get("/api/templates/:slug/raw", templateController.GetRaw)
	app.Post("/api/templates", templateController.Create)

	return app.Listen(":8082")
}
