package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mauriciofsnts/hermes/internal/api"
	"github.com/mauriciofsnts/hermes/internal/types"
)

type HealthControllerInterface interface {
	Health(c *fiber.Ctx) error
}

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (h *HealthController) Health(c *fiber.Ctx) error {
	queue := c.Locals("queue").(types.Queue[types.Mail])

	ping, err := queue.Ping()

	if err != nil {
		return api.Err(c, fiber.StatusInternalServerError, "Queue is not available", err)
	}

	return api.Success(c, fiber.StatusOK, ping)
}
