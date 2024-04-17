package server

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/server/api/health"
	"github.com/mauriciofsnts/hermes/internal/server/api/notify"
	"github.com/mauriciofsnts/hermes/internal/server/api/template"
)

func CreateFiberInstance() *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{AllowMethods: "POST,GET,OPTIONS"}))

	return app
}

var rateLimitsByAPIKey = make(map[string](func(*fiber.Ctx) error))

func Listen(app *fiber.App) error {

	api := app.Group("/api/v1")

	healthController := health.NewHealthController()
	api.Get("/health", healthController.Health)

	app.Use(func(c *fiber.Ctx) error {
		apiKey := c.Get("x-api-key")

		appConfig, ok := config.Hermes.AppsByAPIKey[apiKey]
		if !ok {
			return c.Status(fiber.StatusUnauthorized).SendString("Chave API inválida")
		}

		rateLimiter, ok := rateLimitsByAPIKey[apiKey]
		if !ok {
			rateLimiter = limiter.New(limiter.Config{
				Max:        appConfig.LimitPerIPPerHour,
				Expiration: 1 * time.Hour,
				KeyGenerator: func(c *fiber.Ctx) string {
					return c.Get("x-api-key")
				},
			})
			rateLimitsByAPIKey[apiKey] = rateLimiter
		}

		if err := rateLimiter(c); err != nil {
			return c.Status(fiber.StatusTooManyRequests).SendString("Limite de requisições excedido")
		}

		return c.Next()
	})

	emailController := notify.NewEmailController()
	api.Post("/notify", emailController.SendPlainTextEmail)
	api.Post("/notify/:slug", emailController.SendTemplateEmail)

	templateController := template.NewTemplateController()
	api.Get("/templates/:slug/raw", templateController.GetRaw)
	api.Post("/templates", templateController.Create)

	return app.Listen(fmt.Sprintf(":%d", config.Hermes.Http.Port))
}
