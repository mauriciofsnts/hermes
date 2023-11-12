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
	storage := c.Locals("storage").(types.Storage[types.Email])

	ping, err := storage.Ping()

	if err != nil {
		return api.Err(c, fiber.StatusInternalServerError, "failed to ping storage", err)
	}

	return api.Success(c, fiber.StatusOK, ping)
}
