package health

import (
	"github.com/mauriciofsnts/hermes/internal/server/api"
	"github.com/mauriciofsnts/hermes/internal/types"
)

type HealthController struct {
	queue types.Queue[types.Mail]
}

func NewHealthController(queue types.Queue[types.Mail]) *HealthController {
	return &HealthController{
		queue: queue,
	}
}

var _ api.Controller = &HealthController{}

func (c *HealthController) Route(r api.Router) {
	r.Get("/", c.GetHealth)
}
