package health

import "github.com/mauriciofsnts/hermes/internal/server/api"

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

var _ api.Controller = &HealthController{}

func (c *HealthController) Route(r api.Router) {
	r.Get("/", c.GetHealth)
}
