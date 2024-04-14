package controller

import (
	"context"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/gofiber/fiber/v2"
	"github.com/mauriciofsnts/hermes/internal/api"
	"github.com/mauriciofsnts/hermes/internal/providers"
	"github.com/mauriciofsnts/hermes/internal/smtp"
)

type HealthControllerInterface interface {
	Health(c *fiber.Ctx) error
}

type HealthController struct {
	checker health.Checker
}

func NewHealthController() *HealthController {
	return &HealthController{
		checker: getChecker(),
	}
}

func getChecker() health.Checker {
	checker := health.NewChecker(
		health.WithCacheDuration(1*time.Second),
		health.WithTimeout(10*time.Second),

		health.WithCheck(health.Check{
			Name: "queue",
			Check: func(ctx context.Context) error {
				_, err := providers.Queue.Ping()

				if err != nil {
					return err
				}

				return nil
			},
		}),

		health.WithPeriodicCheck(15*time.Second, 3*time.Second, health.Check{
			Name: "smtp",
			Check: func(ctx context.Context) error {
				err := smtp.Ping()

				if err != nil {
					return err
				}

				return nil
			},
		}),
	)

	return checker
}

func (h *HealthController) Health(c *fiber.Ctx) error {
	result := h.checker.Check(context.Background())

	if result.Status == health.StatusUp {
		api.Success(c, fiber.StatusOK, result)
	} else {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Service unavailable", "detail": result})
	}

	return nil
}
