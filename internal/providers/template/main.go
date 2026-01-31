package template

import (
	"bytes"
	"os"
	"sync"

	tmpl "html/template"
)

type TemplateProvider interface {
	Exists(name string) bool
	Get(name string) ([]byte, error)
	Delete(name string) error
	Create(name string, content []byte) error
	ParseHtmlTemplate(name string, content map[string]interface{}) (*bytes.Buffer, error)
	ClearCache() error
}

type TemplateService struct {
	cache map[string]*tmpl.Template
	mu    sync.RWMutex
}

func NewTemplateService() TemplateProvider {
	return &TemplateService{
		cache: make(map[string]*tmpl.Template),
	}
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
	if err := os.WriteFile(getPath(name), content, 0600); err != nil {
		return err
	}
	// Invalidar cache quando um novo template é criado
	t.mu.Lock()
	delete(t.cache, name)
	t.mu.Unlock()
	return nil
}

func (t *TemplateService) ParseHtmlTemplate(name string, content map[string]any) (*bytes.Buffer, error) {
	// Verificar cache primeiro
	t.mu.RLock()
	cachedTmpl, exists := t.cache[name]
	t.mu.RUnlock()

	var htmlTmpl *tmpl.Template
	var err error

	if exists {
		htmlTmpl = cachedTmpl
	} else {
		// Carregar e fazer parse se não está em cache
		html, err := t.Get(name)
		if err != nil {
			return nil, err
		}

		htmlTmpl, err = tmpl.New(name).Parse(string(html))
		if err != nil {
			return nil, err
		}

		// Armazenar em cache
		t.mu.Lock()
		t.cache[name] = htmlTmpl
		t.mu.Unlock()
	}

	buff := bytes.NewBufferString("")
	err = htmlTmpl.Option("missingkey=error").Execute(buff, content)

	if err != nil {
		return nil, err
	}

	return buff, nil
}

// ClearCache limpa o cache de templates em memória
func (t *TemplateService) ClearCache() error {
	t.mu.Lock()
	t.cache = make(map[string]*tmpl.Template)
	t.mu.Unlock()
	return nil
}

func getPath(name string) string {
	return "templates/" + name + ".html"
}
