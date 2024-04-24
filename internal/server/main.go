package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/ctx"
	"github.com/mauriciofsnts/hermes/internal/server/api/health"
	"github.com/mauriciofsnts/hermes/internal/server/api/notify"
	"github.com/mauriciofsnts/hermes/internal/server/api/template"
	"github.com/mauriciofsnts/hermes/internal/server/middleware"
)

func CreateFiberInstance(providers *ctx.Providers) *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{AllowMethods: "POST,GET,OPTIONS"}))
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("providers", providers)
		return c.Next()
	})

	return app
}

func Listen(app *fiber.App) error {

	api := app.Group("/api/v1")

	healthController := health.NewHealthController()
	api.Get("/health", healthController.Health)

	app.Use(middleware.Ratelimit)

	emailController := notify.NewEmailController()
	api.Post("/notify", emailController.SendPlainTextEmail)
	api.Post("/notify/:slug", emailController.SendTemplateEmail)

	templateController := template.NewTemplateController()
	api.Get("/templates/:slug/raw", templateController.GetRaw)
	api.Post("/templates", templateController.Create)

	return app.Listen(fmt.Sprintf(":%d", config.Hermes.Http.Port))
}
