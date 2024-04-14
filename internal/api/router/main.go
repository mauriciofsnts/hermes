package router

import (
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

	app.Use(cors.New(cors.Config{
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

	app.Use(func(c *fiber.Ctx) error {
		apiKey := c.Get("x-api-key")
		apikeys := config.Envs.Hermes.Apikeys

		for _, v := range apikeys {
			if v == apiKey {
				return c.Next()
			}
		}

		return c.SendStatus(fiber.StatusUnauthorized)
	})

	app.Use(limiter.New(limiter.Config{
		Max:        config.Envs.Hermes.RateLimit,
		Expiration: 30 * time.Second,
	}))

	api.Post("/send", emailController.SendPlainTextEmail)
	api.Post("/send/:slug", emailController.SendTemplateEmail)

	api.Get("/templates/:slug/raw", templateController.GetRaw)
	api.Post("/templates", templateController.Create)

	return app.Listen(":8082")
}
