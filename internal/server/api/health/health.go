package health

import (
	"context"
	"net/http"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/mauriciofsnts/hermes/internal/providers/queue"
	"github.com/mauriciofsnts/hermes/internal/providers/smtp"
	"github.com/mauriciofsnts/hermes/internal/server/helper"
)

type HealthControllerInterface interface {
	Health(w http.ResponseWriter, r *http.Request)
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
				_, err := queue.Queue.Ping()

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

func (h *HealthController) Health(w http.ResponseWriter, r *http.Request) {
	result := h.checker.Check(context.Background())

	if result.Status == health.StatusUp {
		helper.Ok(w, result)
	} else {
		helper.DetailedError(w, helper.InternalServerErr, "Service unavailable")
	}
}
