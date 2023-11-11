package template

import (
	"os"

	"github.com/mauriciofsnts/hermes/internal/config"
)

type TemplateService interface {
	Exists(name string) bool
	Get(name string) ([]byte, error)
	Delete(name string) error
	Create(name string, content []byte) error
}

type templateService struct {
}

func NewTemplateService() TemplateService {
	return &templateService{}
}

func (t *templateService) Exists(name string) bool {
	_, err := os.Stat(getPath(name))
	return err == nil
}

func (t *templateService) Get(name string) ([]byte, error) {
	return os.ReadFile(getPath(name))
}

func (t *templateService) Delete(name string) error {
	return os.Remove(getPath(name))
}

func (t *templateService) Create(name string, content []byte) error {
	return os.WriteFile(getPath(name), content, 0644)
}

func getPath(name string) string {
	return config.Hermes.Location + "/" + name + ".html"
}
