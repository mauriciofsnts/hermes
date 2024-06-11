package providers

import (
	"github.com/mauriciofsnts/hermes/internal/providers/queue/worker"
	"github.com/mauriciofsnts/hermes/internal/providers/template"
	"github.com/mauriciofsnts/hermes/internal/types"
	"gorm.io/gorm"
)

type Providers struct {
	DB      *gorm.DB
	Queue   worker.Queue[types.Mail]
	Storage template.TemplateProvider
}
