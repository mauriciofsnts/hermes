package template

import (
	"bytes"
	"os"

	tmpl "html/template"
)

type TemplateServiceInterface interface {
	Exists(name string) bool
	Get(name string) ([]byte, error)
	Delete(name string) error
	Create(name string, content []byte) error
	ParseTemplate(name string, content map[string]interface{}) (*bytes.Buffer, error)
}

type TemplateService struct {
}

func NewTemplateService() TemplateServiceInterface {
	return &TemplateService{}
}

func (t *TemplateService) Exists(name string) bool {
	_, err := os.Stat(getPath(name))
	return err == nil
}

func (t *TemplateService) Get(name string) ([]byte, error) {
	return os.ReadFile(getPath(name))
}

func (t *TemplateService) Delete(name string) error {
	return os.Remove(getPath(name))
}

func (t *TemplateService) Create(name string, content []byte) error {
	return os.WriteFile(getPath(name), content, 0600)
}

func (t *TemplateService) ParseTemplate(name string, content map[string]any) (*bytes.Buffer, error) {
	html, err := t.Get(name)

	if err != nil {
		return nil, err
	}

	htmlTmpl, err := tmpl.New(name).Parse(string(html))

	if err != nil {
		return nil, err
	}

	buff := bytes.NewBufferString("")
	err = htmlTmpl.Option("missingkey=error").Execute(buff, content)

	if err != nil {
		return nil, err
	}

	return buff, nil
}

func getPath(name string) string {
	return "templates/" + name + ".html"
}
