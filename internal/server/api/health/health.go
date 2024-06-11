package health

import (
	"context"
	"net/http"
	"time"

	healthCheck "github.com/alexliesenfeld/health"
	"github.com/mauriciofsnts/hermes/internal/providers/smtp"
	"github.com/mauriciofsnts/hermes/internal/server/api"
)

func (c *HealthController) GetHealth(r *http.Request) api.Response {

	statusChecker := healthCheck.NewChecker(
		healthCheck.WithCacheDuration(1*time.Second),
		healthCheck.WithTimeout(10*time.Second),

		healthCheck.WithCheck(healthCheck.Check{
			Name: "queue",
			Check: func(ctx context.Context) error {
				_, err := c.queue.Ping()

				if err != nil {
					return err
				}

				return nil
			},
		}),

		healthCheck.WithCheck(healthCheck.Check{
			Name: "smtp",
			Check: func(ctx context.Context) error {
				err := smtp.Ping()

				if err != nil {
					return err
				}

				return nil
			},
		}),
	).Check(context.Background())

	if statusChecker.Status == healthCheck.StatusUp {
		return api.Ok(statusChecker)
	} else {
		return api.DetailedError(api.InternalServerErr, statusChecker)
	}
}
