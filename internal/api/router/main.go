package router

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/controller"
	"github.com/mauriciofsnts/hermes/internal/types"
	"github.com/pauloo27/logger"
)

const storageKey = "storage"

func CreateFiberInstance(storage types.Storage[types.Email]) *fiber.App {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals(storageKey, storage)
		return c.Next()
	})

	app.Use(limiter.New(limiter.Config{
		Max:        115,
		Expiration: 30 * time.Second,
	}))

	allowedOrigins := config.Hermes.AllowedOrigins
	var parsedOrigin string

	if len(allowedOrigins) == 0 {
		parsedOrigin = "*"
	} else {
		parsedOrigin = strings.Join(allowedOrigins, ",")
	}

	logger.Infof("Allowed origins: %s", parsedOrigin)

	app.Use(cors.New(cors.Config{
		AllowOrigins: parsedOrigin,
		AllowMethods: "POST,GET,OPTIONS",
	}))

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
