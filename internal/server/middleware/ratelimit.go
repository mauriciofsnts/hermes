package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/mauriciofsnts/hermes/internal/config"
)

var rateLimitsByAPIKey = make(map[string](func(*fiber.Ctx) error))

func Ratelimit(c *fiber.Ctx) error {
	apiKey := c.Get("x-api-key")

	appConfig, ok := config.Hermes.AppsByAPIKey[apiKey]
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid API Key")
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
		return c.Status(fiber.StatusTooManyRequests).SendString("Too many requests")
	}

	return c.Next()
}
