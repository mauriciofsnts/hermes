package providers

import (
	"github.com/mauriciofsnts/hermes/internal/providers/template"
	"github.com/mauriciofsnts/hermes/internal/types"
	"gorm.io/gorm"
)

type Providers struct {
	DB      *gorm.DB
	Queue   types.Queue[types.Mail]
	Storage template.TemplateProvider
}
