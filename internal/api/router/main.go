package router

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/mauriciofsnts/hermes/internal/api/controller"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/types"
)

func CreateFiberInstance(queue types.Queue[types.Mail]) *fiber.App {
	app := fiber.New()

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
	healthController := controller.NewHealthController()
	emailController := controller.NewEmailController()
	templateController := controller.NewTemplateController()

	api := app.Group("/api")

	api.Get("/health", healthController.Health)

	api.Post("/send", emailController.SendPlainTextEmail)
	api.Post("/send/:slug", emailController.SendTemplateEmail)

	api.Get("/templates/:slug/raw", templateController.GetRaw)
	api.Post("/templates", templateController.Create)

	return app.Listen(":8082")
}
