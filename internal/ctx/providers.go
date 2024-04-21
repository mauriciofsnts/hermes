package ctx

import (
	"context"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/types"
)

type Key string

const (
	ProvidersKey Key = "providers"
)

type Providers struct {
	Config *config.Config
	Queue  types.Queue[types.Mail]
}

func GetProviders(ctx context.Context) *Providers {
	return ctx.Value(ProvidersKey).(*Providers)
}
