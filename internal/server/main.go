package server

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/server/controller"
)

func CreateFiberInstance() *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowMethods: "POST,GET,OPTIONS",
	}))

	return app
}

func Listen(app *fiber.App) error {

	api := app.Group("/api/v1")

	healthController := controller.NewHealthController()
	api.Get("/health", healthController.Health)

	app.Use(func(c *fiber.Ctx) error {
		apiKey := c.Get("x-api-key")

		app := config.Hermes.AppsByAPIKey[apiKey]

		if app != nil {
			return c.Next()
		}

		return c.SendStatus(fiber.StatusUnauthorized)
	})

	app.Use(func(c *fiber.Ctx) error {
		apiKey := c.Get("x-api-key")

		app := config.Hermes.AppsByAPIKey[apiKey]

		if app != nil {
			return limiter.New(limiter.Config{
				Max:        1,
				Expiration: 1 * time.Hour,
			})(c)
		}

		return c.SendStatus(fiber.StatusUnauthorized)
	})

	emailController := controller.NewEmailController()
	api.Post("/notify", emailController.SendPlainTextEmail)
	api.Post("/notify/:slug", emailController.SendTemplateEmail)

	templateController := controller.NewTemplateController()
	api.Get("/templates/:slug/raw", templateController.GetRaw)
	api.Post("/templates", templateController.Create)

	return app.Listen(fmt.Sprintf(":%d", config.Hermes.Http.Port))
}
