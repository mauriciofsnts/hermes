package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mauriciofsnts/hermes/internal/api"
)

type HealthController interface {
	Health(c *fiber.Ctx) error
}

type healthController struct{}

func NewHealthController() *healthController {
	return &healthController{}
}

func (h *healthController) Health(c *fiber.Ctx) error {
	return api.Success(c, fiber.StatusOK, "OK")
}
