package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mauriciofsnts/hermes/internal/api"
	"github.com/mauriciofsnts/hermes/internal/types"
)

type HealthController interface {
	Health(c *fiber.Ctx) error
}

type healthController struct{}

func NewHealthController() *healthController {
	return &healthController{}
}

func (h *healthController) Health(c *fiber.Ctx) error {
	queue := c.Locals("queue").(types.Queue[types.Email])

	ping, err := queue.Ping()

	if err != nil {
		return api.Err(c, fiber.StatusInternalServerError, "Queue is not available", err)
	}

	return api.Success(c, fiber.StatusOK, ping)
}
