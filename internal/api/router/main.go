package router

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/mauriciofsnts/hermes/internal/api/controller"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/types"
)

const queueKey = "queue"

func CreateFiberInstance(queue types.Queue[types.Email]) *fiber.App {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals(queueKey, queue)
		return c.Next()
	})

	app.Use(limiter.New(limiter.Config{
		Max:        config.Hermes.RateLimit,
		Expiration: 30 * time.Second,
	}))

	allowedOrigins := config.Hermes.AllowedOrigins
	var parsedOrigin string

	if len(allowedOrigins) == 0 {
		parsedOrigin = "*"
	} else {
		parsedOrigin = strings.Join(allowedOrigins, ",")
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: parsedOrigin,
		AllowMethods: "POST,GET,OPTIONS",
	}))

	return app
}

func Listen(app *fiber.App) error {
	slog.Info("Starting HTTP server...")

	healthController := controller.NewHealthController()
	emailController := controller.NewEmailController()
	templateController := controller.NewTemplateController()

	api := app.Group("/api")

	api.Post("/send", emailController.SendEmail)
	api.Get("/health", healthController.Health)

	api.Get("/templates/:slug/raw", templateController.GetRaw)
	api.Post("/templates", templateController.Create)

	return app.Listen(":8082")
}
